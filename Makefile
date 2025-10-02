.PHONY: help build up down logs clean kind-up kind-down k8s-deploy k8s-delete docker-build

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Docker Compose commands
build: ## Build all Docker images
	docker compose build

up: ## Start all services with Docker Compose
	docker compose up -d

down: ## Stop all services
	docker compose down

logs: ## Show logs from all services
	docker compose logs -f

clean: ## Remove all containers, volumes, and images
	docker compose down -v --rmi all

# Individual service builds
build-probe-api: ## Build probe-api image
	docker build -t pulseguard/probe-api:latest ./probe-api

build-anomaly-worker: ## Build anomaly-worker image
	docker build -t pulseguard/anomaly-worker:latest ./anomaly-worker

build-webapp: ## Build webapp image
	docker build -t pulseguard/webapp:latest ./webapp

build-notifier: ## Build notifier image
	docker build -t pulseguard/notifier:latest ./notifier

docker-build: build-probe-api build-anomaly-worker build-webapp build-notifier ## Build all service images

# Kubernetes commands
kind-up: ## Create kind cluster
	kind create cluster --name pulseguard

kind-down: ## Delete kind cluster
	kind delete cluster --name pulseguard

kind-load: docker-build ## Load images into kind cluster
	kind load docker-image pulseguard/probe-api:latest --name pulseguard
	kind load docker-image pulseguard/anomaly-worker:latest --name pulseguard
	kind load docker-image pulseguard/webapp:latest --name pulseguard
	kind load docker-image pulseguard/notifier:latest --name pulseguard

k8s-deploy: ## Deploy to Kubernetes
	kubectl apply -f deploy/k8s/base/

k8s-delete: ## Delete Kubernetes resources
	kubectl delete -f deploy/k8s/base/

k8s-status: ## Show status of all pods
	kubectl get pods -n pulseguard

k8s-logs-probe: ## Show probe-api logs
	kubectl logs -n pulseguard -l app=probe-api -f

k8s-logs-worker: ## Show anomaly-worker logs
	kubectl logs -n pulseguard -l app=anomaly-worker -f

k8s-logs-webapp: ## Show webapp logs
	kubectl logs -n pulseguard -l app=webapp -f

k8s-logs-notifier: ## Show notifier logs
	kubectl logs -n pulseguard -l app=notifier -f

port-forward-webapp: ## Forward webapp port to localhost:3000
	kubectl -n pulseguard port-forward svc/webapp 3000:80

port-forward-api: ## Forward probe-api port to localhost:8080
	kubectl -n pulseguard port-forward svc/probe-api 8080:80

# Development helpers
dev-probe-api: ## Run probe-api locally
	cd probe-api && go run cmd/probe-api/main.go

dev-webapp: ## Run webapp locally
	cd webapp && npm run dev

dev-worker: ## Run anomaly-worker locally
	cd anomaly-worker && python app.py

dev-notifier: ## Run notifier locally
	cd notifier && go run cmd/notifier/main.go
