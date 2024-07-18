variable "project_id" {
  description = "The project ID to deploy to"
  type        = string
}

variable "region" {
  description = "The region to deploy to"
  type        = string
  default     = "europe-west1"
}

variable "cluster_name" {
  description = "The name of the GKE cluster"
  type        = string
  default     = "birthday-api-cluster"
}

variable "node_pool_name" {
  description = "The name of the GKE node pool"
  type        = string
  default     = "birthday-api-node-pool"
}

variable "db_instance_name" {
  description = "The name of the Cloud SQL instance"
  type        = string
  default     = "birthday-api-db"
}

variable "db_tier" {
  description = "The machine type of the Cloud SQL instance"
  type        = string
  default     = "db-f1-micro"
}
