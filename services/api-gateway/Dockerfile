FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY services/api-gateway/go.mod services/api-gateway/
RUN go mod download
COPY . .
WORKDIR /app/services/api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /main .
CMD ["./main"]
