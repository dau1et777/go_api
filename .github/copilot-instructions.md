# Copilot / AI agent instructions for this repository

This repo is a small monolithic Go HTTP API. Below are the concrete facts and patterns an AI assistant should follow to be immediately productive.

- Entry point: `cmd/server/main.go` — loads env, calls `database.Connect()` and `routes.RegisterRoutes()` and starts the HTTP server. Prefer edits here only for bootstrapping changes.
- Router: `internal/routes/routes.go` — registers endpoints with `github.com/gorilla/mux`. Example: `r.HandleFunc("/user", handler.CreateUser).Methods("POST")`.
- HTTP handlers: `internal/handler/*.go` — thin layer that currently performs request parsing, validation and direct DB calls (see `CreateUser` in `internal/handler/user_handler.go`). When adding features prefer using the `service` layer instead of direct DB access.
- Service layer: `internal/service/user_service.go` — contains business logic and interacts with `internal/repository`. Prefer adding/using methods here for validation, uniqueness checks, and orchestrating repo calls.
- Repository layer: `internal/repository/user_repository.go` — raw SQL access to Postgres (uses `github.com/lib/pq`). Note: `GetByEmail` returns the user (without password) and the stored password hash separately — handle that carefully when implementing auth.
- Models: `internal/models/user.go` — model struct and validation helpers `ValidateForCreate()` and `ValidateForUpdate()`. Use these methods where appropriate instead of duplicating validation.
- Database: `internal/database/database.go` — global `var DB *sql.DB` and `Connect()` which reads DB_* env vars and sets `DB`. Codebase currently relies on this global DB variable.
- Config: `internal/config/config.go` — wraps `godotenv` and reads `APP_PORT` and `JWT_SECRET`. Environment variables used across the project: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`, `APP_PORT`, `JWT_SECRET`.
- Auth: `cmd/server/auth/jwt.go` — helper to generate JWT tokens; the project uses `github.com/golang-jwt/jwt/v5`.
- Logger: `internal/logger/logger.go` — provides `InitLogger()` and a `Logger` variable but it's not wired consistently across `main.go` (be mindful when adding logs).

Conventions & patterns (follow these exactly):
- Layering: route -> handler -> service -> repository -> database. Prefer adding new logic in `service` (not in handlers) and SQL in `repository`.
- Validation: prefer `models.User.ValidateForCreate()` / `ValidateForUpdate()` for user validation. Do not duplicate regexes or error messages; reuse those functions.
- Error handling: return Go `error` values from service/repo functions. Handlers translate errors to HTTP responses. Follow existing error message styles (short, user-facing strings in handlers).
- DB interactions: repository methods return domain objects (often without password) and explicit errors. `GetByEmail` returns `(user, passwordHash, error)` — caller must compare the hash.
- Routing: add new endpoints in `internal/routes/routes.go` using Gorilla Mux. Keep route registration centralized in this file.

Build, run, and debug commands (Windows PowerShell):
- Build: `go build ./...` from repository root.
- Run locally: set `.env` and then run `go run ./cmd/server` or run the built binary. `cmd/server/main.go` expects `.env` with DB_*/APP_PORT/JWT_SECRET.

Dependencies of note (see `go.mod`): `github.com/gorilla/mux`, `github.com/lib/pq`, `github.com/joho/godotenv`, `github.com/golang-jwt/jwt/v5`.

Typical task checklist when adding an API endpoint (minimal, follow this order):
1. Add route to `internal/routes/routes.go`.
2. Add handler in `internal/handler` that parses request and calls the appropriate `service` method.
3. Add business logic to `internal/service` (validation, uniqueness checks, orchestration).
4. Add SQL operations to `internal/repository` and update models if needed.
5. Update or add model validation in `internal/models` if new fields are introduced.
6. Add or update any environment/config usage in `internal/config`.
7. Run `go build ./...` and test endpoints.

Quick examples from this repo:
- Use `User.ValidateForCreate()` before persisting a new user (see `internal/models/user.go`).
- To check for existing email, the service calls `s.Repo.GetByEmail(u.Email)` and checks return values (see `internal/service/user_service.go`).
- `internal/handler.CreateUser` currently writes to DB via `database.DB.Exec(...)`; prefer calling `service.CreateUser` instead for new work.

Notes for the AI agent:
- Do not change global patterns unless the change is small and consistent (e.g., wiring `logger.InitLogger()` in `main.go` is OK). If you propose architectural changes (like removing `database.DB` global), state the migration plan and update all call sites.
- Be conservative with error messages and HTTP status codes — follow existing usage (`http.StatusBadRequest`, `http.StatusCreated`, `http.StatusInternalServerError`).
- When editing SQL, keep Postgres syntax compatible with `lib/pq` and use parameterized queries (already used with `$1`).

If any behavior is ambiguous, ask for clarification and include the exact file(s) you plan to change.

---
If you'd like, I can now merge this into the repo and run `go build ./...` to verify. Tell me to proceed or request adjustments to the draft.
