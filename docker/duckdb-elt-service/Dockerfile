FROM python:3.11-slim AS builder

# Install system deps
RUN apt-get update \
 && apt-get install -y build-essential \
 && rm -rf /var/lib/apt/lists/*


# Install Python deps
COPY requirements.txt .
RUN pip install --prefix=/install -r requirements.txt --no-cache-dir

# Copy app + dbt project
FROM python:3.11-slim
WORKDIR /service

COPY --from=builder /install /usr/local
COPY app ./app
COPY dbt ./dbt

EXPOSE 8000
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]

