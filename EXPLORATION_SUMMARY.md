# Comprehensive Exploration: Manifesto Library + CLI

## REPOSITORY 1: `/Users/abraxas/manifesto` — Core Library

### Directory Structure

```
manifesto/
├── cmd/
│   ├── container.go           # Root DI wiring (infrastructure: DB, Redis, FS)
│   └── server.go              # Fiber HTTP server template with middleware
│
├── internal/
│   ├── kernel/                # Domain primitives
│   │   ├── common_ids.go      # UserID, TenantID types
│   │   ├── context.go         # AuthContext with scope validation
│   │   ├── store.go           # Store interface for contextual data
│   │   ├── proj_ids.go        # Project-specific ID types (placeholder)
│   │   ├── proj_objvalue.go   # Project object values (placeholder)
│   │   └── common_objvalue.go # Common value objects
│   │
│   ├── errx/                  # Structured error handling
│   │   ├── error.go           # Error struct with code/message/type/httpstatus/details
│   │   ├── http.go            # HTTP status mapping
│   │   ├── types.go           # Error type enum (NOT_FOUND, UNAUTHORIZED, etc.)
│   │   ├── common.go          # Pre-defined common errors
│   │   └── regestry.go        # Error code registry
│   │
│   ├── logx/                  # Structured logging
│   │   ├── api.go             # Package-level API (Trace, Debug, Info, Warn, Error, Fatal)
│   │   ├── logger.go          # Logger struct with level/formatter/output
│   │   ├── levels.go          # Log level constants (Trace, Debug, Info, Warn, Error, Fatal)
│   │   ├── entry.go           # LogEntry for context keys, request IDs, etc.
│   │   ├── config.go          # LogConfig (level, format, json_fields)
│   │   ├── formatter.go       # Formatter interface
│   │   ├── console_formatter.go # Colored terminal output
│   │   └── json_formatter.go  # JSON log output
│   │
│   ├── config/                # Environment-driven configuration
│   │   ├── config.go          # Main Config struct + Load() function
│   │   ├── server.go          # Server config (port, env, log level)
│   │   ├── database.go        # DB config (host, port, user, password, dbname, sslmode)
│   │   ├── auth.go            # JWT, session, OTP config
│   │   ├── email.go           # Email provider config (SES/console)
│   │   ├── jobx.go            # Job queue config (redis address, queues, workers)
│   │   ├── notifx.go          # Notification config (provider, from, region, bucket)
│   │   └── oauth.go           # OAuth providers (Google, Microsoft) secrets
│   │
│   ├── ptrx/                  # Pointer utilities
│   │   └── ptrx.go            # 15+ helpers: Ptr(), Val(), SlicePtr(), MapPtr(), etc.
│   │
│   ├── asyncx/                # Async primitives + concurrency helpers
│   │   ├── doc.go             # Full package documentation
│   │   └── asyncx.go          # Futures, All/AllSettled/Race, Map/ForEach/Pool, 
│   │                           # Retry/RetryWithBackoff, Debounce/Throttle, Once, Do, DoCtx
│   │
│   ├── fsx/                   # File storage abstraction
│   │   ├── fsx.go             # FileReader, FileWriter, FileDeleter, PathOperations,
│   │   │                       # FileSystem, PresignedURLGenerator interfaces
│   │   ├── fsxlocal/          # Local disk implementation
│   │   └── fsxs3/             # AWS S3 implementation
│   │
│   ├── jobx/                  # Async job queue (Redis-backed)
│   │   ├── jobx.go            # Job dispatcher interface
│   │   ├── models.go          # Job struct with ID, type, payload, status, retries
│   │   ├── options.go         # JobOptions (retries, timeout, priority, queue)
│   │   ├── errors.go          # Job-specific errors
│   │   └── jobxredis/         # Redis implementation of dispatcher
│   │
│   ├── notifx/                # Email notification service
│   │   ├── notifx.go          # EmailSender interface (Send, SendBatch)
│   │   ├── models.go          # EmailMessage struct (to, subject, body, attachments)
│   │   ├── template.go        # HTML template rendering for emails
│   │   ├── options.go         # SendOptions (priority, retries, bcc, etc.)
│   │   ├── errors.go          # Notification errors
│   │   ├── notifxses/         # AWS SES implementation
│   │   └── notifxconsole/     # Console (debug) implementation
│   │
│   ├── ai/                    # LLM + embeddings + vector store
│   │   ├── llm/               # OpenAI, Anthropic, Azure, Bedrock, Google Gemini clients
│   │   ├── embedding/         # Embedding service (text → vectors)
│   │   ├── vstore/            # Vector store (Pgvector, other backends)
│   │   ├── document/          # Document struct + chunking strategies
│   │   ├── ocr/               # Document OCR (image → text)
│   │   ├── speech/            # Speech-to-text, text-to-speech
│   │   ├── providers/         # LLM provider configs (models, API keys)
│   │   └── *examples*         # Code examples in /examples/ai/
│   │
│   └── iam/                   # Full IAM system (auth, users, tenants, RBAC, API keys)
│       ├── doc.go             # Comprehensive endpoint reference + architecture
│       ├── iam.go             # IAM package root
│       ├── auth/              # OAuth2 + JWT + refresh tokens
│       │   ├── authinfra/     # PostgreSQL + Redis implementations
│       │   ├── OAuth endpoints (Google, Microsoft)
│       │   └── JWT generation, refresh, logout
│       ├── user/              # User entity + domain logic
│       │   ├── usersrv/       # Service layer (scope assignment, validation)
│       │   └── userinfra/     # PostgreSQL repository
│       ├── tenant/            # Tenant (organization) entity
│       │   ├── tenantsrv/     # Service layer
│       │   └── tenantinfra/   # PostgreSQL repository
│       ├── otp/               # One-time password for passwordless auth
│       │   ├── otpsrv/        # OTP generation, validation, rate limiting
│       │   └── otpinfra/      # PostgreSQL + Redis implementations
│       ├── invitation/        # User invitation flow
│       │   └── Full endpoint support (create, list, validate, revoke, delete)
│       ├── apikey/            # API key authentication
│       │   ├── apikeyapi/     # HTTP handlers
│       │   ├── apikeysrv/     # Service layer
│       │   └── apikeyinfra/   # PostgreSQL repository
│       ├── scopes/            # Scope validation + templates
│       └── iamcontainer/      # DI wiring for all IAM subdomains
│
├── examples/                  # Executable examples
│   ├── ai/                    # LLM, embeddings, vector store usage
│   ├── asyncx/                # Concurrency patterns
│   ├── fsx/                   # File storage examples
│   └── vectorstore/           # Vector search examples
│
├── migrations/                # SQL migrations (for IAM: users, tenants, etc.)
│
├── go.mod                     # Go module declaration + dependencies
├── Makefile                   # Development tasks (test, lint, build, migrate)
├── docker-compose.yml         # PostgreSQL + Redis services
├── README.md                  # Main documentation
└── init-project.sh            # Shell script for initializing new projects
```

