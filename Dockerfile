# Build stage
FROM golang:1.21-alpine AS builder

# Install required packages
RUN apk add --no-cache git curl

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install oapi-codegen
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Copy source code
COPY . .

# Generate API client and build
RUN /go/bin/oapi-codegen -config oapi-codegen.yaml spec/coolify-openapi.yaml
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o coolifyme cmd/*.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates curl

# Create non-root user
RUN adduser -D -s /bin/sh coolify

# Set working directory
WORKDIR /home/coolify

# Copy binary from builder
COPY --from=builder /app/coolifyme /usr/local/bin/coolifyme

# Switch to non-root user
USER coolify

# Create config directory
RUN mkdir -p /home/coolify/.config/coolifyme

ENTRYPOINT ["coolifyme"]
CMD ["--help"] 