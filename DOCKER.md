# Articulate Parser - Docker

A powerful command-line tool for parsing and processing articulate data files, now available as a lightweight Docker container.

## Quick Start

### Pull from GitHub Container Registry

```bash
docker pull ghcr.io/kjanat/articulate-parser:latest
```

### Run with local files

```bash
docker run --rm -v $(pwd):/data ghcr.io/kjanat/articulate-parser:latest /data/input.txt
```

## Usage

### Basic File Processing

```bash
# Process a single file
docker run --rm -v $(pwd):/data ghcr.io/kjanat/articulate-parser:latest /data/document.txt

# Process with output redirection
docker run --rm -v $(pwd):/data ghcr.io/kjanat/articulate-parser:latest /data/input.txt > output.json
```

### Display Help and Version

```bash
# Show help information
docker run --rm ghcr.io/kjanat/articulate-parser:latest --help

# Show version
docker run --rm ghcr.io/kjanat/articulate-parser:latest --version
```

## Available Tags

-   `latest` - Latest stable release
-   `v1.x.x` - Specific version tags
-   `main` - Latest development build

## Image Details

-   **Base Image**: `scratch` (minimal attack surface)
-   **Architecture**: Multi-arch support (amd64, arm64)
-   **Size**: < 10MB (optimized binary)
-   **Security**: Runs as non-root user
-   **Features**: SBOM and provenance attestation included

## Development

### Local Build

```bash
docker build -t articulate-parser .
```

### Docker Compose

```bash
docker-compose up --build
```

## Repository

-   **Source**: [github.com/kjanat/articulate-parser](https://github.com/kjanat/articulate-parser)
-   **Issues**: [Report bugs or request features](https://github.com/kjanat/articulate-parser/issues)
-   **License**: See repository for license details
