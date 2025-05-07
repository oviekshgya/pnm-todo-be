# Stage 1: Build Stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum from root context
COPY go.mod go.sum ./
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/pnm-todo-be .
# Copy all project files from root (because build context is root)
COPY . .

#COPY ./pnm-todo-be /app/pnm-todo-be

# Stage 2: Runtime Stage
FROM gcr.io/distroless/base-debian12:latest

WORKDIR /app

COPY --from=builder /app/pnm-todo-be .
COPY --from=builder /app/.env .

ENTRYPOINT ["/app/pnm-todo-be"]
