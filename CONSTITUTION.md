# Hermes Engineering Constitution

> *Ratified 2026-06-10. Highest-priority benchmark for all projects.*
> *Overridable only by explicit user directive. Not implicitly — explicitly.*

---

## 0. Preamble: Burn VPS, Not Tokens

Every engineering decision begins with a single question: **"Does this consume tokens or compute?"**

Tokens are scarce, metered, and expensive. Compute is abundant, sunk-cost, and already paid for. This VPS has 2 CPU cores, 12GB RAM, and runs 24/7 regardless of whether it's busy or idle. Every CPU cycle burned on this machine is free. Every token sent to an LLM API costs money.

Therefore:

- **Prefer computation over LLM calls.** Sort, filter, aggregate, index, and cache data locally rather than asking an LLM to reason over it.
- **Prefer databases over context windows.** Store facts in PostgreSQL, not in agent memory. Query them. Index them. A `SELECT` with a `WHERE` clause costs zero tokens.
- **Prefer caching over re-fetching.** HTTP responses, API results, computed values — cache them. A stale cache hit is cheaper than a fresh LLM call.
- **Prefer deterministic algorithms over probabilistic reasoning.** If a problem can be solved with a regex, a SQL query, or a hash lookup, solve it that way. Reserve LLMs for problems that genuinely require semantic understanding.
- **Prefer indexing over scanning.** Design schemas with the queries in mind. Every full-table scan is an opportunity for an index.

This is not a suggestion. It is the operating system's defining principle.

---

## 1. Technology Stack

### 1.1 Fixed Standards (Non-Negotiable)

| Layer | Technology | Version Floor | Rationale |
|-------|-----------|---------------|-----------|
| **Frontend** | React + TypeScript + Tailwind CSS | React 18+, TS 5.x, Tailwind 3.x | Type safety, component reuse, utility-first CSS |
| **Backend** | Go | 1.22+ | Single binary, zero-dependency deploys, fast cold starts |
| **Database** | PostgreSQL | 16+ | JSONB, full-text search, row-level security, proven reliability |
| **Container** | Docker + Docker Compose | Latest stable | Reproducible environments, service isolation |
| **CDN/DNS** | Cloudflare | — | DDoS protection, DNS management, Pages/Workers when needed |
| **VCS** | Git + GitHub | — | Industry standard, CI/CD integration |

### 1.2 Banned Technologies

These are explicitly banned unless the user grants a specific exception:

- **No ORMs.** Use raw SQL with `pgx` (Go) or `psycopg2` (Python). ORMs obscure query performance and prevent PostgreSQL-specific optimizations.
- **No MongoDB or NoSQL as primary store.** PostgreSQL handles JSONB natively. Two databases are worse than one.
- **No serverless frameworks on VPS.** This machine runs 24/7. Serverless abstractions add cost and complexity with no benefit.
- **No microservices for single-team projects.** Start with a monolith. Split services only when a measurable bottleneck demands it.
- **No GraphQL unless the frontend has 5+ distinct data shapes.** REST with query parameters covers 90% of use cases with simpler caching and debugging.

### 1.3 Permitted Additions (Require Justification)

- **Redis** — Only if you need sub-millisecond cache access or rate limiting at scale. PostgreSQL can cache most things.
- **Message queues (NATS/RabbitMQ)** — Only if you have async workflows that cannot be handled by PostgreSQL's `LISTEN/NOTIFY` or a simple job table.
- **WebSockets** — Only if polling at 5-second intervals is truly insufficient.
- **Additional languages** — Python for data scripts and AI orchestration is accepted. Node.js for tooling is accepted. Everything else requires explicit approval.

---

## 2. Project Architecture

### 2.1 Directory Structure

Every project follows this structure. Deviations require a comment explaining why.

