output "cluster_name" {
  value       = google_container_cluster.primary.name
  description = "The name of the GKE cluster"
}

output "node_pool_instance_group" {
  value       = google_container_node_pool.primary_nodes.instance_group_urls[0]
  description = "The instance group URL of the GKE node pool"
}
