from pydantic import BaseModel

class LoadRequest(BaseModel):
    source_path: str
    table_name: str