```
project-name/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # Entry point, wire up dependencies
│   ├── internal/
│   │   ├── config/              # Environment-based configuration
│   │   ├── database/            # Connection pool, migrations, queries
│   │   ├── handler/             # HTTP handlers (one file per resource)
│   │   ├── middleware/          # Auth, logging, CORS, rate limiting
│   │   ├── model/               # Domain types, request/response structs
│   │   └── service/             # Business logic (no HTTP concerns)
│   ├── migrations/              # SQL migration files (numbered, reversible)
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/          # Reusable UI components
│   │   ├── pages/               # Route-level components
│   │   ├── hooks/               # Custom React hooks
│   │   ├── lib/                 # API client, utilities, constants
│   │   └── styles/              # Global styles, Tailwind config
│   ├── public/                  # Static assets
│   ├── package.json
│   ├── tsconfig.json
│   ├── tailwind.config.ts
│   ├── vite.config.ts
│   └── Dockerfile
├── docker-compose.yml           # Service orchestration
├── .env.example                 # Template for required env vars
├── Makefile                     # Common commands (build, test, migrate, deploy)
└── README.md                    # Project overview, setup, architecture
```

### 2.2 API Design Standards

- **REST first.** Resources are nouns. HTTP methods are verbs. Status codes are meaningful.
- **Version from day one.** Every endpoint lives under `/api/v1/`. Future versions get `/api/v2/`. Never change a v1 endpoint's contract — add v2 and deprecate v1 with a sunset header.
- **JSON request/response.** No XML. No form-encoded responses. JSON only.
- **Consistent error shape:**
  ```json
  {
    "error": {
      "code": "VALIDATION_ERROR",
      "message": "Human-readable description",
      "details": [{ "field": "email", "reason": "invalid format" }]
    }
  }
  ```
- **Pagination for all list endpoints.** Cursor-based preferred. Offset-based accepted with explicit `limit`/`offset` params.
- **No nested resources beyond one level.** `/api/v1/users/123/orders` is acceptable. `/api/v1/users/123/orders/456/items` is not. Flatten it.

### 2.3 Database Standards

- **Migrations are numbered and reversible.** `001_create_users.up.sql` and `001_create_users.down.sql`. Never edit a migration that has been applied to any environment — create a new migration instead.
- **Foreign keys are mandatory.** Every relationship is enforced at the database level. Application-level "soft references" are bugs waiting to happen.
- **Timestamps are `TIMESTAMPTZ`.** Never `TIMESTAMP` without timezone. Never store timezone-naive dates.
- **UUIDs for primary keys exposed to clients.** Auto-increment integers are fine for internal join tables. Anything visible in a URL or API response uses UUIDv7 (time-ordered).
- **Index every column used in a `WHERE` clause.** Run `EXPLAIN ANALYZE` on every query that touches more than 100 rows. If it's a sequential scan, add an index.
- **JSONB for flexible data.** If a column would require a schema migration every time a new field is added, use JSONB. If the field is queried frequently, extract it to a real column.

---

## 3. Deployment Standards

### 3.1 Docker Compose

Every project ships as a `docker-compose.yml` with exactly three services: `db`, `backend`, `frontend`. No more unless justified.

```yaml
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build: ./backend
    environment:
      DATABASE_URL: postgres://${DB_USER}:${DB_PASSWORD}@db:5432/${DB_NAME}?sslmode=disable
    depends_on:
      db:
        condition: service_healthy
    ports:
      - "${API_PORT:-8080}:8080"

  frontend:
    build: ./frontend
    depends_on:
      backend:
        condition: service_started
    ports:
      - "${FRONTEND_PORT:-3000}:3000"

volumes:
  pgdata:
```

### 3.2 Cloudflare Wrangler

Static sites and edge functions use Cloudflare. Every subdomain follows this pattern:

```toml
# wrangler.toml
name = "project-name"
compatibility_date = "2026-01-01"

[env.production]
route = { pattern = "project.client.arjism.com", custom_domain = true }

[env.production.vars]
API_URL = "https://api.client.arjism.com"
```

