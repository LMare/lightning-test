
FROM python:3.11-slim AS builder

# Install system deps (libpq-dev for psycopg, build-essential for pip)
RUN apt-get update \
 && apt-get install -y build-essential libpq-dev \
 && rm -rf /var/lib/apt/lists/*

# Install Python deps in a custom folder
COPY requirements.txt ./
RUN pip install --target=/custom-packages -r requirements.txt --no-cache-dir

FROM prefecthq/prefect:3-latest

# Copier les dépendances Python sans écraser /usr/local
COPY --from=builder /custom-packages /opt/prefect/custom-packages

# Copier ton code Prefect (flows et prefect.yaml)
COPY flows/ /opt/prefect/flows/

# Ajouter le dossier custom au PYTHONPATH
ENV PYTHONPATH=/opt/prefect/custom-packages:$PYTHONPATH

# TODO : retester sans forcer PREFECT_API_URL
ENV PREFECT_API_URL=http://prefect-server.data.svc.cluster.local:4200/api

# Entrypoint pour forcer PREFECT_API_URL
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
