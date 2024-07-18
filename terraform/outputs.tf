output "gke_cluster_name" {
  value       = module.gke.cluster_name
  description = "The name of the GKE cluster"
}

output "cloud_sql_instance_name" {
  value       = module.cloud_sql.instance_name
  description = "The name of the Cloud SQL instance"
}

output "cloud_sql_connection_name" {
  value       = module.cloud_sql.connection_name
  description = "The connection name of the Cloud SQL instance"
}

output "load_balancer_ip" {
  value       = google_compute_global_address.default.address
  description = "The IP address of the global load balancer"
}
