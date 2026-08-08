# ComeToSee Backend Architecture

This document expands the diagrams in the root [README](../README.md). It is based on the current implementation in `main.go`, `routes/`, `controller/`, `service/`, `repository/`, `middleware/`, and `intailizer/migrations/`.

## Component boundaries

```mermaid
flowchart LR
    subgraph Edge[HTTP edge]
        CORS[CORS]
        Limit[Per-process rate limiter]
        Recover[Panic recovery]
        Mux[Gorilla Mux]
    end
    subgraph App[Application]
        Routes[Route declarations]
        Auth[JWT middleware]
        Admin[Admin JWT middleware]
        Controllers[Controllers]
        Services[Services]
        Repos[Repositories]
    end
    subgraph State[State and persistence]
        Cache[go-cache]
        Manager[WebSocket manager]
        PG[(PostgreSQL/PostGIS)]
    end
    subgraph Providers[External providers]
        Mail[SMTP]
        Payment[eSewa]
        RTC[Agora]
        Firebase[Firebase initialization]
    end

    CORS --> Limit --> Recover --> Mux --> Routes
    Routes --> Auth --> Controllers
    Routes --> Admin --> Controllers
    Routes --> Controllers
    Controllers --> Services --> Repos --> PG
    Services --> Cache
    Controllers --> Manager --> PG
    Services --> Mail
    Controllers --> Payment
    Services --> RTC
    App -. startup .-> Firebase
```

## WebSocket event routing

```mermaid
flowchart TD
    Upgrade[GET /ws/manager] --> JWT[JWT middleware]
    JWT --> Handler[WSController.ServeWS]
    Handler --> Read[Read loop]
    Handler --> Write[Write loop]
    Read --> Route[WebsocketService.RouteEvent]
    Route --> Register[register]
    Route --> Room[change room]
    Route --> Send[send message]
    Route --> Invite[video call invite]
    Route --> History[get history]
    Register --> Manager[(Manager maps: clients, users, rooms)]
    Room --> Manager
    Send --> MessageRepo[(message repository)]
    Send --> Manager
    Invite --> Manager
    History --> MessageRepo
    MessageRepo --> DB[(PostgreSQL)]
    Write --> Client[WebSocket client]
```

## Authentication sequence

```mermaid
sequenceDiagram
    participant Client
    participant Middleware
    participant Handler
    participant Service

    Client->>Middleware: Authorization: Bearer <JWT>
    Middleware->>Middleware: Parse HMAC JWT with SECRET
    Middleware->>Middleware: Require email, authId and username claims
    Middleware->>Handler: Attach claims to request context
    Handler->>Service: Read authId/email/username from context
    Service-->>Handler: Business result
    Handler-->>Client: JSON response
```

Admin middleware follows the same pattern with `ADMIN_SECRET` fallback to `SECRET`, then requires `role == "admin"` and attaches `adminId` and `adminEmail`.

## Data flow for geographic discovery

```mermaid
sequenceDiagram
    participant User
    participant API
    participant Discovery as Discovery/map service
    participant Repo as Location repository
    participant PG as PostgreSQL + PostGIS

    User->>API: Save or update location
    API->>Discovery: Pass authenticated location payload
    Discovery->>Repo: Persist location
    Repo->>PG: Write geography(Point, 4326)
    PG-->>Repo: Saved location
    Repo-->>API: Success
    API-->>User: JSON response
    User->>API: Request map pins or discovery feed
    API->>Repo: Query nearby users/events
    Repo->>PG: Spatial query using geom/index
    PG-->>Repo: Nearby records
    Repo-->>API: Pins/users
    API-->>User: JSON response
```

## Notes for future changes

- Keep route protection explicit and document any newly public endpoint.
- Prefer typed context keys for JWT claims to avoid collisions.
- Move process-local OTP, cache, rate-limit, and WebSocket coordination to shared infrastructure when scaling beyond one instance.
- Keep credentials outside the repository and rotate exposed values before deployment.
