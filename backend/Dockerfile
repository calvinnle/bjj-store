# Build stage
FROM golang:1.21-alpine AS builder

# Install git and certificates
RUN apk --no-cache add git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies with more verbose output
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main .

# Final stage - use minimal image
FROM scratch

# Copy certificates and timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary and config
COPY --from=builder /app/main /main
COPY --from=builder /app/config.yaml /config.yaml

# Expose port
EXPOSE 8080

# Command to run
ENTRYPOINT ["/main"]