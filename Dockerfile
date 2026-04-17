# Build Stage
FROM golang:1.25-alpine AS builder

## Waktu standar & sertifikat gembok SSL (Root CA) ke mesin Alpine
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
## RUN go build -o main cmd/app.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/app.go

# Run Stage
FROM alpine:latest

## Sertifikat gembok CA dan zona waktu hasil tarikan Builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfod

WORKDIR /app
COPY --from=builder /app/main .
#COPY --from=builder /app/cmd/tribun_mapping_real.csv ./cmd/

COPY --from=builder /app/migrations ./migrations

EXPOSE 8081
CMD ["./main"]
