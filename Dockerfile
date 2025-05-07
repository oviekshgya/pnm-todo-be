FROM golang:1.23-alpine AS builder

#RUN apk add --no-cache git
RUN apk add --no-cache mysql-client tzdata

ENV TZ=Asia/Jakarta
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/pnm-todo-be .

# Runtime stage
FROM gcr.io/distroless/base-debian12:latest

WORKDIR /app

COPY --from=builder /app/pnm-todo-be .
COPY --from=builder /app/.env .

ENTRYPOINT ["/app/pnm-todo-be"]