---

### Dependencies (go.mod)

**LLM Providers:**
- `github.com/openai/openai-go/v3` — OpenAI API
- `github.com/anthropics/anthropic-sdk-go` — Anthropic (Claude) API
- `google.golang.org/genai` — Google Gemini API
- `github.com/aws/aws-sdk-go-v2/service/bedrockruntime` — AWS Bedrock
- `github.com/Azure/azure-sdk-for-go/sdk/azcore` — Azure services

**Cloud & Storage:**
- `github.com/aws/aws-sdk-go-v2` — AWS SDK (S3, SES, Bedrock)
- `github.com/aws/aws-sdk-go-v2/service/s3` — S3 file storage
- `github.com/aws/aws-sdk-go-v2/service/ses` — Email via SES

**Web Framework & Auth:**
- `github.com/gofiber/fiber/v2` — Fiber HTTP server
- `github.com/golang-jwt/jwt/v5` — JWT tokens

**Database & Cache:**
- `github.com/lib/pq` — PostgreSQL driver
- `github.com/jmoiron/sqlx` — SQL query builder + mapping
- `github.com/redis/go-redis/v9` — Redis client

**Utilities:**
- `github.com/google/uuid` — UUID generation
- `golang.org/x/crypto` — Cryptography (hashing, encryption)

**Go Version:** 1.25.4

---

### All Packages & Their Responsibilities

| Package | LOC | Purpose | Key Types/Interfaces |
|---------|-----|---------|----------------------|
| **kernel** | 160 | Domain primitives | `UserID`, `TenantID`, `AuthContext` (with scope validation), `Store` |
| **errx** | 300+ | Structured errors | `Error` struct, `Type` enum, HTTP mapping, registry, JSON marshaling |
| **logx** | 813 | Structured logging | `Logger`, `Level`, `Entry`, `Formatter` (console + JSON), colored output |
| **config** | 200+ | Env-driven config | `Config`, `Load()`, database/auth/email/jobx/notifx/oauth config structs |
| **ptrx** | 798 | Pointer utilities | `Ptr()`, `Val()`, `SlicePtr()`, `MapPtr()`, compare/coalesce helpers |
| **asyncx** | ~400 | Concurrency | `Future[T]`, `All()`, `AllSettled()`, `Race()`, `Map()`, `ForEach()`, `Pool()`, `Retry()`, `Debounce()`, `Throttle()`, `Once()` |
| **fsx** | ~150 | File storage | `FileReader`, `FileWriter`, `FileDeleter`, `PathOperations`, `FileSystem`, `PresignedURLGenerator` |
| **fsx/fsxlocal** | ~100 | Local disk | Implements `FileSystem` interface for local directory |
| **fsx/fsxs3** | ~200 | AWS S3 | Implements `FileSystem` + `PresignedURLGenerator` for S3 |
| **jobx** | ~200 | Job queue | `Job` struct, `Dispatcher` interface, job options, Redis implementation |
| **jobx/jobxredis** | ~300 | Redis jobs | Redis-backed async dispatcher with retries, priorities, queues |
| **notifx** | ~180 | Email service | `EmailSender`, `EmailMessage`, template rendering, SES + console |
| **notifx/notifxses** | ~150 | AWS SES | SES-backed email delivery |
| **notifx/notifxconsole** | ~50 | Debug email | Console logger for development |
| **ai** | ~500 | LLM + embeddings | LLM clients (OpenAI, Anthropic, Azure, Bedrock, Gemini), embedding service, vector store, OCR, speech |
| **ai/llm** | ~300 | LLM providers | Provider factories, model selection, token counting |
| **ai/embedding** | ~150 | Text → vectors | `EmbeddingService`, chunking strategies |
| **ai/vstore** | ~200 | Vector DB | `VectorStore` interface, Pgvector implementation |
| **ai/ocr** | ~100 | Image → text | Document OCR using provider services |
| **ai/speech** | ~100 | Audio ↔ text | Speech-to-text, text-to-speech |
| **iam** | ~1500 | Full IAM | OAuth2, JWT, OTP, API keys, users, tenants, invitations, RBAC |
| **iam/auth** | ~400 | Auth flow | OAuth callbacks, JWT refresh, session management, middleware |
| **iam/user** | ~300 | User entity | User struct, DTOs, scope assignment, password resets |
| **iam/tenant** | ~250 | Tenant entity | Multi-tenancy, subscription plans, user limits, lifecycle |
| **iam/otp** | ~200 | Passwordless | OTP generation, verification, rate limiting, email delivery |
| **iam/apikey** | ~200 | API auth | API key generation, validation, rotation, metadata |
| **iam/invitation** | ~150 | Onboarding | Invitation tokens, acceptance, expiration, revocation |
| **iam/scopes** | ~100 | Authorization | Scope validation, templates (admin, viewer, analyst, etc.) |

