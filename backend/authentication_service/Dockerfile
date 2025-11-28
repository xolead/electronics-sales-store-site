FROM golang:1.25.3 AS builder

WORKDIR /app


COPY . .

RUN go build -o auth .

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/auth /app/auth
COPY --from=builder /app/pkg/dbwork/migrations /app/pkg/dbwork/migrations
EXPOSE 8081

CMD ["/app/auth"]
