import subprocess

def run_dbt(command):
    subprocess.run(["dbt", command], cwd="/dbt", check=True)