---

### Key Types, Interfaces & Patterns

#### **Kernel Primitives**
```go
type UserID string         // Typed ID for users
type TenantID string       // Typed ID for tenants

type AuthContext struct {
    UserID   *UserID       // nil for unauthenticated
    TenantID TenantID
    Email    string
    Scopes   []string      // e.g., ["users:read", "reports:write"]
    IsAPIKey bool
}

func (ac *AuthContext) HasScope(scope string) bool      // Checks scope + wildcards
func (ac *AuthContext) IsAdmin() bool                   // Has "*" or "admin:*"
```

#### **Error Handling**
```go
type Error struct {
    Code       string                 // Unique code (e.g., "USER.NOT_FOUND")
    Message    string
    Type       Type                   // NOT_FOUND, UNAUTHORIZED, INVALID, etc.
    HTTPStatus int                    // Auto-mapped (404, 401, 400, 500)
    Details    map[string]interface{} // Context about the error
    Err        error                  // Underlying cause
}

// Usage:
errx.New("User not found", errx.NotFound).
    WithDetail("user_id", id)
```

#### **Structured Logging**
```go
logx.Info("User created")
logx.Infof("User %s created", userID)
logx.WithContext(ctx).Warnf("Request timeout")
logx.WithError(err).Errorf("Database error")

// Outputs colored (dev) or JSON (prod) depending on formatter
```

#### **Concurrency Patterns**
```go
// Future-based async
fut := asyncx.Run(func() (*User, error) {
    return repo.GetByID(ctx, id)
})
user, err := fut.Await()

// Fan-out (wait all)
results, err := asyncx.All(ctx,
    func(ctx context.Context) (*A, error) { ... },
    func(ctx context.Context) (*B, error) { ... },
)

// Worker pool (bounded concurrency)
results, err := asyncx.Pool(ctx, 10, items,
    func(ctx context.Context, item T) (Result, error) { ... },
)

// Retry with backoff
data, err := asyncx.RetryWithBackoff(ctx, 5, 100*time.Millisecond,
    func(ctx context.Context) (*Data, error) { ... },
)
```

#### **File Storage Abstraction**
```go
type FileSystem interface {
    ReadFile(ctx context.Context, path string) ([]byte, error)
    WriteFile(ctx context.Context, path string, data []byte) error
    DeleteFile(ctx context.Context, path string) error
    List(ctx context.Context, path string) ([]FileInfo, error)
    Exists(ctx context.Context, path string) (bool, error)
    CreateDir(ctx context.Context, path string) error
    Join(elem ...string) string
}

type PresignedURLGenerator interface {
    GetPresignedDownloadURL(ctx, path string, expiration time.Duration) (string, error)
    GetPresignedUploadURL(ctx, path string, expiration time.Duration) (string, error)
}

// Can swap implementations: fsxlocal ↔ fsxs3
```

#### **Job Queue Pattern**
```go
type Job struct {
    ID       string           // UUID
    Type     string           // Job type (e.g., "send_email")
    Payload  map[string]interface{}
    Status   JobStatus        // PENDING, PROCESSING, FAILED, COMPLETED
    Retries  int
    Queue    string           // Redis queue name
    Attempt  int
}

type Dispatcher interface {
    Dispatch(ctx context.Context, job *Job, opts ...JobOption) error
    Subscribe(ctx context.Context, queueName string, handler func(*Job) error)
}

// Implementations support retries, timeouts, priorities, queue selection
```

#### **Email Service Pattern**
```go
type EmailMessage struct {
    To          []string
    Subject     string
    Body        string
    Attachments []*Attachment
    TemplateID  string        // For template rendering
}

type EmailSender interface {
    Send(ctx context.Context, msg *EmailMessage) error
    SendBatch(ctx context.Context, msgs []*EmailMessage) error
}

// Can use SES or console (debug)
```

#### **LLM + Embeddings Pattern**
```go
// LLM client abstraction (OpenAI, Anthropic, Bedrock, Gemini, Azure all compatible)
type LLMClient interface {
    Complete(ctx context.Context, prompt string, opts ...LLMOption) (string, error)
    Stream(ctx context.Context, prompt string, handler func(string) error) error
    CountTokens(prompt string) int
}

// Embeddings
type EmbeddingService interface {
    Embed(ctx context.Context, text string) ([]float32, error)
    EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
}

// Vector store (Pgvector)
type VectorStore interface {
    Insert(ctx context.Context, vector []float32, metadata map[string]interface{}) error
    Search(ctx context.Context, vector []float32, limit int) ([]SearchResult, error)
}
```

#### **IAM Architecture**

**Authentication:**
- **OAuth2**: Google, Microsoft
- **Passwordless (OTP)**: Email-based 6-digit codes
- **JWT**: 15-min access tokens + 7-day refresh tokens
- **API Keys**: Machine-to-machine auth

