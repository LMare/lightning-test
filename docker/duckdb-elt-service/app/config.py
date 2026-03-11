import os

MINIO_ENDPOINT = os.getenv("MINIO_ENDPOINT")
MINIO_ACCESS_KEY = os.getenv("MINIO_ACCESS_KEY")
MINIO_SECRET_KEY = os.getenv("MINIO_SECRET_KEY")
MINIO_SSL = os.getenv("MINIO_SSL", "false")
DUCKDB_PATH = os.getenv("DUCKDB_PATH", "/data/warehouse.duckdb")

