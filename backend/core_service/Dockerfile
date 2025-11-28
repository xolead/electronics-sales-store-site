FROM golang:1.25.1 AS builder

WORKDIR /app


COPY . .

RUN go build -o core .

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/core /app/core
COPY --from=builder /app/pkg/dbwork/migrations /app/pkg/dbwork/migrations
EXPOSE 8082

CMD ["/app/core"]

