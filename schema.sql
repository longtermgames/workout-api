-- Workout Journal schema

CREATE TABLE IF NOT EXISTS users (
    id       SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS workouts (
    id         SERIAL PRIMARY KEY,
    exercise   TEXT NOT NULL,
    reps       INT NOT NULL CHECK (reps > 0),
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- A foreign key is not indexed automatically in Postgres; every /workouts query
-- filters by user_id, so this index keeps those lookups fast as data grows.
CREATE INDEX IF NOT EXISTS idx_workouts_user_id ON workouts (user_id);
