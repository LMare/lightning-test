FROM python:3.11-slim AS builder

# Install system deps
RUN apt-get update \
 && apt-get install -y build-essential \
 && rm -rf /var/lib/apt/lists/*


# Install Python deps
COPY requirements.txt .
RUN pip install --prefix=/install -r requirements.txt --no-cache-dir


FROM prefecthq/prefect:2-latest

# Installer tes dépendances Python
COPY --from=builder /install /usr/local

# Copier ton code Prefect
COPY flows/ /opt/prefect/flows/

# Point d’entrée par défaut (Prefect gère l’exécution)
CMD ["python", "-m", "prefect.engine"]
