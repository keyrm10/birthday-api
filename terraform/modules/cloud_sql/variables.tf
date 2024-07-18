variable "project_id" {
  description = "The project ID to deploy to"
  type        = string
}

variable "region" {
  description = "The region to deploy to"
  type        = string
}

variable "instance_name" {
  description = "The name of the Cloud SQL instance"
  type        = string
}

variable "database_version" {
  description = "The database version to use"
  type        = string
}

variable "tier" {
  description = "The machine type of the Cloud SQL instance"
  type        = string
}

variable "vpc_id" {
  description = "The ID of the VPC to which the Cloud SQL instance should be conected"
  type        = string
}

variable "db_password_secret" {
  description = "The secret ID for the database password"
  type        = string
}
