# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies and HTTPS)
RUN apk add --no-cache git ca-certificates tzdata

# Set the working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# Disable CGO for a fully static binary
# Use linker flags to reduce binary size and embed version info
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X github.com/kjanat/articulate-parser/internal/version.Version=${VERSION} -X github.com/kjanat/articulate-parser/internal/version.BuildTime=${BUILD_TIME} -X github.com/kjanat/articulate-parser/internal/version.GitCommit=${GIT_COMMIT}" \
    -o articulate-parser \
    ./main.go

# Final stage - minimal runtime image
FROM scratch

# Copy CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/articulate-parser /articulate-parser

# Create a non-root user (note: in scratch image, we can't create users dynamically)
# The binary will run as root in scratch, which is acceptable for containers

# Set the binary as entrypoint
ENTRYPOINT ["/articulate-parser"]

# Default command shows help
CMD ["--help"]

# Add labels for metadata
LABEL org.opencontainers.image.title="Articulate Parser"
LABEL org.opencontainers.image.description="A tool to parse Articulate Rise courses and export them to various formats"
LABEL org.opencontainers.image.vendor="kjanat"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/kjanat/articulate-parser"
LABEL org.opencontainers.image.documentation="https://github.com/kjanat/articulate-parser/blob/master/README.md"
