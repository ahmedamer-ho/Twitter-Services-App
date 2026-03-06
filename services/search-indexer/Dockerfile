# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# We need to copy the shared package to satisfy local module replacement
# The build context for the entire docker-compose must be the project root
COPY go.work go.work.sum ./
COPY shared ./shared
COPY services/search-indexer ./services/search-indexer

WORKDIR /app/services/search-indexer

RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o searchindexer ./cmd/server

# Final stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/services/search-indexer/searchindexer .

EXPOSE 8084
CMD ["./searchindexer"]
