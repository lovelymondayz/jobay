# Go Docker Build: Pre-built Binary Pattern for Slow VPS

## Problem

`docker compose build --no-cache` with multi-stage Go builds takes **10-15+ minutes** on a 2-core VPS because:
- `go mod download` re-downloads all dependencies
- `go build` compiles everything from scratch with CGO disabled
- First-time terminal tool timeouts kill the process, discarding cache

## Solution: Pre-built Binary → Single-Stage Docker

Build the binary **on host** (fast, uses local module cache), then copy into a minimal Alpine image (17 seconds).

### Step 1: Build on host
```bash
cd backend
go build -o server ./cmd/server
```

### Step 2: Single-stage Dockerfile
```dockerfile
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY server .
COPY migrations/ ./migrations/
EXPOSE 8083
CMD ["./server"]
```

### Step 3: Fast Docker build
```bash
cd /root/hermes/<project>
docker compose build backend   # 17 seconds, not 15 minutes
docker compose up -d
```

## Trade-offs

| | Multi-stage (from scratch) | Pre-built binary |
|---|---|---|
| Build time | 10-15 min | ~17 sec |
| Reproducibility | ✅ Fully reproducible | ⚠️ Host-compiled binary |
| CI/CD suitability | ✅ Perfect | ❌ Requires build step |
| VPS dev/test | ❌ Too slow | ✅ Fast iteration |

**Recommendation:** Use pre-built binary for VPS dev/test. Use multi-stage for CI/CD or production releases.

## Why This Works

- `go build` on host uses the local module cache (`~/go/pkg/mod`) — no re-download
- Docker `COPY server .` is just a file copy — no compilation
- Alpine image is already cached after first pull

## Session Origin

Klikku photobooth (2026-08-20) — Multi-stage build hit 13+ minute timeout twice. Pre-built binary built in ~5 min on host, then Docker image built in 17 sec.

---

# Go Module Path: Local vs GitHub Remote

## Problem

Module declared as `github.com/lovelymondayz/klikku` causes `go mod download` in Docker to fail with:
```
repository 'https://github.com/lovelymondayz/klikku/' not found
```

Because the repo doesn't exist on GitHub yet (local-only project).

## Solution

Use a **local module name** in `go.mod`:
```
module klikku
```

Then imports become:
```go
import "klikku/internal/config"
import "klikku/internal/handlers"
```

## Trade-offs

| | `github.com/...` | `klikku` (local) |
|---|---|---|
| Docker build | ❌ Needs repo to exist | ✅ Always works |
| CI/CD push | ✅ Standard | ⚠️ Must rename before open-sourcing |
| go mod download | ❌ Tries remote | ✅ Local only |

**Recommendation:** Start local. Rename to GitHub path when pushing to remote.

---

# Fiber v2.52 Middleware Package Paths

## Breaking Change

In Fiber v2.52.x, some middleware packages were renamed:

| Old (pre-v2.52) | New (v2.52+) |
|---|---|
| `middleware/ratelimit` | `middleware/limiter` |
| `ratelimit.New()` | `limiter.New()` |
| `ratelimit.Config` | `limiter.Config` |

`cors` and `logger` paths remain unchanged.

## Fix

```go
// WRONG (v2.52+)
import "github.com/gofiber/fiber/v2/middleware/ratelimit"
return ratelimit.New(ratelimit.Config{...})

// RIGHT
import "github.com/gofiber/fiber/v2/middleware/limiter"
return limiter.New(limiter.Config{...})
```

## Session Origin

Klikku photobooth (2026-08-20) — `go build` failed with `undefined: ratelimit` after `go mod tidy` resolved Fiber to v2.52.5.
