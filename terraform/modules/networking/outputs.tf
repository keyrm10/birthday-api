output "vpc_id" {
  value       = google_compute_network.vpc_network.id
  description = "The ID of the VPC"
}

output "subnet_id" {
  value       = google_compute_subnetwork.subnet.id
  description = "The ID of the subnet"
}
