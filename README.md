# ComeToSee Backend

ComeToSee is a Go backend for a social/event-discovery application. It exposes JSON HTTP APIs for authentication, profiles, posts, connections, subscriptions, payments, verification, maps, QR validation, and discovery. It also exposes a JWT-protected WebSocket endpoint for chat and call invitations.

## Technology stack

- Go 1.25.7
- Gorilla Mux for HTTP routing
- PostgreSQL via `database/sql` and `lib/pq`
- Goose embedded SQL migrations
- JWT authentication with separate user and admin middleware
- Gorilla WebSocket for real-time messaging
- PostGIS for geographic queries
- eSewa payment integration
- Agora RTC token generation
- Firebase Admin initialization / Firestore support
- In-memory cache, Prometheus metrics, and `net/http/pprof`

## Repository layout

| Directory | Responsibility |
| --- | --- |
| `main.go` | Startup, middleware, dependency injection, route registration, metrics and profiling servers |
| `routes/` | Gorilla Mux endpoint declarations |
| `controller/` | HTTP/WebSocket request handlers and response serialization |
| `service/` | Application/business logic |
| `repository/` | PostgreSQL queries and persistence abstractions |
| `model/` | Request, response, database and WebSocket models |
| `di/` | Manual controller-service-repository wiring |
| `middleware/` | JWT, admin JWT, rate limiting and panic recovery |
| `intailizer/` | Environment loading, database connection and Goose migrations |
| `intailizer/migrations/` | Ordered PostgreSQL/PostGIS schema migrations |
| `firebase/` | Firebase application and Firestore client initialization |
| `features/` | External feature adapters such as Agora token generation |
| `common/` | Shared context helpers, JSON responses, mail and room-id helpers |
| `config/` | Cache and garbage-collector configuration |
| `test/` | Integration-style tests |

## Architecture

```mermaid
flowchart TB
    Client[Web or mobile client]
    Admin[Admin client]
    Router[Gorilla Mux router]
    HTTP[HTTP middleware\nCORS -> rate limit -> panic recovery]
    JWT[JWT middleware\nuser claims in context]
    AdminJWT[Admin JWT middleware\nrole=admin]
    Controllers[Controllers\nHTTP and WebSocket handlers]
    Services[Services\nbusiness rules and orchestration]
    Repositories[Repositories\nSQL persistence]
    DB[(PostgreSQL + PostGIS)]
    WS[In-memory WebSocket manager\nclients, users, rooms, handlers]
    Cache[(In-memory cache)]
    Firebase[Firebase Admin / Firestore]
    SMTP[SMTP Gmail]
    Esewa[eSewa]
    Agora[Agora RTC]

    Client --> Router
    Admin --> Router
    Router --> HTTP
    HTTP --> JWT
    HTTP --> AdminJWT
    HTTP --> Controllers
    JWT --> Controllers
    AdminJWT --> Controllers
    Controllers --> Services
    Services --> Repositories
    Repositories --> DB
    Services --> Cache
    Controllers --> WS
    WS --> DB
    Services --> SMTP
    Controllers --> Esewa
    Services --> Agora
    Controllers -. initialized at startup .-> Firebase
```

The normal synchronous request path is `route -> middleware (where configured) -> controller -> service -> repository -> PostgreSQL -> response`. Dependency injection is manual and happens in `main.go` through the constructors in `di/`.

## Startup and runtime flow

```mermaid
sequenceDiagram
    participant Process as Go process
    participant Env as .env / environment
    participant Cache as Cache
    participant DB as PostgreSQL
    participant Goose as Goose migrations
    participant Firebase as Firebase Admin
    participant Router as API server
    participant Metrics as Metrics :2112
    participant Pprof as pprof :6060

    Process->>Env: Load environment variables
    Process->>Cache: Initialize 5-minute cache
    Process->>DB: Open pooled connection and Ping
    Process->>Goose: Apply embedded migrations
    Goose->>DB: Create/update schema
    Process->>Firebase: Initialize app from credentials and project ID
    Process->>Router: Build controllers and register routes
    Process->>Metrics: Start Prometheus server
    Process->>Pprof: Start localhost profiling server
    Process->>Router: Listen on :PORT
```

## HTTP request flow

