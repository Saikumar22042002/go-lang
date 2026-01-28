# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go-lang-app .

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

# Copy the built binary from the builder stage
COPY --from=builder /go-lang-app .

# Ensure the binary is owned by the non-root user
RUN chown appuser:appgroup /home/appuser/go-lang-app

# Switch to the non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Set the entrypoint for the container
CMD ["./go-lang-app"]
