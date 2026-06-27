# Training Management API

REST API for workout management built with Go and PostgreSQL, with a frontend UI.

**Live API:** https://workout-api-production-c1f1.up.railway.app  
**Live Frontend:** https://longtermgames.github.io/workout-api

## Features

- JWT authentication (register, login)
- bcrypt password hashing
- CRUD operations for workouts
- Each user sees only their own workouts
- CORS middleware for frontend access
- Frontend UI (HTML/JS)

## Technologies

- Go
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
CREATE TABLE users (id SERIAL PRIMARY KEY, username TEXT NOT NULL UNIQUE, password TEXT NOT NULL);
CREATE TABLE workouts (id SERIAL PRIMARY KEY, exercise TEXT NOT NULL, reps INT NOT NULL, user_id INTEGER REFERENCES users(id));
```

3. Run the server:

```bash
DATABASE_URL="user=youruser dbname=workout_db sslmode=disable" go run main.go
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
