FROM golang:1.23.4 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o rate-limiter ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/rate-limiter .
COPY --from=builder /app/cmd/.env .
EXPOSE 8080
CMD ["./rate-limiter"]