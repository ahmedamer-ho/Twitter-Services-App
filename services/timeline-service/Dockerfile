FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/services/timeline-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /main .
CMD ["./main"]
