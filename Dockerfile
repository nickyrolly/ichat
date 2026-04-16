# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/app.go

# Run Stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/cmd/tribun_mapping_real.csv ./cmd/
EXPOSE 8081
CMD ["./main"]
