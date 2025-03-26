# Use the official Go image as the base
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
# Copy the templates directory
COPY --from=builder /app/templates ./templates
# Copy the assets directory
COPY --from=builder /app/assets ./assets

# Expose the port
EXPOSE 4747

# Run the application
CMD ["./main"] 