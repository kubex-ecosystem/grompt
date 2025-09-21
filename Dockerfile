# Multi-stage Dockerfile for Grompt Gateway
# Image: gemx-grompt:latest
# Stage 1: Build frontend
FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package*.json ./

# Install dependencies (including dev dependencies for build)
RUN npm ci --no-fund --no-audit --loglevel=error

# Copy frontend source
COPY frontend/ ./

# Build frontend
RUN npm run build:static

# Stage 2: Build Go backend
FROM golang:1.25-alpine AS backend-builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Download dependencies
RUN go mod tidy

# Build the gateway binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=dev" \
    -o dist/grompt_linux_amd64 \
    ./cmd

# Stage 3: Final runtime image
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1001 grompt && \
    adduser -D -u 1001 -G grompt grompt

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=backend-builder /app/dist/grompt_linux_amd64 .

# Copy frontend assets from builder
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Copy configuration files
COPY config/ ./config/

# Set ownership
RUN chown -R grompt:grompt /app

# Switch to non-root user
USER grompt

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/v1/health || exit 1

# Set entrypoint
ENTRYPOINT ["bash", "./dist/grompt_linux_amd64", "gateway", "serve", "--config", "./config/config.example.yml", "&", "bash"]