```mermaid
sequenceDiagram
    participant C as Client
    participant CORS as CORS handler
    participant RL as Rate limiter
    participant Recovery as Panic recovery
    participant Auth as JWT / admin JWT
    participant Ctrl as Controller
    participant Svc as Service
    participant Repo as Repository
    participant DB as PostgreSQL

    C->>CORS: HTTP request
    CORS->>RL: Forward request
    RL->>Recovery: Forward if under limit
    Recovery->>Auth: Route-specific auth check
    Auth->>Ctrl: Request plus context claims
    Ctrl->>Svc: Validate and invoke use case
    Svc->>Repo: Read or mutate domain data
    Repo->>DB: SQL query
    DB-->>Repo: Rows or mutation result
    Repo-->>Svc: Domain result
    Svc-->>Ctrl: Result or error
    Ctrl-->>C: JSON response
```

Public routes are not automatically protected. Authentication is applied per route in `routes/`; verify a route before exposing it to untrusted clients. In particular, several connection, discovery, user-filter, requirement, and Agora routes currently do not wrap handlers with JWT middleware.

## Main user journeys

### Registration and login

```mermaid
sequenceDiagram
    participant U as User
    participant API as Auth controller
    participant S as Auth service
    participant DB as PostgreSQL
    participant Mail as SMTP

    U->>API: POST /signup
    API->>S: Validate registration data
    S->>S: Generate OTP and keep pending signup in memory
    S->>Mail: Send verification OTP
    API-->>U: Signup accepted
    U->>API: POST /verifyemail with OTP
    API->>S: Verify OTP
    S->>DB: Insert user into cometoseeauth
    API-->>U: Email verified
    U->>API: POST /login
    API->>S: Check password and account status
    S->>DB: Load user by email
    S-->>API: JWT with authId, email and username
    API-->>U: JWT
```

### Posts, slots and social interactions

```mermaid
sequenceDiagram
    participant U as Authenticated user
    participant API as Post controller
    participant S as Post service
    participant DB as PostgreSQL
    participant Mail as SMTP

    U->>API: POST /post, /post/like, /post/comment or /post/share
    API->>S: Extract authId from JWT context
    S->>DB: Create post or interaction
    DB-->>S: Persisted record
    S-->>API: Result
    API-->>U: JSON response
    U->>API: POST /createslot or /joinslot
    API->>S: Create/join event slot
    S->>DB: Update post_slots and slot_participants
    DB-->>S: Updated slot
    S-->>API: Slot result
    API-->>U: JSON response
    U->>API: POST /post/cancel
    S->>DB: Cancel post and load participant emails
    S->>Mail: Notify participants
```

### Real-time messaging

```mermaid
sequenceDiagram
    participant U1 as User A
    participant WS as /ws/manager
    participant M as In-memory manager
    participant S as WebSocket service
    participant DB as PostgreSQL
    participant U2 as User B

    U1->>WS: Upgrade with JWT
    WS->>M: Register client/user/room
    U1->>WS: send message event
    WS->>S: RouteEvent
    S->>DB: Save message
    S->>M: Enqueue new message for recipient/room
    M-->>U2: new message event
    U1->>WS: get history event
    S->>DB: Query paginated room history
    DB-->>S: Messages
    S-->>U1: history response event
```

### eSewa subscription payment

```mermaid
sequenceDiagram
    participant U as User
    participant API as Backend
    participant DB as PostgreSQL
    participant E as eSewa

    U->>API: POST /esewa/initiate with JWT
    API->>DB: Create pending payment
    API-->>U: Signed eSewa form/redirect data
    U->>E: Complete payment
    E->>API: GET /esewa/verify
    API->>API: Verify eSewa signature and transaction
    API->>DB: Mark payment confirmed
    API->>DB: Create or extend subscription
    API-->>E: Verification result
    E-->>U: Redirect to success URL
```

## Route inventory

