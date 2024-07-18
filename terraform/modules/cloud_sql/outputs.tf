output "instance_name" {
  value       = google_sql_database_instance.instance.name
  description = "The name of the Cloud SQL instance"
}

output "connection_name" {
  value       = google_sql_database_instance.instance.connection_name
  description = "The connection name of the Cloud SQL instance"
}
