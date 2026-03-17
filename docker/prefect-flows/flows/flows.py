from datetime import datetime

from prefect import flow
import io

from flows import config as cfg, tasks as t, utils as u

#TODO relire & check chunking params....

# Classe runner globale pour l'orchestration de toutes les tables
class ELTFlowRunner:
    def __init__(self):
        pass

    def read_state(self):
        try:
            return u.read_json_from_storage(cfg.MINIO_BUCKET, "state.json")
        except Exception:
            # TODO Log
            return {"tables": {}}

    def write_state(self, state: dict):
        import json
        client = u.client_minio()
        data = json.dumps(state).encode("utf-8")
        client.put_object(
            bucket_name=cfg.MINIO_BUCKET,
            object_name="state.json",
            data=io.BytesIO(data),
            length=len(data),
            content_type="application/json",
        )

    def run_flow(self):
        state = self.read_state()
        current_datetime = datetime.now().isoformat()
        status = {}
        for table in u.read_tables_config().keys():
            state_table = state.get("tables", {}).get(table, {})
            status[table] = elt_table_flow(table, current_datetime, state_table)
        status["dbt"] = dbt_flow(cfg.DUCKDB_SERVICE)
        self.write_state(state)
        return status

# Runner par table
class ELTTableFlowRunner:
    def __init__(self, table: str, current_datetime: str, state_table: dict):
        self.table = table
        self.current_datetime = current_datetime
        self.state_table = state_table
        self.table_cfg = u.read_tables_config().get(table)
        self.chunk_size = self.table_cfg.get("chunk_size", 10000)
        self.time_field = self.table_cfg.get("time_field", "updated_at")
        self.last_export = state_table.get("last_export")
        self.last_load = state_table.get("last_load")

    def run_table(self):
        parquet_files, export_timestamps = self._export_chunks()
        loaded_files = self._load_files(parquet_files)
        new_last_export = export_timestamps[-1] if export_timestamps else self.last_export
        new_last_load = max([f["timestamp"] for f in loaded_files], default=self.last_load)
        missing_loads = [
            {"filename": f[0], "timestamp": self.current_datetime}
            for f in parquet_files if self.current_datetime > (self.last_load or "")
        ] if not loaded_files else []
        return {
            "parquet_files": [f[0] for f in parquet_files],
            "loaded_files": loaded_files,
            "last_export": new_last_export,
            "last_load": new_last_load,
            "missing_loads": missing_loads,
        }

    def _export_chunks(self):
        export_timestamps = []
        chunk_start = self.last_export
        chunk_end = self.current_datetime
        chunk_idx = 0
        parquet_files = []
        while True:
            params = {"last_run": chunk_start, "chunk_end": chunk_end, "chunk_size": self.chunk_size}
            filename = f"{self.table}/{self.current_datetime}_chunk{chunk_idx}.parquet"
            export_result = t.export_postgres_to_minio.fn(self.table, params, filename, self.time_field)
            parquet_files.append((filename, export_result["url"]))
            export_timestamps.append(self.current_datetime)
            if export_result["row_count"] < self.chunk_size or not export_result["max_time"]:
                break
            chunk_start = export_result["max_time"]
            chunk_idx += 1
        return parquet_files, export_timestamps

    def _load_files(self, parquet_files):
        loaded_files = []
        missing_loads = self.state_table.get("missing_loads", [])
        for missing in missing_loads:
            fname = missing["filename"]
            file_ts = missing.get("timestamp", self.current_datetime)
            load_result = t.trigger_duckdb_load.fn(
                cfg.DUCKDB_SERVICE,
                source_path=cfg.OUTPUT_PATH + fname,
                table_name=self.table,
            )
            loaded_files.append({"file": fname, "result": load_result, "timestamp": file_ts})
        for fname, parquet_url in parquet_files:
            file_ts = self.current_datetime
            if file_ts > (self.last_load or ""):
                load_result = t.trigger_duckdb_load.fn(
                    cfg.DUCKDB_SERVICE,
                    source_path=cfg.OUTPUT_PATH + fname,
                    table_name=self.table,
                )
                loaded_files.append({"file": fname, "result": load_result, "timestamp": file_ts})
        return loaded_files

# Runner pour le flow dbt
class DBTFlowRunner:
    def __init__(self, service_url: str):
        self.service_url = service_url

    def run(self):
        return t.trigger_duckdb_dbt.fn(self.service_url)



# -----------------------------
# PRINCIPAL FLOW
# -----------------------------
# export all the table from postgres to minio & trigger duckdb load 
# trigger the run of dbt
@flow(name="elt_master_flow")
def elt_master_flow():
    runner = ELTFlowRunner()
    return runner.run_flow()


@flow(name="elt_table_flow")
def elt_table_flow(table: str, current_datetime: str, state_table: dict):
    runner = ELTTableFlowRunner(table, current_datetime, state_table)
    return runner.run_table()


@flow(name="dbt_flow")
def dbt_flow(service_url: str):
    runner = DBTFlowRunner(service_url)
    return runner.run()


if __name__ == "__main__":
    elt_master_flow()

