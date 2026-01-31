FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY services/users-service/go.mod services/users-service/
RUN go mod download
COPY . .
WORKDIR /app/services/users-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /main .
CMD ["./main"]
