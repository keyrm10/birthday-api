data "google_secret_manager_secret_version" "db_password" {
  secret = var.db_password_secret
}

resource "google_sql_database_instance" "instance" {
  name             = var.instance_name
  region           = var.region
  database_version = var.database_version

  settings {
    tier = var.tier

    ip_configuration {
      ipv4_enabled    = false
      private_network = var.vpc_id
    }

    backup_configuration {
      enabled                        = true
      start_time                     = "02:00"
      location                       = var.region
      point_in_time_recovery_enabled = true
    }
  }
}

resource "google_sql_database" "database" {
  name       = "birthdays"
  instance   = google_sql_database_instance.instance.name
  depends_on = [google_sql_database_instance.instance]
}

resource "google_sql_user" "users" {
  name       = "postgres"
  instance   = google_sql_database_instance.instance.name
  password   = data.google_secret_manager_secret_version.db_password.secret_data
  depends_on = [google_sql_database_instance.instance]
}
