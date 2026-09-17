# Workout Journal

Workout tracking app: a REST API in Go (standard library only) backed by PostgreSQL, with a small web frontend.

**Live app:** https://longtermgames.github.io/workout-api
**API:** https://workout-api-production-c1f1.up.railway.app

## Features

- JWT authentication (register, login) with bcrypt password hashing
- Each user sees only their own workouts (every query is scoped by `user_id`)
- Input validation and correct HTTP status codes (400 / 401 / 404 / 409)
- Workout journal grouped by day (Today / Yesterday / date)
- Per-exercise progression chart + exercise-name autocomplete
- CRUD for workouts; CORS middleware for the frontend

## Tech

- Go (`net/http`, method-based routing from Go 1.22)
- PostgreSQL (`github.com/lib/pq`)
- JWT (`github.com/golang-jwt/jwt/v5`), bcrypt (`golang.org/x/crypto`)
- Railway (API) + GitHub Pages (frontend)

## Architecture

```
main.go        server setup, routing, CORS, timeouts, config from env
auth/          JWT issue/verify, password hashing
db/            data access (queries) + connection pool
handlers/      HTTP handlers, one per method+route, request validation
```

## Design notes

- **Secrets from the environment, fail-fast.** The server refuses to start without `JWT_SECRET` instead of silently signing tokens with an empty key.
- **No detail leaks to clients.** Internal errors are logged server-side; the client gets a generic 500. Login returns the same message for an unknown user and a wrong password (no user enumeration).
- **JWT signing method is verified** on parse to guard against `alg` confusion.
- **Server timeouts** (read/write/idle) protect against slow clients.
- **Pooled DB connections** and an index on `workouts(user_id)`.

## Run locally

1. Start PostgreSQL and create the database + tables:

```bash
createdb workout_db
psql workout_db -f schema.sql
```

2. Copy env and fill it in:

```bash
cp .env.example .env
```

3. Run the server (`JWT_SECRET` is required):

```bash
DATABASE_URL="user=youruser dbname=workout_db sslmode=disable" \
JWT_SECRET="$(openssl rand -base64 32)" \
go run .
```

## Tests

```bash
go test ./...
```

## API

All `/workouts` and `/workout` endpoints require `Authorization: Bearer <token>`.

| Method | Path             | Body                          | Notes                    |
|--------|------------------|-------------------------------|--------------------------|
| POST   | `/register`      | `{username, password}`        | 201, 409 if taken        |
| POST   | `/login`         | `{username, password}`        | returns `{token}`        |
| GET    | `/workouts`      | —                             | list, newest first       |
| POST   | `/workouts`      | `{exercise, reps}`            | 201                      |
| GET    | `/workout?id=1`  | —                             | single                   |
| PUT    | `/workout?id=1`  | `{exercise, reps}`            | 404 if not yours         |
| DELETE | `/workout?id=1`  | —                             | 404 if not yours         |
