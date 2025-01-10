# Use the official Golang image to build the application
FROM golang:1.20 as builder

# Set the working directory
WORKDIR /app

# Copy the source code
COPY main.go .

# Build the Go application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -o hello-twitter -a -installsuffix cgo main.go

# Use a smaller base image for the runtime
FROM gcr.io/distroless/base-debian11

# Copy the compiled binary from the builder stage
COPY --from=builder /app/hello-twitter /

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["/hello-twitter"]