**Authorization:**
- Scope-based RBAC (e.g., `users:read`, `reports:write`)
- Scope templates (super_admin, tenant_admin, viewer, analyst, etc.)
- Wildcard support (`*`, `resource:*`)

**Multi-Tenancy:**
- Every user belongs to exactly one tenant per email
- Same email can exist in multiple tenants independently
- Subscription plans enforce user limits (TRIAL: 5, BASIC: 5, PROF: 50, ENTERPRISE: 500)

**User Lifecycle:**
- Invited via token
- Signup (OAuth or OTP)
- Active/Suspended status
- Scope assignment
- Sessions + refresh tokens
- API keys per tenant

**Endpoints (Summary):**
- `POST /auth/login` — OAuth initiation
- `GET /auth/callback/:provider` — OAuth redirect
- `POST /auth/refresh` — Token refresh
- `POST /auth/logout` — Revoke session
- `POST /auth/passwordless/signup/initiate` — OTP signup
- `POST /auth/passwordless/login/initiate` — OTP login
- `POST /invitations` — Send invitation
- `POST /api-keys` — Create API key
- `GET /api-keys` — List API keys
- `PUT /api-keys/:id` — Update API key
- `POST /api-keys/:id/revoke` — Revoke API key

---

### Sample Architecture (cmd/)

**container.go:**
```go
type Container struct {
    Config    *config.Config
    DB        *sqlx.DB          // PostgreSQL
    Redis     *redis.Client
    FileSystem fsx.FileSystem    // Local or S3
    S3Client  *s3.Client
    
    // Module containers (added via `manifesto add`)
    // UserContainer *user.Container
    // TenantContainer *tenant.Container
}

func NewContainer(cfg *config.Config) *Container {
    c := &Container{Config: cfg}
    c.initInfrastructure()  // DB, Redis, FS
    c.initModules()         // Bounded context containers
    return c
}
```

**server.go:**
```go
func main() {
    cfg, _ := config.Load()
    container := NewContainer(cfg)
    
    app := fiber.New()
    setupMiddleware(app, cfg)  // CORS, request ID, logging, recovery
    
    app.Get("/health", healthCheck)
    app.Get("/", info)
    
    registerRoutes(app, container)  // Domain routes
    startServer(app, cfg)            // Graceful shutdown
}
```

---

### Configuration Pattern

**Environment variables** → Load via `config.Load()` → Type-safe config struct

```go
// pkg/config/config.go marker system:
type Config struct {
    Server    ServerConfig
    Database  DatabaseConfig
    Redis     RedisConfig
    // manifesto:config-fields  ← modules inject here
}

func Load() *Config {
    cfg.Server = loadServerConfig()
    cfg.Database = loadDatabaseConfig()
    cfg.Redis = loadRedisConfig()
    // manifesto:config-loads  ← modules inject here
    return cfg
}
```

---

## REPOSITORY 2: `/Users/abraxas/manifesto-cli` — Scaffolding Tool

### Directory Structure

```
manifesto-cli/
├── cmd/
│   └── manifesto/
│       └── main.go            # CLI entry point + version flag
│
├── internal/
│   ├── cli/
│   │   ├── root.go            # Root command + command registration
│   │   ├── init.go            # `manifesto init` command
│   │   ├── add.go             # `manifesto add` (module wiring OR domain scaffolding)
│   │   ├── install.go         # Deprecated alias to `add`
│   │   └── modules.go         # `manifesto modules` listing command
│   │
│   ├── config/
│   │   ├── manifest.go        # manifest.yaml structure + 12-module registry
│   │   └── wiring.go          # 6 wireable module specs + code injection templates
│   │
│   ├── scaffold/
│   │   ├── project.go         # InitProject() — orchestrates full project creation
│   │   ├── domain.go          # GenerateDomain() — DDD domain vertical generation
│   │   ├── module.go          # InstallModule() — module installation
│   │   └── wire.go            # WireModule() — code injection at marker points
│   │
│   ├── templates/
│   │   ├── embed.go           # Embedded template filesystem
│   │   ├── project/           # 4 project templates (container, server, makefile, docker)
│   │   └── domain/            # 8 domain templates (entity, port, errors, service, postgres, handler, container, ids)
│   │
│   ├── remote/
│   │   └── github.go          # GitHub archive downloader + import rewriting
│   │
│   └── ui/
│       ├── ui.go              # Colored output, spinners, progress messages
│       └── multiselect.go     # Interactive module selection UI
│
├── go.mod                     # CLI dependencies
├── go.sum
├── Makefile                   # Build, install, release tasks
├── .goreleaser.yaml           # Release automation
├── README.md                  # CLI documentation
└── .github/workflows/         # CI/CD pipelines
```

---

### Dependencies (go.mod)

- **CLI Framework**: `github.com/spf13/cobra v1.8.1`
- **Terminal Colors**: `github.com/fatih/color v1.18.0`
- **YAML Parsing**: `gopkg.in/yaml.v3`
- **Go Version**: 1.24.0+

---

### All CLI Commands

#### **1. `manifesto init <name> --module <go-module>`**

Creates a new project with all core libraries (kernel, errx, logx, config, ptrx).

**Flags:**
- `--module <path>` (required) — Go module path (e.g., `github.com/user/myapp`)
- `--with <mods>` — Comma-separated modules to wire: `fsx,asyncx,ai,jobx,notifx,iam`
- `--all` — Wire all optional modules
- `--quick` — Lightweight (excludes IAM, migrations)
- `--ref <tag|branch>` — Pin manifesto version (default: `main`)

