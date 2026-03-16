import os


MINIO_ENDPOINT = os.getenv("MINIO_ENDPOINT", "minio.data.svc.cluster.local:9000")
MINIO_BUCKET = os.getenv("MINIO_BUCKET", "raw")
MINIO_ACCESS_KEY = os.getenv("MINIO_ACCESS_KEY", "minio")
MINIO_SECRET_KEY = os.getenv("MINIO_SECRET_KEY", "minio123")
OUTPUT_PATH = os.getenv("OUTPUT_PATH", "exports/")
DUCKDB_SERVICE = os.getenv("DUCKDB_SERVICE", "http://duckdb-elt-service.data.svc.cluster.local")
PG_CONN = os.getenv("PG_CONN", "postgresql://user:pass@postgres.data.svc.cluster.local:5432/db")