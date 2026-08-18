FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o monitor-bot ./cmd/main.go


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/monitor-bot .

RUN mkdir -p /app/storage

CMD ["./monitor-bot"]