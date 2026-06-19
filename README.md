# Training Management API

REST API for workout management built with Go and PostgreSQL.

**Live:** https://workout-api-production-c1f1.up.railway.app

## Features

- Create workouts
- Get all workouts
- Get workout by ID
- Update workouts
- Delete workouts

## Technologies

- Go
- PostgreSQL
- Railway (deployment)

## Run locally

1. Start PostgreSQL
2. Create database:

```sql
CREATE DATABASE workout_db;
```

3. Run the server:

```bash
DATABASE_URL="user=youruser dbname=workout_db sslmode=disable" go run main.go
```

## Endpoints

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
```
