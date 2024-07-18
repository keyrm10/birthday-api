terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5.1"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}

module "cloud_sql" {
  source             = "./modules/cloud_sql"
  project_id         = var.project_id
  region             = var.region
  vpc_id             = module.networking.vpc_id
  instance_name      = var.db_instance_name
  database_version   = "POSTGRES_15"
  tier               = var.db_tier
  db_password_secret = google_secret_manager_secret.db_password.id
}

module "gke" {
  source                = "./modules/gke"
  project_id            = var.project_id
  region                = var.region
  vpc_id                = module.networking.vpc_id
  subnet_id             = module.networking.subnet_id
  cluster_name          = var.cluster_name
  node_pool_name        = var.node_pool_name
  service_account_email = module.iam.gke_service_account_email
}

module "iam" {
  source       = "./modules/iam"
  project_id   = var.project_id
  cluster_name = var.cluster_name
}

module "networking" {
  source     = "./modules/networking"
  project_id = var.project_id
  region     = var.region
}

resource "google_cloud_armor_policy" "policy" {
  name        = "birthday-api-armor-policy"
  description = "Cloud Armor policy for Birthday API"

  rule {
    description = "SQL Injection protection"
    action      = "deny(403)"
    priority    = "1000"
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('sqli-v33-stable')"
      }
    }
  }

  rule {
    description = "Cross-site scripting protection"
    action      = "deny(403)"
    priority    = "1001"
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('xss-v33-stable')"
      }
    }
  }

  rule {
    description = "Local file inclusion protection"
    action      = "deny(403)"
    priority    = "1002"
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('lfi-v33-stable')"
      }
    }
  }

  rule {
    description = "Remote file inclusion protection"
    action      = "deny(403)"
    priority    = "1003"
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('rfi-v33-stable')"
      }
    }
  }

  rule {
    description = "Remote code execution protection"
    action      = "deny(403)"
    priority    = "1004"
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('rce-v33-stable')"
      }
    }
  }

  rule {
    action      = "rate_based_ban"
    priority    = "1005"
    description = "Rate limiting rule"
    rate_limit_options {
      conform_action = "allow"
      exceed_action  = "deny(429)"
      enforce_on_key = "IP"
      rate_limit_threshold {
        count        = 100
        interval_sec = 60
      }
    }
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
  }

  rule {
    description = "default rule"
    action      = "allow"
    priority    = "2147483647"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
  }
}

resource "google_compute_global_address" "default" {
  name = "birthday-api-ip"
}

resource "google_compute_global_forwarding_rule" "default" {
  name                  = "birthday-api-forwarding-rule"
  target                = google_compute_target_http_proxy.default.id
  port_range            = "80"
  ip_address            = google_compute_global_address.default.address
  load_balancing_scheme = "EXTERNAL"
}

resource "google_compute_target_http_proxy" "default" {
  name    = "birthday-api-target-proxy"
  url_map = google_compute_url_map.default.id
}

resource "google_compute_url_map" "default" {
  name            = "birthday-api-url-map"
  default_service = google_compute_backend_service.default.id
}

resource "google_compute_backend_service" "default" {
  name            = "birthday-api-backend"
  protocol        = "HTTP"
  port_name       = "http"
  timeout_sec     = 10
  health_checks   = [google_compute_health_check.default.id]
  security_policy = google_cloud_armor_policy.policy.id

  backend {
    group = module.gke.node_pool_instance_group
  }
  depends_on = [module.gke]
}

resource "google_compute_health_check" "default" {
  name               = "http-health-check"
  check_interval_sec = 5
  timeout_sec        = 5

  http_health_check {
    port         = "8080"
    request_path = "/healthz"
  }
}

resource "kubernetes_secret" "db_secret" {
  metadata {
    name = "birthday-api-db-secret"
  }

  data = {
    password = data.google_secret_manager_secret_version.db_password.secret_data
  }
}

resource "helm_release" "birthday_api" {
  name      = "birthday-api"
  chart     = "../helm"
  namespace = "default"

  set {
    name  = "database.instanceConnectionName"
    value = module.cloud_sql.connection_name
  }

  set {
    name  = "database.secretName"
    value = kubernetes_secret.db_secret.metadata[0].name
  }
}

resource "random_password" "db_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "google_secret_manager_secret" "db_password" {
  secret_id = "db-password"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "db_password" {
  secret      = google_secret_manager_secret.db_password.id
  secret_data = random_password.db_password.result
}
