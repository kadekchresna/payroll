FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o payroll-service .

# ----------------------------------------

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache curl ca-certificates && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz -o migrate.tar.gz && \
    tar -xzf migrate.tar.gz && \
    mv migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate && \
    rm migrate.tar.gz


COPY --from=builder /app/payroll-service .
COPY --from=builder /app/migration ./migration

RUN chmod +x payroll-service

EXPOSE 8081

CMD migrate -path ./migration -database postgres://postgres:secret@postgres:5432/payroll_db?sslmode=disable up && \
    ./payroll-service web
