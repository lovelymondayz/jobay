"""
Pipeline configuration — all settings in one place.
"""
import os

def _load_env():
    """Load .env file if it exists."""
    env_path = os.path.join(os.path.dirname(__file__), ".env")
    if os.path.exists(env_path):
        for line in open(env_path).read().strip().split(chr(10)):
            line = line.strip()
            if line and not line.startswith(chr(35)) and chr(61) in line:
                k, v = line.split(chr(61), 1)
                os.environ.setdefault(k.strip(), v.strip())

_load_env()

def _read_file(path):
    if os.path.exists(path):
        return open(path).read().strip()
    return ""

def _get(key, default=""):
    v = os.environ.get(key, "")
    return v if v else default

# PostgreSQL
PG_HOST = _get("PG_HOST", "localhost")
PG_PORT = int(_get("PG_PORT", "5433"))
PG_DB = _get("PG_DB", "leadgen")
PG_USER = _get("PG_USER", "leadgen")
PG_PASSWORD = _get("PG_PASSWORD", "")

# Firecrawl
def _get_fc_key():
    k = os.environ.get("FIRECRAWL_API_KEY", "")
    if k:
        return k
    p = "/root/.firecrawl/.env"
    if os.path.exists(p):
        for line in open(p).read().strip().split(chr(10)):
            if line.startswith("FIRECRAWL_API_KEY" + chr(61)):
                return line.split(chr(61), 1)[1].strip()
    return ""

FIRECRAWL_API_KEY = _get_fc_key()
FIRECRAWL_BASE = "https://api.firecrawl.dev/v2"

# GitHub
GITHUB_USER = "lovelymondayz"
GITHUB_PAT = _read_file("/root/.github/pat")
GITHUB_API = "https://api.github.com"

# Cloudflare (DNS managed manually via Cloudflare dashboard)
# Pages auto-deploys from GitHub — no API needed
CLIENT_BASE_DOMAIN = "client.arjism.com"

# Email SMTP — Gmail
# Use App Password (not regular password): https://myaccount.google.com/apppasswords
SMTP_HOST = _get("SMTP_HOST", "smtp.gmail.com")
SMTP_PORT = int(_get("SMTP_PORT", "587"))
SMTP_USER = os.environ.get("SMTP_USER", "arjithedev@gmail.com")
SMTP_PASSWORD = _get("SMTP_PASSWORD", "")
SMTP_FROM_NAME = _get("SMTP_FROM_NAME", "Arjism Web Studio")
PING_EMAIL_TO = _get("PING_EMAIL_TO", "arjithedev@gmail.com")

# Discord
DISCORD_WEBHOOK = os.environ.get("DISCORD_WEBHOOK", "")

# Scouting
TARGET_CITIES = ["Jakarta", "BSD", "Tangerang", "Bogor", "Depok"]
TARGET_CATEGORIES = [
    "restaurant", "gym", "dental clinic", "beauty salon",
    "barber shop", "cafe", "bakery", "pet clinic", "laundry", "tailor",
]
MIN_RATING = 4.0
MIN_REVIEWS = 20
DAILY_BUILD_LIMIT = 10
DAILY_OUTREACH_LIMIT = 15

# Paths
BASE_DIR = "/root/hermes/leadgen-localbiz"
TEMPLATES_DIR = BASE_DIR + "/templates"
ASSETS_DIR = BASE_DIR + "/assets"
LOGS_DIR = BASE_DIR + "/logs"
