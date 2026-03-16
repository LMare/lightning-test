from prefect import task
from minio import Minio
import requests
import io

from . import utils as u


# -----------------------------
# 1) Export Postgres → Parquet → MinIO
# -----------------------------

from typing import Dict

@task
def export_postgres_to_minio(table: str, params: Dict[str, str], filename: str, time_field: str) -> dict:
    """
    Exporte un chunk de table Postgres en Parquet, le stocke sur MinIO, retourne un dict (url, max_time, row_count).
    """
    result = u.export_table_postgres_to_parquet_buffer(table, params, time_field)
    url = u.store_buffer_in_minio(result["buffer"], filename)
    return {"url": url, "max_time": result["max_time"], "row_count": result["row_count"]}


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


