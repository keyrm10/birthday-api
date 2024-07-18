.PHONY: deploy deploy-infra deploy-app undeploy-app destroy-infra cleanup help

# Variables
DOCKER_REPO := ghcr.io/keyrm10/birthday-api
GCP_PROJECT := $(shell gcloud config get-value project)
GKE_CLUSTER := birthday-api-cluster
GKE_ZONE := europe-west1-b
TF_DIR := ./terraform
HELM_DIR := ./helm

# Deploy infrastructure using Terraform
deploy-infra:
	@echo "Deploying infrastructure..."
	@cd $(TF_DIR) && terraform init
	@cd $(TF_DIR) && terraform apply -auto-approve

# Deploy application to GKE
deploy-app: get-gke-credentials
	@echo "Deploying application to GKE..."
	@helm upgrade --install birthday-api $(HELM_DIR) \
		--set image.repository=$(DOCKER_REPO) \
		--set image.tag=latest

# Undeploy application from GKE
undeploy-app:
	@echo "Undeploying application from GKE..."
	@helm uninstall birthday-api

# Clean up resources
cleanup:
	@echo "Cleaning up resources..."
	@cd $(TF_DIR) && terraform destroy -auto-approve

# Get GKE cluster credentials
get-gke-credentials:
	@echo "Fetching GKE cluster credentials..."
	@gcloud container clusters get-credentials $(GKE_CLUSTER) --zone $(GKE_ZONE) --project $(GCP_PROJECT)

# Help target
help:
	@echo "Available targets:"
	@echo "  deploy-infra  - Deploy infrastructure using Terraform"
	@echo "  deploy-app    - Deploy application to GKE"
	@echo "  undeploy-app  - Undeploy application from GKE"
	@echo "  cleanup       - Clean up resources"
	@echo "  help          - Show this help message"
