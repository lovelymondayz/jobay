# 2026-06-15: VPS Storage Analysis and Cleanup

## Sessions
- **20260615_095918_01a98b70** — "Greeting and Introduction" (Discord, 22 messages, ~09:59–10:12 UTC)

## What Was Built / Done

### VPS Storage Analysis & Cleanup
- kvinn asked why VPS storage was large and requested an ROI breakdown to decide what to delete
- Ran comprehensive storage audit: `df -h`, `du -sh /*`, `docker system df`, `docker images`, journal analysis, `/root/.hermes/` breakdown
- Found 46GB total, 32GB used (73%)

### Storage Breakdown
| Category | Size | Notes |
|---|---|---|
| Docker images (10 active) | 9.26 GB | All in use, 0% reclaimable |
| Docker rootfs (containers) | 9.4 GB | Container writable layers |
| System journal logs | 3.7 GB | /var/log/journal |
| /root/.hermes/lsp | 106 MB | Language server cache |
| /root/.hermes/state.db | 40 MB | Hermes state DB |
| /root/hermes/digital-yearbook | 421 MB | Project source |
| /root/hermes/wedding-invitation | 243 MB | Project source |
| Old backup zips (3x) | 49 MB | Jun 11-13 |
| Go tarball in /tmp | 66 MB | Leftover installer |

### Cleanup Executed (kvinn approved "delete all safe, don't touch Immich")
| Action | Freed |
|---|---|
| Deleted 3 old backup zips | 49 MB |
| Deleted Go tarball from /tmp | 66 MB |
| Trimmed journal logs (3.7GB → 500MB cap) | ~3.1 GB |
| Pruned unused Docker volumes | 0 B (already clean) |
| **Total freed** | **~3.2 GB** |

### Result
- **Before:** 32 GB used (73%)
- **After:** 29 GB used (66%)
- **15 GB free** (up from 12 GB)

## Running Services (all healthy)
- unifi-controller, yearbook-frontend/backend, wedding-frontend/backend, leadgen-postgres, immich_server, immich_microservices, immich_redis, immich_postgres, immich_machine_learning

## Key Learnings
1. **Journal logs were the biggest win.** 3.7 GB of systemd journal was the single largest reclaimable item. `journalctl --vacuum-size=500M` freed 3.1 GB instantly.
2. **Docker images were all honest.** Every image had a running container — no orphaned images to clean.
3. **Immich is the storage heavyweight.** immich-server (3.09 GB) + immich-machine-learning (1.9 GB) + pgvecto-rs (927 MB) = ~5.9 GB for one service family. Worth knowing for future capacity planning.
4. **Old backup zips accumulate quietly.** 3 days of zips = 49 MB. Should add a retention policy (keep only latest, or move off-disk).

## Mistakes Overcome
- First `du -sh /*` command timed out (15s) — likely due to deep directory traversal. Subsequent targeted commands worked fine.
- `docker volume prune` reclaimed 0 B — volumes were already clean, but worth checking.

## Tech Stack
- Proxmox VE host (LVM-thin: `/dev/mapper/pve-vm--101--disk--0`)
- Docker (13 containers, 10 images)
- systemd journal
- Hermes Agent on VPS (2 core, 12GB RAM)
