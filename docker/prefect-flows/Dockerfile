
FROM python:3.11-slim AS builder

# Install system deps (libpq-dev for psycopg, build-essential for pip)
RUN apt-get update \
 && apt-get install -y build-essential libpq-dev \
 && rm -rf /var/lib/apt/lists/*

# Install Python deps in a custom folder
COPY requirements.txt ./
RUN pip install --target=/custom-packages -r requirements.txt --no-cache-dir

FROM prefecthq/prefect:2-latest

# Copier les dépendances Python sans écraser /usr/local
COPY --from=builder /custom-packages /opt/prefect/custom-packages

# Copier ton code Prefect (flows et prefect.yaml)
COPY flows/ /opt/prefect/flows/
COPY prefect.yaml /opt/prefect/prefect.yaml

# Ajouter le dossier custom au PYTHONPATH
ENV PYTHONPATH=/opt/prefect/custom-packages:$PYTHONPATH

