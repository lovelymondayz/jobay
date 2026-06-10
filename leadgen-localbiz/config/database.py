"""
PostgreSQL connection for Lead Gen Pipeline.
Docker: leadgen-postgres on port 5433
"""
import os
import psycopg2
from psycopg2.extras import RealDictCursor
from contextlib import contextmanager

# Auto-load .env so DB_PASSWORD is always available
_env_path = os.path.join(os.path.dirname(__file__), ".env")
if os.path.exists(_env_path):
    for _line in open(_env_path).read().strip().split("\n"):
        _line = _line.strip()
        if _line and not _line.startswith("#") and "=" in _line:
            _k, _v = _line.split("=", 1)
            os.environ.setdefault(_k.strip(), _v.strip())

DB_CONFIG = {
    "host": os.getenv("PG_HOST", "localhost"),
    "port": int(os.getenv("PG_PORT", 5433)),
    "dbname": os.getenv("PG_DB", "leadgen"),
    "user": os.getenv("PG_USER", "leadgen"),
    "password": os.getenv("PG_PASSWORD", ""),
}

@contextmanager
def get_conn():
    conn = psycopg2.connect(**DB_CONFIG)
    try:
        yield conn
        conn.commit()
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()

@contextmanager
def get_cursor(cursor_factory=RealDictCursor):
    with get_conn() as conn:
        cur = conn.cursor(cursor_factory=cursor_factory)
        try:
            yield cur
        finally:
            cur.close()

def query(sql, params=None, fetch="all"):
    with get_cursor() as cur:
        cur.execute(sql, params)
        if fetch == "all":
            return cur.fetchall()
        elif fetch == "one":
            return cur.fetchone()
        return None

def execute(sql, params=None):
    with get_conn() as conn:
        cur = conn.cursor()
        cur.execute(sql, params)
        cur.close()
