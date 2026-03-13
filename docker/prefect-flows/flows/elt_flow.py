from prefect import flow, task
import pandas as pd
import psycopg
from minio import Minio
import requests
import pyarrow as pa
import pyarrow.parquet as pq
import io


# -----------------------------
# 1) Export Postgres → Parquet → MinIO
# -----------------------------
@task
def export_postgres_to_minio(
    pg_conn: str,
    query: str,
    minio_endpoint: str,
    minio_bucket: str,
    minio_access_key: str,
    minio_secret_key: str,
    output_path: str,
):
    # 1. Lire Postgres
    with psycopg.connect(pg_conn) as conn:
        df = pd.read_sql(query, conn)

    # 2. Convertir en Parquet (en mémoire)
    table = pa.Table.from_pandas(df)
    buffer = io.BytesIO()
    pq.write_table(table, buffer)
    buffer.seek(0)

    # 3. Upload MinIO
    client = Minio(
        minio_endpoint,
        access_key=minio_access_key,
        secret_key=minio_secret_key,
        secure=False,
    )

    client.put_object(
        bucket_name=minio_bucket,
        object_name=output_path,
        data=buffer,
        length=len(buffer.getvalue()),
        content_type="application/octet-stream",
    )

    return f"s3://{minio_bucket}/{output_path}"


# -----------------------------
# 2) Appeler /load du service DuckDB
# -----------------------------
@task
def trigger_duckdb_load(service_url: str, source_path: str, table_name: str):
    payload = {
        "source_path": source_path,
        "table_name": table_name,
    }
    r = requests.post(f"{service_url}/load", json=payload)
    r.raise_for_status()
    return r.json()


# -----------------------------
# 3) Appeler /dbt/run du service DuckDB
# -----------------------------
@task
def trigger_duckdb_dbt(service_url: str):
    r = requests.post(f"{service_url}/dbt/run")
    r.raise_for_status()
    return r.json()


# -----------------------------
# FLOW PRINCIPAL
# -----------------------------
@flow(name="postgres_to_duckdb_pipeline")
def main_flow():
    # -------------------------
    # Paramètres (à mettre en env vars dans ton Deployment)
    # -------------------------
    PG_CONN = "postgresql://user:pass@postgres.data.svc.cluster.local:5432/db"
    QUERY = "SELECT * FROM my_table"

    MINIO_ENDPOINT = "minio.data.svc.cluster.local:9000"
    MINIO_BUCKET = "raw"
    MINIO_ACCESS_KEY = "minio"
    MINIO_SECRET_KEY = "minio123"
    OUTPUT_PATH = "exports/my_table.parquet"

    DUCKDB_SERVICE = "http://duckdb-elt-service.data.svc.cluster.local"

    # -------------------------
    # Pipeline
    # -------------------------
    parquet_path = export_postgres_to_minio(
        PG_CONN,
        QUERY,
        MINIO_ENDPOINT,
        MINIO_BUCKET,
        MINIO_ACCESS_KEY,
        MINIO_SECRET_KEY,
        OUTPUT_PATH,
    )

    load_result = trigger_duckdb_load(
        DUCKDB_SERVICE,
        source_path=OUTPUT_PATH,
        table_name="my_table",
    )

    dbt_result = trigger_duckdb_dbt(DUCKDB_SERVICE)

    return {
        "parquet": parquet_path,
        "load": load_result,
        "dbt": dbt_result,
    }


if __name__ == "__main__":
    main_flow()