**Output:**
```
<name>/
├── cmd/
│   ├── container.go      # DI wiring template
│   └── server.go         # Fiber server template
├── pkg/
│   ├── kernel/           # Typed IDs (UserID, TenantID)
│   ├── errx/             # Error handling
│   ├── logx/             # Logging
│   ├── config/           # Env config with manifesto markers
│   └── ptrx/             # Pointer utilities
├── go.mod
├── go.sum
├── .gitignore
├── Makefile              # 40+ development commands
├── docker-compose.yml    # PostgreSQL 15 + Redis 7
├── manifesto.yaml        # Project metadata + module tracking
└── (if --with fsx,asyncx,ai,jobx,notifx,iam: their source packages)
```

---

#### **2. `manifesto add <module>`** — Wire a Module

Wires one of 6 optional modules into an existing project.

**Wireable Modules:**

| Module | Dependencies | Injected Into | Effect |
|--------|--------------|---------------|--------|
| `fsx` | — | config, container | File storage (local/S3) |
| `asyncx` | — | — | Concurrency primitives |
| `ai` | fsx | — | LLM, embeddings, vector store, OCR, speech |
| `jobx` | asyncx | config, container, makefile | Redis job queue |
| `notifx` | — | config, container, makefile | Email (SES/console) |
| `iam` | — | container, server, makefile | Full auth system + migrations |

**Process:**
1. Resolve dependencies (e.g., adding `jobx` automatically downloads `asyncx`)
2. Download module source from GitHub (Abraxas-365/manifesto repo)
3. Rewrite Go imports to project module path
4. Inject code at marker comments in 4 files:
   - `pkg/config/config.go` — config struct fields + Load() code
   - `cmd/container.go` — imports, fields, initialization code
   - `cmd/server.go` — route registration, middleware
   - `Makefile` — environment variable documentation
5. Install Go dependencies
6. Update `manifesto.yaml`

**Cross-Module Bridges:**
- If `jobx` + `notifx` wired → auto-injects email handler into job dispatcher

**Example:**
```bash
manifesto add jobx
# Downloads jobx + asyncx (dep)
# Injects into config, container, server, makefile
# Updates go.mod with new dependencies
```

---

#### **3. `manifesto add <domain-path>`** — Scaffold Domain

Generates a full DDD domain vertical.

**Input:** Domain path like `pkg/recruitment/candidate`

**Output:** 7 files in the domain:
```
pkg/recruitment/candidate/
├── candidate.go              # Entity struct, DTOs, methods
├── port.go                   # Repository interface (contract)
├── errors.go                 # Error registry (errx codes)
├── candidatesrv/
│   └── service.go           # Business logic layer
├── candidateinfra/
│   └── postgres.go          # PostgreSQL repository
├── candidateapi/
│   └── handler.go           # HTTP GET/POST/PUT/DELETE handlers
└── candidatecontainer/
    └── container.go         # Domain-specific DI wiring
```

**Auto-Wiring:**
1. Appends typed ID to `pkg/kernel/ids.go`: `type CandidateID = kernel.ID`
2. Injects into `cmd/container.go` — field + initialization
3. Injects into `cmd/server.go` — route registration (CRUD endpoints)

**Name Transformations:**
- `candidate` (input) → `Candidate` (PascalCase for entity)
- `candidate` → `candidate` (snake_case for package)
- `candidate` → `CANDIDATE` (UPPER_SNAKE for error codes)
- `candidate` → `candidates` (plural for table names)

---

#### **4. `manifesto modules`**

Lists all available modules with their wire status.

**Output:**
```
Core Libraries (always present):
  • kernel      — Domain primitives, typed IDs, auth context
  • errx        — Structured error handling
  • logx        — Structured logging
  • config      — Environment-driven configuration
  • ptrx        — Pointer utilities

Wireable Modules:
  ● fsx         — File storage (local/S3)
  ● asyncx      — Concurrency primitives
  ● ai          — LLM, embeddings, vector store
  ○ jobx        — Async job queue (not wired)
  ○ notifx      — Email notifications (not wired)
  ○ iam         — Full IAM system (not wired)
```

(`●` = wired, `○` = not wired)

---

#### **5. `manifesto version`**

Prints CLI version (set via build flags in main.go)

---

#### **6. `manifesto install <module>`** (Deprecated)

Alias to `manifesto add`. Kept for backward compatibility.

---

### Project Templates Created

#### **cmd/container.go template:**
```go
type Container struct {
    Config     *config.Config
    DB         *sqlx.DB
    Redis      *redis.Client
    FileSystem fsx.FileSystem
    // manifesto:container-fields  ← modules inject here
}

func (c *Container) initInfrastructure() { /* DB, Redis, FS */ }
func (c *Container) initModules() { /* module containers */ }
func (c *Container) StartBackgroundServices(ctx context.Context) { /* async jobs, cleanups */ }
func (c *Container) Cleanup() { /* DB/Redis close */ }
```

#### **cmd/server.go template:**
```go
// Fiber HTTP server with:
// - Global middleware: CORS, request ID, logging, recovery, panic handler
// - Health check endpoint: GET /health
// - Info endpoint: GET /
// - Route registration: registerRoutes(app, container)
// - Graceful shutdown: SIGINT/SIGTERM handling
// - Custom error handler: Maps errx.Error → HTTP status
// - Body limit: 10MB for file uploads
```