- **DNS:** `*.client.arjism.com` CNAME → Cloudflare Pages or origin server.
- **HTTPS:** Always on. Never serve over plain HTTP. Cloudflare's "Always Use HTTPS" rule is enabled by default.
- **Caching:** Static assets get `Cache-Control: public, max-age=31536000, immutable` with content-hash filenames.

---

## 4. Coding Standards

### 4.1 Go

- **`gofmt` is not optional.** CI rejects unformatted code.
- **Errors are values.** Never `panic()` in library code. Return errors up the stack. The only acceptable `panic` is in `main()` when a critical dependency (database, config) fails to initialize.
- **No global state.** Pass dependencies explicitly via struct fields or function parameters. No `init()` functions that mutate global variables.
- **Context-aware everything.** Every function that does I/O accepts a `context.Context` as its first parameter. Every goroutine has a context.
- **Connection pooling.** Use `pgxpool` with min=2, max=10 connections. Never open a new connection per request.
- **Structured logging.** Use `slog` (stdlib). Every log line includes a request ID. No `fmt.Println` in production code.

### 4.2 TypeScript / React

- **Strict mode.** `tsconfig.json` has `"strict": true`. No exceptions. No `any` without a comment explaining why.
- **Components are functions.** No class components. One component per file. File name matches component name.
- **Props are typed.** Every component has an explicit interface. No implicit `any` props.
- **API calls are typed.** Request and response types are defined in a shared `types.ts`. Never `as` cast API responses — validate them at runtime or trust your backend contract.
- **Tailwind only.** No CSS modules. No styled-components. No inline `style={{}}` except for dynamic values. All styling goes through Tailwind classes.
- **No state management library unless needed.** React Context + `useReducer` covers 90% of state. Add Zustand or Jotai only when you have cross-cutting state that causes prop drilling across 3+ levels.

### 4.3 Python (Scripts & AI Orchestration)

- **Type hints on all function signatures.** The script may be throwaway but the next person reading it won't know that.
- **No class unless you have state.** Functions are fine. Classes with one method and `__init__` are not.
- **`urllib` from stdlib over `requests`.** One fewer dependency. One fewer CVEs to track.
- **`pathlib` over `os.path`.** It's 2026. String concatenation for paths is a code smell.

### 4.4 SQL

- **Keywords in UPPERCASE.** `SELECT`, `FROM`, `WHERE`, `JOIN`. Identifiers in lowercase.
- **Explicit column lists.** Never `SELECT *` in application code. Database schema changes shouldn't break your Go structs silently.
- **Parameterized queries.** Never string-concatenate user input into SQL. `$1`, `$2`, `$3` — always.
- **Migrations are idempotent.** Use `IF NOT EXISTS`, `IF EXISTS`, `CREATE OR REPLACE`.

---

## 5. Security

### 5.1 Secrets Management

- **Never in source code.** Not in comments. Not in "temporary" variables. Not in `config.go` defaults.
- **Never in fact_store or memory.** These are accessible to AI agents. Secrets go in `.env` files (loaded via `hermes config set`) or Cloudflare environment variables.
- **Environment variables are the single source of truth.** Every secret has exactly one canonical location. If a secret appears in two places, one of them is a bug.
- **GitHub PATs are scoped.** Minimum permissions. Classic tokens with expiration. Fine-grained tokens preferred.

### 5.2 API Security

- **JWT for authentication.** Access tokens: 15-minute expiry. Refresh tokens: 7-day expiry. Rotate refresh tokens on use.
- **CORS is restrictive.** Allow only the specific frontend origin. Never `Access-Control-Allow-Origin: *`.
- **Rate limiting on all public endpoints.** Token bucket or sliding window. 100 requests/minute for unauthenticated, 1000/minute for authenticated.
- **Input validation at the handler layer.** Never trust request bodies. Validate before they touch business logic.

### 5.3 Infrastructure Security

- **PostgreSQL never exposed to the internet.** Bind to `127.0.0.1` or Docker network. No port forwarding on the host.
- **SSH key authentication only.** Password login disabled. Port 22 is the only open port unless a service explicitly needs exposure.
- **Automatic updates for base images.** Rebuild Docker images weekly. Alpine security patches ship fast — use them.

