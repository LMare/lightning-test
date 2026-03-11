from fastapi import FastAPI
from .schemas import LoadRequest
from .duckdb_utils import run_load
from .dbt_utils import run_dbt

app = FastAPI()

@app.get("/health")
def health():
    return {"status": "ok"}

@app.post("/load")
def load(req: LoadRequest):
    run_load(req.source_path, req.table_name)
    return {"status": "load done"}

@app.post("/dbt/run")
def dbt_run():
    run_dbt("run")
    return {"status": "dbt run done"}

@app.post("/dbt/test")
def dbt_test():
    run_dbt("test")
    return {"status": "dbt test done"}

