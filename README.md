
# Lightning Playground


[![CI](https://github.com/LMare/lightning-playground/actions/workflows/ci.yml/badge.svg)](https://github.com/LMare/lightning-playground/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/LMare/lightning-playground/branch/master/graph/badge.svg)](https://codecov.io/gh/LMare/lightning-playground)



This project started as a small Go/HTMX web app interacting with a Lightning Network node on simnet.
It has since grown into a full Kubernetes playground featuring Prometheus/Grafana monitoring and a modern ELT pipeline (Postgres → MinIO → DuckDB → dbt) orchestrated with Prefect.

## Fonctionnal Prupose
Little web application which interract with a network of lightning servers running on simnet to create channel and  send transactions between nodes.


## Checking the prerequis
```bash
make check-deps
```

The application can be launch with `docker compose` or `Kubernetes (kind)` 


## Lauch the app

```bash
# Docker compose
docker buildx bake
docker compose up -d
# Kubernetes
make all
```
Go to : http://localhost:3000/ (docker compose)

Go to : http://lightning-playground.local/ (Kubernetes)

Grafana (Kubernetes) : 
| Clé      | Valeur |
|----------|--------|
| URL      | https://grafana.lightning-playground.local/d/admzlxj/lightning-playground-backend-dashboard?orgId=1&from=now-30m&to=now&timezone=browser&refresh=5s |
| Login    | `admin` |
| Password | `kubectl get secret -n monitoring monitoring-grafana -o jsonpath="{.data.admin-password}" \| base64 --decode` |


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
Copy the address then replace the value of `miningaddr` in the service `btcd` of `docker-compose.yml` or `kubernetes/manifests/btcd/btcd-statefulset.yaml`.
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

### Adding lnd node (increase de replicas number)
```bash
kubectl scale statefulset lnd --replicas=5
```

---------------------------------------------------------------------------------------------

## Note :
The app works with a LND customised [LMare/lnd](https://github.com/LMare/lnd/tree/feature/gRPC-alias-color).
This version allow to modify the alias and the color of the node LND by gRPC call.

-------------------------------------------------------------------------------------------
## Working in progress

StatefulSet: duckdb
  - [] container: duckdb-sql-server : execute and return des SQL request on the warehouse.duckdb
  - [X] container: duckdb-elt-service : mini server for the 
      * [X] load with incremental merge
      * [X] dbt
        - [] dbt scripts 
  - [X] volume: warehouse.duckdb

StatefulSet: [X] minio (helm)

StatefulSet: [X] postges (helm)

Deployment: prefect-agent
  - [X] prefect-server : orchestrate flow 
  - [X] prefect-agent  : client that call the server in loop and then create ephemeral pod to execute the flow
  - [] flow
       - [X] export posgres to minio : format parquet
       - [X] trigger duckdb-elt-service : load
       - [X] trigger duckdb-elt-service : dbt
  - [] load the flow in the server

Deployment: [] Dashboard backend  : go app that call duckdb-sql-server to return analytics data
Deployment: [] Dashboard frontend : with React

Deployment: Backend (lightning)
  - [] trace event in posgres to have data for the process   


---------------------------------------------------------------------------------------------
## TODO 
  - Warehouse + ELT : Stack identified : DuckDB + dbt + Perfect (Airflow in a second time ?)
  - archi hexagonal to complete TI for test cover
  - monitoring (alerts)
  - monitoring (logs)
  - SSE with Deployement Backend broken -> need to put a broker (ex: Redis)
  - Functionnal : 
    - Wallet creation/unlock handle by the backend via gRPC (`WalletUnlocker` service)
    - Close channel
    - Display amount
  - Define NetworkPolicies to restrict communication paths (frontend ↔ backend ↔ LND ↔ btcd) ?
  - infra as code (Terraform) ?
  - LLM

