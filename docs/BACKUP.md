# Hermes Backup & Restore

Every Docker project on this VPS, backed up nightly and pushed off-site.

## What is protected, and what is not

| Asset | Where it lives | Recoverable? |
|---|---|---|
| Source code | Per-project GitHub repos | ✅ `git clone` |
| **Databases** | Postgres volumes | ✅ nightly `pg_dump` → local + GitHub |
| **`.env` secrets** | Gitignored on purpose | ✅ captured in `config/` |
| `docker-compose.yml` | Project repos | ✅ also captured for convenience |
| Immich photos (2.6GB) | `/home/root/immich-data` | ⚠️ **see Immich section** |
| Docker images | Local only | ❌ rebuilt from source |

The gap this closes: **code was already safe, data was not.** Losing the
Postgres volumes would lose every booking, member and lead — and no `git clone`
brings those back.

## Schedule

```
0  2 * * *   qlio-platform/scripts/backup.sh     (qlio's own, pre-existing)
0  3 * * *   hermes/scripts/backup-all.sh        all databases + config
30 3 * * *   hermes/scripts/backup-offsite.sh    push newest to GitHub
```

Retention: **14 days** locally (`/root/hermes-backups`), **7 backups** in the
private `lovelymondayz/BackupHermes` repo.

## Coverage

| Project | Container | DB | Method |
|---|---|---|---|
| wedding-invitation | `wedding-db` | wedding | pg_dump |
| digital-yearbook | `yearbook-db` | yearbook | pg_dump |
| members | `members-db` | members | pg_dump |
| noidk | `noidk-db` | noidk | pg_dump |
| qlio-platform | `qlio-db` | qlio | pg_dump |
| say-less | `sayless-db` | sayless | pg_dump (**user is `postgres`**, not `sayless`) |
| pico | `pico-postgres` | pico | pg_dump |
| leadgen-localbiz | *(no container)* | leadgen_pgdata | **cold tar** — cluster is stopped |
| mentengdutch | *(static site)* | — | config only |

## Daily use

```bash
/root/hermes/scripts/backup-all.sh            # run now
/root/hermes/scripts/backup-all.sh --list     # what exists
/root/hermes/scripts/restore-one.sh --list    # projects + available backups
/root/hermes/scripts/restore-one.sh qlio      # restore newest
/root/hermes/scripts/restore-one.sh qlio 20260819_031510
```

## Safety design

`backup-all.sh` refuses to keep a bad dump:

| Guard | Behaviour |
|---|---|
| Container not running | Skips, logs it — never writes a fake file |
| `pg_dump` non-zero | Deletes the partial file |
| Dump under 300 bytes | Deletes it — means wrong user/db |
| `gzip -t` fails | Deletes it |
| Any failure | Exits **non-zero** so cron surfaces it |

`restore-one.sh` takes a **safety dump of current state first**, requires typing
`RESTORE`, then recreates the app container so it gets a clean connection pool.

`backup-offsite.sh` **checks the repo is private before pushing** and aborts if
not — the payload contains `.env` secrets and customer PII.

## Verified

Not assumed — tested by destroying real data:

```
noidk: 14 tables, 17 rows in `contributions`
DROP SCHEMA public CASCADE     → 0 tables
restore from backup            → 14 tables, 17 rows, 0 psql errors
app health                     → HTTP 200
```

Off-site push verified by listing the files back through the GitHub API.

## Immich (handled separately)

Immich is **2.6GB of photos** — too large for a git repo, and it needs its own
strategy. Its database is small and is worth dumping; the photo library needs a
real off-site target.

```bash
# database only (259MB raw, ~98MB gzipped)
tar -czf /root/immich-backup/immich-db-$(date +%F).tar.gz \
  -C /home/root/immich-data db
```

⚠️ **The photo library has no off-site copy.** For real durability point rclone
at any cloud provider:

```bash
rclone copy /home/root/immich-data/upload remote:immich-photos
```

Immich image tags are **pinned to v2.5.6** — never let them float to `:release`,
because Immich runs irreversible schema migrations on boot.

## Restore-from-scratch (total VPS loss)

1. Provision a host with Docker + Compose v2.
2. `git clone https://github.com/lovelymondayz/BackupHermes` (needs the PAT).
3. For each project: clone its own repo, drop the matching `config/<project>/.env`
   back into place.
4. `docker compose build && docker compose up -d` per project — sequentially, not
   in parallel; on 2 cores parallel builds push load past 15.
5. Load each database:
   ```bash
   zcat databases/<stamp>/databases/<name>.sql.gz \
     | docker exec -i <container> psql -U <user> -d <db> -v ON_ERROR_STOP=0
   ```
6. Re-add the Cloudflare Tunnel public hostnames (`http://localhost:<port>` —
   scheme **http**, keep the port).
7. Reinstall the cron lines from the Schedule section above.

## Off-site gap to close

Backups live on the same disk as the databases and in one GitHub repo. That
survives a bad migration, an accidental prune, and disk-level corruption — but
**not** losing the GitHub account. A third copy (rclone to cloud storage) is the
remaining improvement.
