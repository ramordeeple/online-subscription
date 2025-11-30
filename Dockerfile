FROM golang:1.25.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o online-subscription ./cmd/online-subscription

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/online-subscription .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations /app/migrations

CMD ["./online-subscription"]
