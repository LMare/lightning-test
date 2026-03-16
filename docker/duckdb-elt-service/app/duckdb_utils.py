import duckdb

from .config import (
    MINIO_ENDPOINT,
    MINIO_ACCESS_KEY,
    MINIO_SECRET_KEY,
    MINIO_SSL,
    DUCKDB_PATH,
)


def run_load(source_path: str, table_name: str):
    con = duckdb.connect(f"{DUCKDB_PATH}")
    con.execute(f"SET s3_endpoint='{MINIO_ENDPOINT}'")
    con.execute(f"SET s3_access_key_id='{MINIO_ACCESS_KEY}'")
    con.execute(f"SET s3_secret_access_key='{MINIO_SECRET_KEY}'")
    con.execute(f"SET s3_use_ssl={MINIO_SSL}")

    # On suppose que la table a un champ updated_at
    con.execute(f"""
        MERGE INTO {table_name} AS target
        USING read_parquet('{source_path}') AS source
        ON target.id = source.id
        WHEN MATCHED AND source.updated_at > target.updated_at THEN
            UPDATE SET *
        WHEN NOT MATCHED THEN
            INSERT *
    """)

    con.close()