#### **Makefile template (40+ commands):**
```makefile
make dev              # Run local server
make dev-watch        # Hot reload (air)
make build            # Compile binary
make test             # Run tests
make lint             # golangci-lint

make up               # Docker: bring up postgres+redis
make down             # Docker: bring down
make health           # Health check

make migrate          # Run SQL migrations (IAM)
make migrate-create   # New migration file
make db-reset         # Clean + migrate + seed

make psql             # Open psql shell
make redis-cli        # Open redis-cli
make env              # Display config variables

# Modules inject additional targets:
# - jobx: make jobx-stats, make jobx-clear
# - iam: make iam-admin, make iam-seed
```

#### **docker-compose.yml template:**
```yaml
version: '3'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: <project>
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
```

---

### Domain Templates Created

#### **entity.go:**
```go
type Candidate struct {
    ID        CandidateID
    CreatedAt time.Time
    UpdatedAt time.Time
    // domain fields
}

type CandidateDTO struct { /* DTO for API responses */ }
type CreateCandidateReq struct { /* Create request body */ }
type UpdateCandidateReq struct { /* Update request body */ }
```

#### **port.go (Repository Interface):**
```go
type CandidateRepository interface {
    Save(ctx context.Context, candidate *Candidate) error
    FindByID(ctx context.Context, id CandidateID) (*Candidate, error)
    FindAll(ctx context.Context, filters map[string]interface{}, page, size int) ([]*Candidate, int, error)
    Update(ctx context.Context, candidate *Candidate) error
    Delete(ctx context.Context, id CandidateID) error
}
```

#### **errors.go:**
```go
var (
    ErrCandidateNotFound    = errx.New("Candidate not found", errx.NotFound)
    ErrCandidateExists      = errx.New("Candidate already exists", errx.Conflict)
    ErrInvalidCandidateData = errx.New("Invalid candidate data", errx.InvalidInput)
)
```

#### **service.go:**
```go
type CandidateService struct {
    repo CandidateRepository
}

func (s *CandidateService) Create(ctx context.Context, req *CreateCandidateReq) (*Candidate, error) { ... }
func (s *CandidateService) GetByID(ctx context.Context, id CandidateID) (*Candidate, error) { ... }
func (s *CandidateService) List(ctx context.Context, page, size int) ([]*Candidate, error) { ... }
func (s *CandidateService) Update(ctx context.Context, id CandidateID, req *UpdateCandidateReq) (*Candidate, error) { ... }
func (s *CandidateService) Delete(ctx context.Context, id CandidateID) error { ... }
```

#### **postgres.go (Infrastructure):**
```go
type PostgresCandidateRepository struct {
    db *sqlx.DB
}

func (r *PostgresCandidateRepository) Save(ctx context.Context, c *Candidate) error {
    return r.db.NamedExecContext(ctx, insertSQL, c).Err
}

func (r *PostgresCandidateRepository) FindByID(ctx context.Context, id CandidateID) (*Candidate, error) {
    var c Candidate
    return &c, r.db.GetContext(ctx, &c, selectSQL, id)
}
// ... other CRUD methods
```

#### **handler.go (HTTP API):**
```go
type CandidateHandlers struct {
    svc *CandidateService
}

func (h *CandidateHandlers) Create(c *fiber.Ctx) error { ... }    // POST
func (h *CandidateHandlers) GetByID(c *fiber.Ctx) error { ... }   // GET /:id
func (h *CandidateHandlers) List(c *fiber.Ctx) error { ... }      // GET with pagination
func (h *CandidateHandlers) Update(c *fiber.Ctx) error { ... }    // PUT /:id
func (h *CandidateHandlers) Delete(c *fiber.Ctx) error { ... }    // DELETE /:id

func (h *CandidateHandlers) RegisterRoutes(app *fiber.App) {
    api := app.Group("/api/v1")
    api.Post("/candidates", h.Create)
    api.Get("/candidates", h.List)
    api.Get("/candidates/:id", h.GetByID)
    api.Put("/candidates/:id", h.Update)
    api.Delete("/candidates/:id", h.Delete)
}
```

#### **container.go (Domain DI):**
```go
type CandidateContainer struct {
    repo      CandidateRepository
    service   *CandidateService
    handlers  *CandidateHandlers
}

func NewCandidateContainer(db *sqlx.DB) *CandidateContainer {
    repo := &PostgresCandidateRepository{db: db}
    svc := &CandidateService{repo: repo}
    hdl := &CandidateHandlers{svc: svc}
    return &CandidateContainer{repo, svc, hdl}
}

func (c *CandidateContainer) RegisterRoutes(app *fiber.App) {
    c.handlers.RegisterRoutes(app)
}
```

---

### Wiring System (Marker-Based Code Injection)

All module wiring uses comment markers for **idempotent** code injection.

**Target Files & Markers:**

1. **pkg/config/config.go**
   ```go
   // manifesto:config-fields     ← Inject config struct fields
   // manifesto:config-loads      ← Inject Load() assignments
   ```

2. **cmd/container.go**
   ```go
   // manifesto:container-imports         ← Package imports
   // manifesto:container-fields          ← Container struct fields
   // manifesto:module-init               ← Module initialization code
   // manifesto:background-start          ← Background service startup
   // manifesto:container-helpers         ← Helper methods (Cleanup, etc.)
   ```

3. **cmd/server.go**
   ```go
   // manifesto:server-imports        ← Package imports
   // manifesto:public-routes         ← Public (unauthenticated) routes
   // manifesto:route-registration    ← Protected routes
   // manifesto:auth-middleware       ← Auth/scope middleware chains
   ```

4. **Makefile**
   ```makefile
   # manifesto:env-config       ← Environment variable documentation
   # manifesto:env-display      ← make env target variables
   ```

