CLUSTER = lightning-playground

# Load the Docker images into the kind cluster
load-images:
	@VERSION=$$(cat version.txt); \
	IMAGES=" \
		LMare/lightning-playground-frontend:$$VERSION \
		LMare/lightning-playground-backend:$$VERSION \
		LMare/lightning-playground-sidecar-backend:$$VERSION \
		LMare/lightning-playground-sidecar-btcd:$$VERSION \
		LMare/lightning-playground-sidecar-lnd:$$VERSION \
		LMare/duckdb-elt-service:$$VERSION \
		LMare/lnd:v0.20.0-beta-custom \
		btcsuite/btcd:v0.25.0"; \
	for img in $$IMAGES; do \
		kind load docker-image $$img --name $(CLUSTER); \
	done

create-cluster:
	@if kind get clusters | grep -q $(CLUSTER); then \
		echo "Cluster '$(CLUSTER)' already exists. Skipping creation."; \
	else \
		kind create cluster --name $(CLUSTER) --config kubernetes/cluster/cluster.yaml; \
		echo "cluster created"; \
	fi

delete-cluster:
	kind delete cluster --name $(CLUSTER)

deploy-lightning-playground:
	kubectl apply -f kubernetes/manifests/namespace.yaml; \
	while [ "$$(kubectl get ns lightning-playground -o jsonpath='{.status.phase}')" != "Active" ]; do sleep 1; done; \
	kubectl apply -R -f kubernetes/manifests; \
	kubectl config set-context --current --namespace=lightning-playground


bake-all:
	docker buildx bake

bake-app:
	docker buildx bake backend frontend sidecars data


# This target adds an entry to the /etc/hosts file to map lightning-playground.local
add-in-hosts:
	@app="127.0.0.1	lightning-playground.local"; \
	grafana="127.0.0.1	grafana.lightning-playground.local"; \
	if ! grep -q "$$app" /etc/hosts; then \
		echo "Adding '$$app' to /etc/hosts..."; \
		echo "$$app" | sudo tee -a /etc/hosts; \
		echo "done"; \
	fi; \
	if ! grep -q "$$grafana" /etc/hosts; then \
		echo "Adding '$$grafana' to /etc/hosts..."; \
		echo "$$grafana" | sudo tee -a /etc/hosts; \
		echo "done"; \
	fi

all: check-deps bake-all create-cluster load-images deploy-postgres deploy-lightning-playground deploy-minio add-in-hosts create-monitoring-cluster
refresh-new-version: bump-version bake-app load-images deploy-lightning-playground

check-deps:
	@command -v docker >/dev/null 2>&1 || { echo "docker is missing"; exit 1; } && \
	command -v kind   >/dev/null 2>&1 || { echo "kind is missing"; exit 1; } && \
	command -v kubectl >/dev/null 2>&1 || { echo "kubectl is missing"; exit 1; } && \
	docker buildx version >/dev/null 2>&1 || { echo "docker buildx is missing"; exit 1; } && \
	helm version >/dev/null 2>&1 || { echo "helm is missing"; exit 1; } && \
	echo "All dependencies are installed."

# Increment the version number in all relevant files or set it to a specific value if provided
bump-version:
	@old=$$(cat version.txt); \
	if [ -n "$(VERSION)" ]; then \
		new="$(VERSION)"; \
	else \
		new=$$(echo $$old | awk -F- '{print $$1 "-" $$2+1}'); \
	fi; \
	echo "Ancienne version: $$old"; \
	echo "Nouvelle version: $$new"; \
	echo $$new > version.txt; \
	find . -type f \( -name "*.yaml" -o -name "*.yml" -o -name "*.env" -o -name "docker-bake.*" \) \
		-exec sed -i "s/$$old/$$new/g" {} + && \
	echo "Version updated in all files."


deploy-monitoring-stack:
	helm dependency update kubernetes/helm/monitoring
	helm upgrade --install monitoring kubernetes/helm/monitoring -n monitoring

deploy-postgres:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm upgrade --install postgres bitnami/postgresql -n lightning-playground \
		--set auth.postgresPassword=superpassword \
		--set auth.username=appuser \
		--set auth.password=apppassword \
		--set auth.database=appdb


deploy-minio:
	helm repo add minio https://charts.min.io/
	helm repo update
	helm upgrade --install minio minio/minio -n data -f kubernetes/helm/minio/values.yaml \
		--set accessKey.password=superpassword \
		--set secretKey.password=superpassword


.PHONY: load-images create-cluster delete-cluster deploy-lightning-playground add-in-hosts all check-deps bake-all bake-app bump-version refresh-new-version deploy-monitoring-stack deploy-postgres deploy-minio
