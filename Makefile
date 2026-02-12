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

deploy:
	kubectl apply -f kubernetes/manifests/namespace.yaml; \
	while [ "$$(kubectl get ns lightning-playground -o jsonpath='{.status.phase}')" != "Active" ]; do sleep 1; done; \
	kubectl apply -R -f kubernetes/manifests; \
	kubectl config set-context --current --namespace=lightning-playground


bake-all:
	docker buildx bake

bake-app:
	docker buildx bake backend frontend sidecars


# This target adds an entry to the /etc/hosts file to map lightning-playground.local
add-in-hosts:
	@line="127.0.0.1	lightning-playground.local"; \
	if ! grep -q "$$line" /etc/hosts; then \
		echo "Adding '$$line' to /etc/hosts..."; \
		echo "$$line" | sudo tee -a /etc/hosts; \
		echo "done"; \
	fi

all: bake-all create-cluster load-images deploy add-in-hosts
refresh-new-version: bump-version bake-app load-images deploy

check-deps:
	@command -v docker >/dev/null 2>&1 || { echo "docker is missing"; exit 1; } && \
	command -v kind   >/dev/null 2>&1 || { echo "kind is missing"; exit 1; } && \
	command -v kubectl >/dev/null 2>&1 || { echo "kubectl is missing"; exit 1; } && \
	docker buildx version >/dev/null 2>&1 || { echo "docker buildx is missing"; exit 1; } && \
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


.PHONY: load-images create-cluster delete-cluster deploy add-in-hosts all check-deps bake-all bake-app bump-version refresh-new-version 
