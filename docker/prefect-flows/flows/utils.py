from minio import Minio
import io
import psycopg

from flows import config as cfg

import yaml
import json


import pandas as pd
import pyarrow as pa
import pyarrow.parquet as pq


# read the tables config from yaml
def read_tables_config():
    return read_yaml_file("tables.yaml")["tables"]

def read_yaml_file(path: str):
    with open(path, "r") as f:
        return yaml.safe_load(f)

def client_minio():
    return Minio(
        cfg.MINIO_ENDPOINT,
        access_key=cfg.MINIO_ACCESS_KEY,
        secret_key=cfg.MINIO_SECRET_KEY,
        secure=False,
    )



def read_json_from_storage(bucket: str, object_name: str):
    client = client_minio()
    response = client.get_object(bucket, object_name)
    return json.loads(response.read().decode("utf-8"))



def export_table_postgres_to_parquet_buffer(table: str, params: dict, time_field: str) -> dict:
    """
    Exporte un chunk de table Postgres en Parquet (fenêtre temporelle via params)
    Retourne un dict: {buffer, max_time, row_count}
    """
    table_cfgs = read_tables_config()
    table_cfg = table_cfgs.get(table)
    if not table_cfg:
        raise ValueError(f"Table config not found for {table}")
    query = table_cfg.get("query")
    # Remplacer les paramètres dans la requête
    for k, v in params.items():
        query = query.replace(f":{k}", f"'{v}'")
    with psycopg.connect(cfg.PG_CONN) as conn:
        df = pd.read_sql(query, conn)
    table_arrow = pa.Table.from_pandas(df)
    buffer = io.BytesIO()
    pq.write_table(table_arrow, buffer)
    buffer.seek(0)
    max_time = None
    row_count = len(df)
    if row_count > 0 and time_field in df.columns:
        max_time = str(df[time_field].max())
    return {"buffer": buffer, "max_time": max_time, "row_count": row_count}

def store_buffer_in_minio(buffer: io.BytesIO, output_path: str) -> str:
    client = client_minio()
    client.put_object(
        bucket_name=cfg.MINIO_BUCKET,
        object_name=output_path,
        data=buffer,
        length=len(buffer.getvalue()),
        content_type="application/octet-stream",
    )
    return f"s3://{cfg.MINIO_BUCKET}/{output_path}"