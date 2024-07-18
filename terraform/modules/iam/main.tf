resource "google_service_account" "gke_workload_identity" {
  account_id   = "${var.cluster_name}-wi-sa"
  display_name = "GKE Workload Identity Service Account for ${var.cluster_name}"
}

resource "google_project_iam_member" "gke_workload_identity" {
  for_each = toset([
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
    "roles/monitoring.viewer",
    "roles/stackdriver.resourceMetadata.writer"
  ])

  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.gke_workload_identity.email}"
}

resource "google_service_account_iam_binding" "workload_identity_binding" {
  service_account_id = google_service_account.gke_workload_identity.name
  role               = "roles/iam.workloadIdentityUser"
  members = [
    "serviceAccount:${var.project_id}.svc.id.goog[default/default]"
  ]
}
