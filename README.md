# Training Management API

REST API for workout management built with Go and PostgreSQL.

**Live:** https://workout-api-production-c1f1.up.railway.app

## Features

- JWT authentication (register, login)
- Create workouts
- Get all workouts
- Get workout by ID
- Update workouts
- Delete workouts
- Each user sees only their own workouts

## Technologies

- Go
- PostgreSQL
- JWT (github.com/golang-jwt/jwt)
- bcrypt password hashing
- Railway (deployment)

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

## Endpoints

### Register

```http
POST /register
```

Body:

```json
{
  "username": "tengiz",
  "password": "12345"
}
```

### Login

```http
POST /login
```

Returns:

```json
{
  "token": "eyJhbGci..."
}
```

### Get all workouts

```http
GET /workouts
Authorization: Bearer <token>
```

### Get workout by ID

```http
GET /workout?id=1
Authorization: Bearer <token>
```

### Create workout

```http
POST /workouts
Authorization: Bearer <token>
```

Body:

```json
{
  "exercise": "Push Ups",
  "reps": 20
}
```

### Update workout

```http
PUT /workout?id=1
Authorization: Bearer <token>
```

Body:

```json
{
  "exercise": "Push Ups",
  "reps": 25
}
```

### Delete workout

```http
DELETE /workout?id=1
Authorization: Bearer <token>
```