**Example: Wiring jobx**
```
Before:
  pkg/config/config.go:  // manifesto:config-fields
  cmd/container.go:      // manifesto:module-init

After:
  pkg/config/config.go:
    Jobx JobxConfig  // injected
  cmd/container.go:
    jobs.Dispatcher // injected field
    c.initJobx()    // injected call
```

---

### Remote Integration (GitHub)

**Repository:** `Abraxas-365/manifesto` (hardcoded)

**Download Strategy:**
1. Try GitHub API tag archive: `https://github.com/Abraxas-365/manifesto/archive/refs/tags/{version}.tar.gz`
2. Fallback to branch archive: `https://github.com/Abraxas-365/manifesto/archive/refs/heads/{version}.tar.gz`
3. Default ref: `main` branch

**Process:**
1. Download full repo archive as `.tar.gz`
2. Extract specific paths (e.g., `internal/fsx`, `internal/jobx`)
3. Rewrite Go imports: `github.com/Abraxas-365/manifesto` → project module path
4. Write to project's `pkg/` directory
5. Run `go get` to fetch new dependencies

**Version Pinning:**
- `--ref <tag|branch>` in `manifesto init` or `manifesto add`
- Stored in `manifesto.yaml` for reproducibility

---

### Config & Server Setup

#### **manifesto.yaml**
```yaml
project_name: myapp
go_module: github.com/user/myapp
manifesto_version: v1.2.0  # pinned version

modules:
  core:
    - kernel
    - errx
    - logx
    - config
    - ptrx

  wired:
    - fsx
    - jobx
    - notifx
    - iam

  timestamp: 2024-02-26T15:30:00Z
```

#### **Environment Variables (auto-loaded in config.go)**

**Server:**
```
SERVER_PORT=3000
SERVER_ENVIRONMENT=development
SERVER_LOG_LEVEL=debug
```

**Database:**
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=myapp
DB_SSLMODE=disable
```

**Redis:**
```
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

**If fsx wired:**
```
STORAGE_MODE=local|s3
UPLOAD_DIR=/uploads
AWS_REGION=us-east-1
AWS_BUCKET=my-bucket
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
```

**If jobx wired:**
```
JOBX_CONCURRENCY=10
JOBX_QUEUES=default,email,background
JOBX_POLL_INTERVAL=1s
```

**If notifx wired:**
```
NOTIFX_PROVIDER=ses|console
NOTIFX_FROM_ADDRESS=noreply@myapp.com
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
```

**If iam wired:**
```
JWT_SECRET_KEY=...
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=7d

OAUTH_GOOGLE_CLIENT_ID=...
OAUTH_GOOGLE_CLIENT_SECRET=...

OAUTH_MICROSOFT_CLIENT_ID=...
OAUTH_MICROSOFT_CLIENT_SECRET=...

OTP_LENGTH=6
OTP_EXPIRATION=5m

SESSION_TTL=7d
REFRESH_TOKEN_ROTATION=enabled
```

---

### README.md Contents

**Quick Summary:**
- Scaffold production-grade Go apps with DDD architecture
- Companion tool to Manifesto Architecture (github.com/Abraxas-365/manifesto)

**Install:**
```bash
go install github.com/Abraxas-365/manifesto-cli/cmd/manifesto@latest
```
Requires Go 1.24+.

**What You Get:**
- **Fiber HTTP server** with error handling, CORS, request IDs, graceful shutdown
- **Dependency injection** wired and ready
- **PostgreSQL + Redis** with pooling, health checks
- **Docker Compose** + **Makefile** (40+ commands)
- **Structured logging** (colored + JSON)
- **Rich error handling** (errx) with HTTP mapping
- **DDD architecture** enforced by templates

**Core Libraries (Always):**
| Name | Purpose |
|------|---------|
| kernel | Domain primitives, typed IDs, pagination, auth context |
| errx | Structured errors, HTTP mapping, code registry |
| logx | Colored/JSON logging, formatters |
| ptrx | Pointer utilities |
| config | Env-driven config |

**Optional Modules (Add on Demand):**
| Name | Purpose |
|------|---------|
| fsx | File storage: local disk or AWS S3 |
| asyncx | Futures, fan-out, pools, retry, timeout, rate-limiting |
| ai | LLM clients, embeddings, vector store, OCR, speech |
| jobx | Redis-backed async job queue |
| notifx | Email: AWS SES or console (debug) |
| iam | OAuth2, OTP, JWT, API keys, RBAC, multi-tenant users, sessions, invitations |

**Auto-Dependency Resolution:**
- `manifesto add jobx` → downloads asyncx + jobx
- `manifesto add ai` → downloads fsx + ai
- Cross-module bridges auto-inject (jobx + notifx → email handler)

**Generated Project Structure:**
```
myapp/
├── cmd/
│   ├── container.go        # DI wiring
│   └── server.go           # Fiber, middleware, handlers
├── pkg/
│   ├── kernel/             # Shared types
│   ├── errx/               # Error handling
│   ├── logx/               # Logging
│   ├── config/             # Configuration
│   └── ptrx/               # Pointer utilities
│   └── (module packages added on demand)
├── migrations/             # SQL migrations (if iam)
├── docker-compose.yml
├── Makefile
└── manifesto.yaml
```

**Commands Reference:**
```
manifesto init <name> --module <go-module>
manifesto init <name> --module <go-module> --with fsx,jobx,iam
manifesto init <name> --module <go-module> --all
manifesto init <name> --module <go-module> --quick

manifesto add <module>      # Wire: fsx, asyncx, ai, jobx, notifx, iam
manifesto add <path>        # Scaffold domain: pkg/foo/bar

manifesto modules           # List available + wire status
manifesto version           # Show CLI version
```

