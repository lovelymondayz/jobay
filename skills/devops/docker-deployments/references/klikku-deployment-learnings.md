# Klikku Deployment Learnings

## Prebuilt Binary Pattern for Slow VPS

On slow VPS (2c/12GB), multi-stage Go builds inside Docker take 5+ min and often timeout. **Prefer the prebuilt binary pattern:** build on host (~5 min once), then single-stage Docker copies the binary (~17s).

```dockerfile
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY server .
COPY migrations/ ./migrations/
EXPOSE 8083
CMD ["./server"]
```

```bash
cd backend && CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
```

## Go Binary: Alpine (musc) vs Debian (glibc)

**Symptom:** Container starts but immediately crashes: `exec ./server: no such file or directory`

**Cause:** Go binary built on Ubuntu/Debian host is dynamically linked against glibc, but Alpine uses musl libc. Even with `CGO_ENABLED=0`, the binary may still be dynamically linked.

**Fix:** Use `debian:bookworm-slim` instead of `alpine`:

```dockerfile
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY server .
COPY migrations/ ./migrations/
EXPOSE 8083
CMD ["./server"]
```

**Verify:** `file server` should show "statically linked" for Alpine, or "dynamically linked" for Debian (which is fine on Debian).

## MinIO Requires x86-64-v2 CPU

**Symptom:** MinIO container crashes: `Fatal glibc error: CPU does not support x86-64-v2`

**Fix:** Replace MinIO with filesystem-based storage (bind-mounted volume). Use a simple Storage wrapper:

```go
type Storage struct { basePath string }
func (s *Storage) Upload(bucket, objectName string, data []byte, contentType string) error {
    dir := filepath.Join(s.basePath, bucket)
    os.MkdirAll(dir, 0755)
    return os.WriteFile(filepath.Join(dir, objectName), data, 0644)
}
```

Verify CPU: `grep avx2 /proc/cpuinfo` — empty means MinIO won't work.

## nginx proxy_pass Trailing Slash

**Symptom:** API calls return 404 with "Cannot POST /auth/login" (missing `/api/` prefix).

**Cause:** `proxy_pass http://backend:8083/;` (WITH trailing slash) strips the `/api/` prefix before forwarding. So `/api/auth/login` becomes `/auth/login` at the backend.

**Fix:** Remove trailing slash from `proxy_pass`:

```nginx
location /api/ {
    proxy_pass http://backend:8083;  # NO trailing slash — preserves full path
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

## Tailwind CSS Not Generating

**Symptom:** Frontend loads but no styles (CSS < 10KB).

**Causes:**
1. Missing `postcss.config.js` (or `postcss.config.cjs` when `"type": "module"` in package.json)
2. CSS variables in `tailwind.config.js` break opacity modifiers (`bg-primary/50` fails). Use direct hex values.
3. `@import` must precede `@tailwind` directives.

**Debug:** Built CSS < 10KB = Tailwind not processing.

## golang-migrate pgx Driver Dirty State

**Symptom:** `migrate: Dirty database version 2. Fix and force version.`

**Cause:** The golang-migrate `pgx` driver can get stuck on dirty migrations, especially when migration files have `-- +migrate Up` / `-- +migrate Down` comments that confuse the parser.

**Fix:** Replace with a simple custom migration runner:

```go
func RunMigrations(pool *pgxpool.Pool, dbName string) error {
    db := stdlib.OpenDBFromPool(pool)
    defer db.Close()
    
    // Create tracking table
    db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
        version INT PRIMARY KEY,
        applied_at TIMESTAMPTZ DEFAULT NOW()
    )`)
    
    // Read .sql files, check if applied, execute in transaction
    // Record migration in schema_migrations
}
```

## Super Admin Seed Ordering

**Symptom:** `Warning: Could not seed super admin: no rows in result set`

**Cause:** The `SeedSuperAdmin` function tries to create a user with `merchant_id` FK, but no merchants exist yet.

**Fix:** Create a default merchant FIRST if none exist:

```go
func SeedSuperAdmin(pool *pgxpool.Pool, cfg *config.Config) error {
    // ... check if user exists ...
    
    // Get or create default merchant
    var merchantID string
    err := pool.QueryRow(ctx, "SELECT id FROM merchants LIMIT 1").Scan(&merchantID)
    if err == sql.ErrNoRows {
        // Create default merchant first
        err = pool.QueryRow(ctx, 
            "INSERT INTO merchants (business_name, slug, subscription_status) VALUES ($1, $2, $3) RETURNING id",
            "Default Merchant", "default", "active").Scan(&merchantID)
        if err != nil {
            return fmt.Errorf("create default merchant: %w", err)
        }
    }
    
    // Now create super admin with valid merchant_id
    // ...
}
```

## Response Format Mismatch (Backend → Frontend)

**Symptom:** Login returns 200 OK but frontend stays on login page (no redirect).

**Cause:** Backend wraps responses as `{success: true, data: {access_token, ...}}` but frontend reads `res.data.access_token` instead of `res.data.data.access_token`.

**Fix:** In the Zustand store, unwrap the response:

```typescript
login: async (email, password) => {
    const res = await api.post('/auth/login', { email, password })
    const data = res.data.data || res.data  // Handle both wrapped and unwrapped
    get().setAuth({
        token: data.access_token,
        refreshToken: data.refresh_token,
        // ...
    })
}
```

## ImageMagick for Photo Composition

For composing multiple photos into a final photobooth image:

```go
// Use ImageMagick montage command
cmd := exec.Command("montage", 
    photoPaths...,
    "-geometry", fmt.Sprintf("%dx%d+10+10", width/len(photoPaths)-20, height/len(photoPaths)-20),
    "-tile", fmt.Sprintf("%dx1", len(photoPaths)),
    "-background", "white",
    "-gravity", "center",
    outputPath,
)
```

**Fallback:** If ImageMagick is not installed, use Go's `image/draw` package for basic composition.

## Docker Compose Cleanup

- Remove `version:` attribute (obsolete, triggers warning)
- `docker compose down --remove-orphans` cleans up removed services
- `docker compose up -d` rejected by terminal tool as long-lived → use `background=True`