---

## 6. Testing

### 6.1 What to Test

| Priority | Test Type | When | Coverage Target |
|----------|-----------|------|-----------------|
| **Critical** | API integration tests | Every endpoint | 100% of handlers |
| **High** | Business logic unit tests | Every service function | 80% of service layer |
| **Medium** | Database migration tests | Every up+down migration | 100% of migrations |
| **Low** | Frontend component tests | Critical user flows only | Happy paths |
| **Skip** | Frontend unit tests for simple components | Tailwind-styled divs | — |

### 6.2 Testing Rules

- **Tests must pass before merge.** No exceptions. Broken tests are broken builds.
- **Tests use a real database.** An in-memory SQLite is not PostgreSQL. Spin up a test container or use a test database.
- **No mocks for the database layer.** Mock at the HTTP or service boundary. The database is not an implementation detail — it's part of the system.
- **Test error paths.** Happy-path tests prove the feature works. Error-path tests prove it fails safely.

---

## 7. Performance

### 7.1 Targets

| Metric | Target | Measured At |
|--------|--------|-------------|
| API response (p50) | < 50ms | Backend handler |
| API response (p99) | < 500ms | Backend handler |
| Database query | < 20ms | `EXPLAIN ANALYZE` |
| Frontend LCP | < 2.5s | Lighthouse |
| Frontend TTI | < 3.5s | Lighthouse |
| Cold start | < 5s | `docker compose up` |

### 7.2 Caching Strategy

- **Static assets: 1 year.** Content-hash filenames (`main.a3f2b1c.js`). Immutable cache headers.
- **API responses: case-by-case.** List endpoints with stable data get `Cache-Control: public, max-age=60`. Mutation endpoints never get cached.
- **Database query results: application-level.** Use a TTL cache (hashmap + expiry) for queries that run more than 10 times per second. PostgreSQL's query cache handles the rest.

### 7.3 Database Performance

- **Connection pooling.** `pgxpool` with min=2, max=10. Never one-connection-per-request.
- **Index before query.** If a query appears in application code, its columns appear in an index.
- **`EXPLAIN ANALYZE` on every query during development.** Sequential scans on tables with more than 1000 rows are a bug.

---

## 8. Documentation

### 8.1 Minimum Viable Documentation

Every project requires:
- **`README.md`** — What this project does, how to run it locally, how to deploy it.
- **`ARCHITECTURE.md`** — High-level system design, data flow, service boundaries.
- **API documentation** — Every endpoint with request/response examples. Generated from code if possible (Go doc comments, Swagger).
- **Environment variables** — `.env.example` with every variable, its purpose, and its default (if any).

### 8.2 When to Write ADRs

Architecture Decision Records (ADRs) are required for:
- Choosing between two viable technologies
- Changing an existing architectural pattern
- Adding a new service or data store
- Any decision that would surprise a new team member

ADRs follow the format: **Title, Status (Proposed/Accepted/Deprecated/Superseded), Context, Decision, Tradeoffs, Alternatives Considered.**

---

## 9. AI Usage & Token Optimization

### 9.1 The Token Hierarchy

Not all token spending is equal. Spend tokens where they create compounding value:

| Tier | Use Case | Value |
|------|----------|-------|
| **Invest** | Architecture decisions, system design, security review | High leverage — one good decision saves weeks |
| **Build** | Feature implementation, complex debugging | Medium leverage — accelerates development |
| **Analyze** | Code review, test generation, documentation | Medium leverage — improves quality |
| **Assist** | Simple scripts, boilerplate, formatting | Low leverage — could be a template |
| **Waste** | Asking LLM to sort a list, query a DB, do math | Negative leverage — computation is free |

### 9.2 Agent Usage Principles