---

## Summary: How the Two Repos Work Together

### **Manifesto Library** (`/Users/abraxas/manifesto`)
- **Purpose**: Provide reusable, production-grade Go packages for building web services
- **Scope**: 12 libraries across 5000+ lines of Go
  - **Foundational**: kernel, errx, logx, ptrx, config
  - **Service-layer**: asyncx, fsx, jobx, notifx
  - **Feature-complete**: ai (LLM + embeddings), iam (full auth system)
- **Deployment**: Published as Go module `github.com/Abraxas-365/manifesto`
- **Types**: All types, interfaces, implementations are in internal packages

### **Manifesto CLI** (`/Users/abraxas/manifesto-cli`)
- **Purpose**: Scaffold new Go projects and wire modules
- **Scope**: 16 .go files (3000+ LOC)
- **Core features**:
  1. **Project scaffolding**: `manifesto init` creates a production-ready DDD project
  2. **Module wiring**: `manifesto add <module>` downloads from upstream + injects code
  3. **Domain generation**: `manifesto add <path>` creates 7-file DDD verticals
  4. **Dependency resolution**: Auto-downloads transitive deps (jobx → asyncx)
  5. **Code injection**: Marker-based system (idempotent) for wiring into config/container/server/makefile
  6. **GitHub integration**: Downloads manifesto library, rewrites imports, embeds in binary
- **Deployment**: Single binary via GoReleaser

### **Workflow:**
```
User: manifesto init myapp --module github.com/user/myapp --with fsx,jobx
      ↓
CLI downloads: kernel, errx, logx, ptrx, config, fsx, asyncx, jobx
              from github.com/Abraxas-365/manifesto@main
      ↓
CLI rewrites imports & generates:
  - DI container (cmd/container.go)
  - HTTP server (cmd/server.go)
  - Config loader (pkg/config/config.go with markers)
  - Makefile (40+ commands)
  - docker-compose.yml
  - manifesto.yaml
      ↓
Project ready to develop with:
  - `make dev` — start server
  - `manifesto add pkg/recruitment/candidate` — add domain
  - `manifesto add iam` — add IAM system
      ↓
All types from manifesto library are available in pkg/
```

---

## Key Insights

### **Manifesto Library Strengths:**
1. ✅ **Rich, reusable packages** — Each package solves a specific problem (logging, errors, async, file storage, jobs, email, LLM, IAM)
2. ✅ **Composable interfaces** — All storage/messaging abstractions allow swapping implementations (FileSystem: local ↔ S3)
3. ✅ **Production-ready** — Includes error mapping, structured logging, retries, rate limiting, graceful shutdown
4. ✅ **LLM-first** — Dedicated ai package with multi-provider support (OpenAI, Anthropic, Bedrock, Gemini, Azure)
5. ✅ **Full IAM** — Complete auth system (OAuth2, OTP, JWT, API keys, RBAC, multi-tenancy)

### **Manifesto CLI Strengths:**
1. ✅ **Opinionated DDD** — Enforces layered architecture (handler → service → domain → repo → infra)
2. ✅ **Smart scaffolding** — Generates 7-file domain verticals with auto-injection into root container/server
3. ✅ **Marker-based wiring** — Idempotent code injection allows running `manifesto add jobx` multiple times safely
4. ✅ **Dep resolution** — Recursive dependency solver (jobx → asyncx)
5. ✅ **Cross-module bridges** — jobx + notifx auto-inject email handler
6. ✅ **Makefile as first-class** — Generated Makefile includes 40+ dev/test/build/db commands
7. ✅ **GitHub integration** — Downloads library, rewrites imports, self-contained binary

### **Gap Analysis (What's Missing from CLI):**

1. **Additional module scaffolding**:
   - `manifesto add service <name>` — Generate service-only layer (no handlers/repo)
   - `manifesto add migration <name>` — Create new SQL migration files
   - `manifesto add middleware <name>` — Generate custom Fiber middleware templates

2. **Enhanced configuration**:
   - `manifesto config set KEY VALUE` — Persist config to manifesto.yaml
   - `manifesto config validate` — Validate manifesto.yaml against schema
   - `manifesto config diff` — Show what would change before applying

3. **Code generation improvements**:
   - `manifesto generate swagger` — Auto-generate OpenAPI/Swagger from domain handlers
   - `manifesto generate tests` — Generate unit tests for domain templates
   - `manifesto generate migrations` — Auto-generate SQL migrations from iam module

4. **Project maintenance**:
   - `manifesto lint` — Check project structure compliance with DDD rules
   - `manifesto upgrade` — Update manifesto library version in existing projects
   - `manifesto clean` — Remove unused generated code

5. **Multi-repo support**:
   - `--repo <url>` flag — Allow custom manifesto repository (instead of hardcoded Abraxas-365/manifesto)
   - `manifesto registry` — List custom/community module registries

6. **Better module discovery**:
   - `manifesto search <query>` — Search for modules by name/description
   - `manifesto module info <name>` — Show detailed module documentation
   - `manifesto module deps <name>` — Show transitive dependency graph

7. **Live validation**:
   - `manifesto validate` — Check generated project for Go build/lint errors
   - `manifesto test` — Run tests in generated project
   - `manifesto check-compat` — Verify module compatibility before wiring

These improvements would make the CLI more powerful for project lifecycle management while maintaining the current simplicity of core commands.

