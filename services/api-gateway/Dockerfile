FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY services/api-gateway/go.mod services/api-gateway/
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main services/api-gateway/cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /main .
CMD ["./main"]