- **One task, one agent.** Don't ask an agent to "build the whole app." Decompose into discrete, verifiable tasks.
- **Delegate expensive reasoning.** Use `delegate_task` with specialist agents for complex subtasks. The orchestrator stays lean.
- **Verify, don't trust.** Every subagent output with external side effects must be independently verified.
- **Session_search before asking.** The session database has what was said. Query it before making the user repeat themselves.
- **Skills encode patterns.** When you solve a recurring problem, encode it as a skill. Skills are cached knowledge — tokens you only spend once.

### 9.3 Token Budget

- **Every agent interaction should produce a verifiable artifact.** If the output is "let me think about that," it wasn't worth the tokens.
- **Prefer `terminal` over `execute_code` for simple operations.** `execute_code` loads an interpreter. `terminal` runs a shell command.
- **Batch file reads.** Use parallel `read_file` calls instead of sequential reads. Every round trip costs tokens.
- **Archive, don't accumulate.** When a skill, agent, or config is no longer relevant, archive it. Don't keep it loaded in prompt context.

### 9.4 When to Use Which Model

| Model | Use For | Max Tokens/Task |
|-------|---------|-----------------|
| `deepseek/deepseek-v4-pro` | Architecture, complex debugging, multi-step reasoning | High |
| `openrouter/owl-alpha` | Specialist subagents, feature implementation | Medium |
| Cheapest available model | Formatting, translations, simple classification | Low |

---

## 10. Governance

### 10.1 Amendment Process

1. Any agent may propose an amendment to this constitution.
2. The amendment must cite the specific section being changed, the rationale, and a tradeoff analysis.
3. The user (kvinn) must explicitly approve the amendment.
4. Approved amendments are appended to the end of this document with a version bump and date.

### 10.2 Constitutional Precedence

- This constitution overrides all skill files, agent SOUL files, and project-level conventions.
- If a skill contradicts this constitution, the constitution wins. Update the skill.
- If a project-level decision contradicts this constitution, the project must document the exception and justification.
- Only the user can override this constitution. Agents cannot.

### 10.3 Review Cycle

This constitution shall be reviewed:
- After every major project completion
- When a new technology is adopted
- At minimum, every 6 months

---

## 11. Appendices

### A. Quick Reference: Technology Choices

| Problem | Solution | Why |
|---------|----------|-----|
| State management | React Context + useReducer | Sufficient for 90% of UIs |
| API calls | Fetch + typed wrapper | No dependencies |
| CSS | Tailwind utility classes | No CSS file management |
| Auth | JWT + refresh tokens | Stateless, scalable |
| File uploads | S3-compatible / Cloudflare R2 | Cheap, CDN-ready |
| Background jobs | PostgreSQL job table + poller | No extra infrastructure |
| Search | PostgreSQL full-text search | Good enough for <1M rows |
| Caching | In-memory TTL map + Cloudflare CDN | Simple, effective |
| Monitoring | Structured logs + health check endpoint | Minimum viable observability |
| CI/CD | GitHub Actions | Free for public repos, integrated |

### B. Quick Reference: File Naming

| Type | Convention | Example |
|------|-----------|---------|
| Go files | `snake_case.go` | `user_handler.go` |
| React components | `PascalCase.tsx` | `UserProfile.tsx` |
| SQL migrations | `NNN_description.up.sql` | `001_create_users.up.sql` |
| Dockerfiles | `Dockerfile` (no extension) | `Dockerfile` |
| Environment files | `.env` (gitignored), `.env.example` (committed) | `.env.example` |
| whatILearn logs | `YYYYMMDD-camelCase.md` | `20260610-leadGenDeployment.md` |

### C. Quick Reference: Ports

| Service | Port |
|---------|------|
| PostgreSQL | 5432 (inside Docker), mapped to host as needed |
| Go API | 8080 |
| React dev server | 5173 (Vite) |
| React production (Nginx) | 3000 |
| SSH | 22 |

---

*Constitution v1.0 — Ratified 2026-06-10*
*Last amended: —*
*Next review: 2026-12-10*
