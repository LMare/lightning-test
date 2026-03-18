#!/bin/bash
# Force PREFECT_API_URL avant de lancer Prefect
export PREFECT_API_URL="http://prefect-server.data.svc.cluster.local:4200/api"
exec "$@"
