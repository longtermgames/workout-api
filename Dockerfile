# --- build stage ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# copy dependency files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of the source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# --- run stage ---
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
