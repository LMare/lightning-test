FROM python:3.11-slim
RUN pip install dbt-core dbt-duckdb
WORKDIR /dbt
COPY . .
CMD ["dbt", "run"]

