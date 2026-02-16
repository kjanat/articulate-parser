# Build stage
FROM golang:1.26-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies and HTTPS)
RUN apk add --no-cache git ca-certificates tzdata file

# Create a non-root user for the final stage
RUN adduser -D -u 1000 appuser

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
# Docker buildx automatically provides these for multi-platform builds
ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

# Debug: Show build information
RUN echo "Building for platform: $TARGETPLATFORM (OS: $TARGETOS, Arch: $TARGETARCH, Variant: $TARGETVARIANT)" \
&& echo "Build platform: $BUILDPLATFORM" \
&& echo "Go version: $(go version)"

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
-ldflags="-s -w -X github.com/kjanat/articulate-parser/internal/version.Version=${VERSION} -X github.com/kjanat/articulate-parser/internal/version.BuildTime=${BUILD_TIME} -X github.com/kjanat/articulate-parser/internal/version.GitCommit=${GIT_COMMIT}" \
-o articulate-parser \
./main.go

# Verify the binary architecture
RUN file /app/articulate-parser || echo "file command not available"

# Final stage - minimal runtime image
FROM scratch

# Copy CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Add a minimal /etc/passwd file to support non-root user
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary
COPY --from=builder /app/articulate-parser /articulate-parser

# Switch to non-root user (appuser with UID 1000)
USER appuser

# Set the binary as entrypoint
ENTRYPOINT ["/articulate-parser"]

# Default command shows help
CMD ["--help"]

# Add labels for metadata
LABEL org.opencontainers.image.title="Articulate Parser"
LABEL org.opencontainers.image.description="A powerful CLI tool to parse Articulate Rise courses and export them to multiple formats (Markdown, HTML, DOCX). Supports media extraction, content cleaning, and batch processing for educational content conversion."
LABEL org.opencontainers.image.vendor="kjanat"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/kjanat/articulate-parser"
LABEL org.opencontainers.image.documentation="https://github.com/kjanat/articulate-parser/blob/master/DOCKER.md"