| Area | Endpoints |
| --- | --- |
| Auth | `POST /signup`, `/login`, `/forgetpassword`, `/resetpassword`, `/verifyemail`; `GET /getprofile` |
| Admin | `POST /admin/login` |
| Posts | `POST /post`, `/post/like`, `/post/comment`, `/post/share`, `/getpost`, `/createslot`, `/joinslot`, `/post/cancel`; `GET /latestlike`, `/slot/participant`, `/chats/joined` |
| Connections | `/connection/send`, `/accept`, `/block`, `/get`, `/reject`, `/received`, `/sended`, `/unsend`, `/filter`, `/connectedpeople`, `/discoveredpeople` |
| Profiles and discovery | `/userdetailinfo`, `/userlocation`, `/profilestatus`, `/updateuserdetailinfo`, `/updatelocation`, `/fetchuserprofile`, `GET /userfilter`, `GET /feed/dicovery`, `GET /map/pins` |
| Subscription | `POST /delete`; `GET /status` |
| Payments | `POST /esewa/initiate`; `GET /esewa/verify`, `/esewa/failure` |
| Verification | User upload/read endpoints under `/verification`; admin review endpoints under `/admin/verification` |
| Requirements and QR | CRUD under `/requirements`; `/qr/joined`, `/qr/my-posts`, `/qr/verify` |
| Calls | `POST /createcall`, `/startcall`, `/endcall` |
| WebSocket | `GET /ws/manager` |

See `routes/` for the authoritative HTTP methods and middleware assignments.

## Database model

```mermaid
erDiagram
    cometoseeauth ||--o{ messages : sends
    cometoseeauth ||--o{ subscriptions : owns
    cometoseeauth ||--o{ payments : makes
    cometoseeauth ||--o{ connections : participates
    cometoseeauth ||--o| userdetailinfo : has
    userdetailinfo ||--o{ location : has
    cometoseeauth ||--o{ post : creates
    post ||--o{ comments : receives
    post ||--o{ post_likes : receives
    post ||--o{ post_shares : receives
    post ||--o| post_slots : schedules
    post_slots ||--o{ slot_participants : contains
    cometoseeauth ||--o{ slot_participants : joins
    cometoseeauth ||--o| user_verification : submits
    cometoseeauth ||--o{ player_documents : uploads
    post ||--o{ post_requirements : defines
    cometoseeauth ||--o{ video_call_sessions : joins
    admin_auth ||--o{ user_verification : reviews

    cometoseeauth { int auth_id PK }
    post { int post_id PK int auth_id FK }
    userdetailinfo { int user_detail_id PK int auth_id FK }
    location { int location_id PK int user_detail_id FK geography geom }
    post_slots { int slot_id PK int post_id FK }
    slot_participants { int slot_id FK int auth_id FK }
```

## Configuration

Copy the required variable names into a local `.env` or deployment secret store. Never commit real credentials.

```dotenv
PORT=8081
DB_URL=postgresql://user:password@host:5432/database
SECRET=replace-user-jwt-secret
ADMIN_SECRET=replace-admin-jwt-secret
FIREBASE_PROJECT_ID=your-project-id
Location=path/to/serviceAccountkey.json
APP_PASSWORD=mail-app-password
Agora_APP_ID=your-agora-app-id
Agora_APP_CERTIFICATE=your-agora-certificate
ESEWA_SECRET=your-esewa-secret
ESEWA_PRODUCT_CODE=your-product-code
ESEWA_SUCCESS_URL=http://localhost:8081/esewa/verify
ESEWA_FAILURE_URL=http://localhost:8081/esewa/failure
GOGC=100
MEM_LIMIT_MB=512
```

The current process also references `GOOGLE_APPLICATION_CREDENTIALS` in test/development setup. Firebase initialization reads `Location`.

## Run locally

Prerequisites: Go, a reachable PostgreSQL database with PostGIS available, and a Firebase service-account JSON file if Firebase initialization is enabled.

```powershell
go mod download
go test ./...
go run .
```

On startup the backend loads `.env`, opens PostgreSQL, applies embedded Goose migrations, initializes Firebase, starts the API on `:${PORT}`, Prometheus metrics on `:2112`, and pprof on `localhost:6060`.

## Operational and security notes

- Rotate any credentials that may have been committed in `.env` or `firebase/serviceAccountkey.json`; use a secret manager in deployment.
- Restrict pprof to trusted operators and do not expose it publicly.
- Replace the permissive WebSocket `CheckOrigin` implementation before production deployment.
- Review route-level JWT coverage before exposing the API publicly; middleware is not global.
- OTP pending users and WebSocket connection state are process-local and are lost on restart; horizontal scaling needs shared state or sticky routing.
- The rate limiter is process-local and keyed by `RemoteAddr`, so a multi-instance deployment needs a shared limiter or edge protection.
- Firebase is initialized during startup, but no production repository call to Firestore was found in this backend scan; PostgreSQL is the primary persistence store.

## Verification checklist

```powershell
go test ./...
go vet ./...
```

The migrations are embedded into the binary, so deployment does not require a separate migration directory at runtime.
