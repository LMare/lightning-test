
# Lightning Playground


[![CI](https://github.com/LMare/lightning-playground/actions/workflows/ci.yml/badge.svg)](https://github.com/LMare/lightning-playground/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/LMare/lightning-playground/branch/master/graph/badge.svg)](https://codecov.io/gh/LMare/lightning-playground)



Personnal projet to discover and improve skill on  :
  - Golang
  - gRPC (use & extend gRPC API)
  - Lnd
  - HTMX
  - SSE
  - dockerfile
  - docker compose
  - docker bake
  - CI

TODO :
  - CD / infra as code (Terraform + kubernetes with Kind)
  - modules with Go
  - increase the test cover by implementing TI

---------------------------------------------------------------------------------------------

## TODO List for Kubernetes Migration (Lightning Network Project)

issues :
 - SSE with Deployement Backend broken -> need to put a broker

### 1. Cluster Setup
- [X] Spin up a local Kubernetes cluster (Kind, Minikube, or k3s).
- [X] Configure `kubectl` and create a dedicated namespace (e.g. `lightning`).

### 2. Core Components
- [X] **Frontend:** Deployment + Service + Ingress (stateless, scalable).
- [X] **Backend:** Deployment + Service (stateless, responsible for discovering LND pods and unlocking/creating wallets).
- [X] **btcd:** StatefulSet + PVC + Service (Bitcoin full node).
- [X] **LND:** StatefulSet + PVC + headless Service (multiple replicas, each with its own wallet/certs).

### 3. Data & Secrets Management
- [X] Define **PersistentVolumeClaims** for each LND and btcd.
- [X] Create a **Secret bundle** (`lnd-credentials`) to store certs/macaroons for all **LND** pods.
- [X] Create a **Secret bundle** (`btcd-credentials`) to store certs for all (only one BTCD for now)  **BTCD** pods.
- [X] Implement a **lnd-sidecar** that copies certs/macaroons from **LND** PVCs into the `lnd-credentials`.
- [X] Implement a **lnd-sidecar**  that copies certs from `btcd-credentials` into the **LND**
- [X] Implement a **backend-sidecar**  that copies certs/macaroons from `lnd-credentials` into the **Backend**.
- [X] Implement a **btcd-sidecar** that copies certs from **BTCD** PVCs into the `btcd-credentials`.
- [X] Mount the `btcd-credentials` in the **lnd-sidecar** in read-only mode.
- [X] Mount the `lnd-credentials` in the **backend-sidecar** in read-only mode.

### 4. Backend Responsibilities
- [X] Discover LND pods.
- [X] **backend-sidecar** Read the corresponding certs/macaroons from the Secret bundle.
- [ ] Use gRPC to **create or unlock wallets** via the `WalletUnlocker` service.
- [X] Replace static `nodes.yaml` with dynamic discovery logic.

### 5. Networking & Service Discovery
- [X] Configure a **headless Service** for LND to provide stable DNS per pod.
- [X] Ensure the backend can dynamically map endpoints (`lnd-N`) to certs/macaroons.
- [X] Use Ingress to expose frontend APIs externally.

### 6. Security & Best Practices
- [X] Restrict RBAC permissions for the job that **updates** the Secret bundle **LND**.
- [X] Restrict RBAC permissions for the job that **read** the Secret bundle **LND**.
- [X] Separate configs: ConfigMaps for non-sensitive data, Secrets for sensitive data.
- [ ] Add liveness/readiness probes for backend and LND.

### 7. Scalability & Monitoring
- [X] Test scaling: `kubectl scale statefulset lnd --replicas=5`
- [X] Verify the backend adapts automatically to new pods.
- [ ] Add monitoring (Prometheus + Grafana) and centralized logging.
- [ ] Define NetworkPolicies to restrict communication paths (frontend ↔ backend ↔ LND ↔ btcd).

### 8. Finalization
- [ ] Organize manifests into folders (`frontend/`, `backend/`, `lnd/`, `btcd/`).
- [ ] Deploy everything with `kubectl apply -f ./manifests`.
- [X] Validate end-to-end flow: frontend → backend → LND → btcd.
- [ ] Document the workflow for reproducibility (CI/CD, Helm charts, etc.).

---

## ✅ Expected Result

- Scaling LND pods is done with a single command (`kubectl scale`).
- Each LND pod has its own wallet and certs, stored securely in PVCs and synced into a Secret bundle.
- The backend dynamically discovers LND pods via DNS and uses the correct certs/macaroons from the Secret bundle.
- Wallet creation/unlock is handled by the backend via gRPC (`WalletUnlocker` service).
- Frontend and backend are exposed externally via Ingress.
- Secrets are mounted read-only, RBAC is restricted, and configs are separated (ConfigMap vs Secret).
- Monitoring, probes, and NetworkPolicies ensure production-grade reliability and observability.
- The architecture is clean, Kubernetes-native, and ready to evolve toward production.


-----------------------------------------------------------------------------------------------
## Deploiment avec Kubernetes IN Docker
```bash
docker buildx bake
echo "127.0.0.1	lightning-playground.local" >> /etc/hosts
cd kubernetes
export VERSION=v0.2.0-14
# cluster + docker images
kind create cluster --name lightning-playground --config cluster.yaml
kind load docker-image LMare/lightning-playground-frontend:$VERSION        --name lightning-playground &
kind load docker-image LMare/lightning-playground-backend:$VERSION         --name lightning-playground &
kind load docker-image LMare/lightning-playground-sidecar-backend:$VERSION --name lightning-playground &
kind load docker-image LMare/lightning-playground-sidecar-btcd:$VERSION    --name lightning-playground &
kind load docker-image LMare/lightning-playground-sidecar-lnd:$VERSION     --name lightning-playground &
kind load docker-image LMare/lnd:v0.20.0-beta-custom                       --name lightning-playground &
kind load docker-image btcsuite/btcd:v0.25.0                               --name lightning-playground &

kubectl create namespace lightning-playground
# frontend
kubectl apply -f frontend-deployment.yaml -n lightning-playground
kubectl apply -f frontend-service.yaml    -n lightning-playground
kubectl apply -f frontend-ingress.yaml    -n lightning-playground
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
# backend
kubectl apply -f backend-deployment.yaml      -n lightning-playground
kubectl apply -f backend-service.yaml         -n lightning-playground
kubectl apply -f backend-service-account.yaml -n lightning-playground
# btcd
kubectl apply -f btcd-cm1-configmap.yaml    -n lightning-playground
kubectl apply -f btcd-role-binding.yaml     -n lightning-playground
kubectl apply -f btcd-role.yaml             -n lightning-playground
kubectl apply -f btcd-secret.yaml           -n lightning-playground
kubectl apply -f btcd-service-account.yaml  -n lightning-playground
kubectl apply -f btcd-service-headless.yaml -n lightning-playground
kubectl apply -f btcd-statefulset.yaml      -n lightning-playground
# lnd
kubectl apply -f lnd-rolebinding.yaml      -n lightning-playground
kubectl apply -f lnd-role.yaml             -n lightning-playground
kubectl apply -f lnd-secret.yaml           -n lightning-playground
kubectl apply -f lnd-service-account.yaml  -n lightning-playground
kubectl apply -f lnd-service-headless.yaml -n lightning-playground
kubectl apply -f lnd-statefulset.yaml      -n lightning-playground
```
Go to http://lightning-playground.local/
-----------------------------------------------------------------------------------------------

## Prupose
Be able to do a little web application to interract with and a lightning serveur running on simnet

## Lauch the app

```bash
docker buildx bake
docker compose up -d
```
Go to : http://localhost:3000/

### First launch
To use the lnd fonctionalities, you will need at least 2 lnd nodes with a wallet :
```bash
# Docker compose
docker exec -it lightning-playground-lnd1-1 lncli --network=simnet create
docker exec -it lightning-playground-lnd2-1 lncli --network=simnet create
# Kubernetes
kubectl exec -it lnd-0 -n lightning-playground -- lncli --network=simnet create
kubectl exec -it lnd-1 -n lightning-playground -- lncli --network=simnet create
kubectl exec -it lnd-2 -n lightning-playground -- lncli --network=simnet create
```

To mine with one of these address do :
```bash
# Docker compose
docker exec -it lightning-playground-lnd1-1 lncli --network=simnet newaddress np2wkh
# Kubernetes
kubectl exec -it lnd-0 -n lightning-playground -- lncli --network=simnet newaddress np2wkh
```
Copy the address then replace the value of `miningaddr` in the service `btcd` of `docker-compose.yml` or `kubernetes/btcd-statefulset.yaml`.
And reload the containers
```bash
# Docker compose
docker compose up -d
# Kubernetes
kubectl apply -f btcd-statefulset.yaml -n lightning-playground
kubectl delete pod btcd-0 -n lightning-playground
```

Mine enough block to activate taproot
```bash
# Docker compose
docker exec -it lightning-playground-btcd-1 btcctl --simnet generate 1500
# Kubernetes
kubectl exec -it btcd-0 -n lightning-playground -- btcctl --simnet generate 1500
```

### Unlock the wallet
After each up of the lnd containers, the wallet must be unlock
```bash
# Docker compose
docker exec -it lightning-playground-lnd1-1 lncli --network=simnet unlock
docker exec -it lightning-playground-lnd2-1 lncli --network=simnet unlock
# Kubernetes
kubectl exec -it lnd-0 -n lightning-playground -- lncli --network=simnet unlock
kubectl exec -it lnd-1 -n lightning-playground -- lncli --network=simnet unlock
kubectl exec -it lnd-2 -n lightning-playground -- lncli --network=simnet unlock
```

### Generate a new block
In the simnet network the news blocks must to be mine manually. Run this cmd (every 10 minutes) to keep the lnd node synchronised with the btcd node    
```bash
# Docker compose
docker exec -it lightning-playground-btcd-1 btcctl --simnet generate 1
# Kubernetes
kubectl exec -it btcd-0 -n lightning-playground -- btcctl --simnet generate 1
```


## Stop the app
```bash
# Docker compose
docker-compose down
# Kubernetes
docker stop lightning-playground-control-plane
```

## Using the application

![Lightning-Playground](https://github.com/LMare/lightning-playground/blob/master/Lightning-Playground.png)


Steps :
 1. Add new pair with the URI of other nodes
 2. Create channels between pairs.
    After creating the channel to pass it in `active` state, generate some blocs with :
```bash
# Docker compose
docker exec -it lightning-playground-btcd-1 btcctl --simnet generate 10
# Kubernetes
kubectl exec -it btcd-0 -n lightning-playground -- btcctl --simnet generate 1
```
  3. Generate an invoice
  4. Import the invoice on another node and pay-it



## Note :
The app works with a LND customised [LMare/lnd](https://github.com/LMare/lnd/tree/feature/gRPC-alias-color).
This version allow to modify the alias and the color of the node LND by gRPC call.
