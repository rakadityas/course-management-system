import os
import re
import mysql.connector

class DB:
    def __init__(self, dsn: str):
        self.dsn = dsn

    def _parse(self):
        m = re.match(r"(?P<user>[^:]+):(?P<pw>[^@]+)@tcp\((?P<host>[^:]+):(?P<port>\d+)\)/(?P<db>[^?]+)", self.dsn)
        if m:
            return {
                "user": m.group("user"),
                "password": m.group("pw"),
                "host": m.group("host"),
                "port": int(m.group("port")),
                "database": m.group("db"),
            }
        url = self.dsn
        if url.startswith("mysql://"):
            url = url[len("mysql://"):]
        parts = re.match(r"(?P<user>[^:]+):(?P<pw>[^@]+)@(?P<host>[^:]+):(?P<port>\d+)/(?P<db>[^?]+)", url)
        if parts:
            return {
                "user": parts.group("user"),
                "password": parts.group("pw"),
                "host": parts.group("host"),
                "port": int(parts.group("port")),
                "database": parts.group("db"),
            }
        raise ValueError("Unsupported DATABASE_URL format")

    def get_connection(self):
        cfg = self._parse()
        return mysql.connector.connect(
            user=cfg["user"], password=cfg["password"], host=cfg["host"], port=cfg["port"], database=cfg["database"]
        )

def init_db_from_env() -> 'DB':
    dsn = os.getenv("DATABASE_URL", "")
    if not dsn:
        raise RuntimeError("DATABASE_URL is not set")
    return DB(dsn)
