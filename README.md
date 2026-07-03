# Workout Journal

Workout tracking app: REST API built with Go and PostgreSQL, with a web frontend.

**Live app:** https://longtermgames.github.io/workout-api
**API:** https://workout-api-production-c1f1.up.railway.app

## Features

- JWT authentication (register, login) with bcrypt password hashing
- Each user sees only their own workouts
- Workout journal grouped by day (Today / Yesterday / date)
- Per-exercise progression chart
- Exercise name autocomplete
- CRUD operations for workouts
- CORS middleware for frontend access

## Technologies

- Go (net/http, standard library)
- PostgreSQL
- JWT (github.com/golang-jwt/jwt)
- bcrypt
- Railway (API deployment)
- GitHub Pages (frontend deployment)

## Run locally

1. Start PostgreSQL
2. Create database and tables:

```sql
CREATE DATABASE workout_db;
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);
CREATE TABLE workouts (
  id SERIAL PRIMARY KEY,
  exercise TEXT NOT NULL,
  reps INT NOT NULL,
  user_id INTEGER REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW()
);
```

3. Run the server:

```bash
DATABASE_URL="user=youruser dbname=workout_db sslmode=disable" \
JWT_SECRET="your-secret-key" \
go run main.go
```

## API Endpoints

All workout endpoints require `Authorization: Bearer <token>` header.

### Register
```http
POST /register
```
```json
{ "username": "tengiz", "password": "12345" }
```

### Login
```http
POST /login
```
Returns:
```json
{ "token": "eyJhbGci..." }
```

### Get all workouts
```http
GET /workouts
```
Returns workouts with creation dates:
```json
[{ "id": 1, "exercise": "Push Ups", "reps": 20, "created_at": "2026-07-04T10:00:00Z" }]
```

### Get workout by ID
```http
GET /workout?id=1
```

### Create workout
```http
POST /workouts
```
```json
{ "exercise": "Push Ups", "reps": 20 }
```

### Update workout
```http
PUT /workout?id=1
```
```json
{ "exercise": "Push Ups", "reps": 25 }
```

### Delete workout
```http
DELETE /workout?id=1
```
