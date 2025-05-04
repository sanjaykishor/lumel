FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lumel ./cmd/lumel

FROM alpine:latest

RUN apk add --no-cache postgresql-client ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/lumel .

COPY --from=builder /app/.env.example ./.env

RUN mkdir -p /data
VOLUME /data

EXPOSE 8080

ENV DB_HOST=postgres \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=sales_data \
    SERVER_PORT=8080 \
    CSV_PATH=/data/sales.csv

CMD ["./lumel"] 