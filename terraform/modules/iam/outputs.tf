output "gke_service_account_email" {
  value       = google_service_account.gke_workload_identity.email
  description = "The email address of the GKE Workload Identity service account"
}
