#!/usr/bin/env bash
set -euo pipefail

# Get the *real* path to the script, even if called via symlink
SCRIPT_PATH="$(readlink -f "$0" 2>/dev/null || realpath "$0")"
SCRIPT_DIR="$(dirname "$SCRIPT_PATH")"
PARENT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PARENT_DIR"

# Default values
OS=("darwin" "freebsd" "linux" "windows")
ARCH=("amd64" "arm64")
OUTDIR="build"
ENTRYPOINT="main.go"
JOBS=4
SHOW_TARGETS=false
SHOW_HELP=false
VERBOSE=false
DEFAULT_LDFLAGS="-s -w"

# Function to show help
show_help() {
    cat <<'EOF'
articulate-parser Build Script (Bash)
=====================================

SYNOPSIS:
    build.sh [OPTIONS] [GO_BUILD_FLAGS...]

DESCRIPTION:
    Cross-platform build script for articulate-parser. Builds binaries for multiple
    OS/architecture combinations in parallel with embedded version information.

OPTIONS:
    -h, --help          Show this help message and exit
    -j <number>         Number of parallel jobs (default: 4)
    -o <directory>      Output directory for binaries (default: build)
    -e <file>           Entry point Go file (default: main.go)
    -v, --verbose       Enable verbose output for debugging
    --show-targets      Show available Go build targets and exit

EXAMPLES:
    # Basic build with default settings
    ./scripts/build.sh

    # Build with 8 parallel jobs
    ./scripts/build.sh -j 8

    # Build to custom directory
    ./scripts/build.sh -o my_builds

    # Build with custom entry point
    ./scripts/build.sh -e test_entry.go

    # Build with verbose output
    ./scripts/build.sh -v

    # Build with Go build flags and version info
    ./scripts/build.sh -ldflags "-s -w -X main.version=1.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"

    # Show available targets
    ./scripts/build.sh --show-targets

    # Build with Go build flags and version info
    ./scripts/build.sh -ldflags "-s -w -X github.com/kjanat/articulate-parser/internal/version.Version=1.0.0 -X github.com/kjanat/articulate-parser/internal/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"

    # Build with custom ldflags (overrides default -s -w)
    ./scripts/build.sh -ldflags "-X github.com/kjanat/articulate-parser/internal/version.Version=1.0.0"

    # Build without any ldflags (disable defaults)
    ./scripts/build.sh -ldflags ""

DEFAULT TARGETS:
    Operating Systems: darwin, freebsd, linux, windows
    Architectures:     amd64, arm64
    
    This creates 8 binaries total (4 OS × 2 ARCH)

GO BUILD FLAGS:
    Any additional arguments are passed directly to 'go build'.
    Default: -ldflags "-s -w" (strip debug info for smaller binaries)
    Common flags include:
        -ldflags    Link flags (e.g., version info, optimization)
        -tags       Build tags
        -v          Verbose Go build output
        -race       Enable race detector
        -trimpath   Remove file system paths from executable

    To override default ldflags, specify your own -ldflags argument.
    To disable ldflags entirely, use: -ldflags ""

OUTPUT:
    Binaries are named: articulate-parser-{OS}-{ARCH}[.exe]
    Build logs for failed builds: {BINARY_NAME}.log

NOTES:
    - Requires Go to be installed and in PATH
    - Removes and recreates the output directory
    - Failed builds create .log files with error details
    - Uses colored output with real-time status updates
    - Entry point validation ensures file exists before building
    - Supports custom entry points (compiles single file to avoid conflicts)

EOF
}

# Function to show available Go build targets
show_targets() {
    echo "Available Go Build Targets:"
    echo "=========================="
    echo
    if command -v go >/dev/null 2>&1; then
        echo "Getting targets from 'go tool dist list'..."
        echo

        # Get all targets and format them nicely
        local targets
        targets=$(go tool dist list 2>/dev/null)

        if [ $? -eq 0 ] && [ -n "$targets" ]; then
            # Show formatted output
            printf "%-15s %-10s %s\n" "OS" "ARCH" "STATUS"
            printf "%-15s %-10s %s\n" "---------------" "----------" "------"

            # Track our default targets
            local default_targets=()
            for os in "${OS[@]}"; do
                for arch in "${ARCH[@]}"; do
                    default_targets+=("$os/$arch")
                done
            done

            # Display all targets with status
            echo "$targets" | sort | while IFS='/' read -r os arch; do
                local status="available"
                if printf '%s\n' "${default_targets[@]}" | grep -q "^$os/$arch$"; then
                    status="default"
                fi
                printf "%-15s %-10s %s\n" "$os" "$arch" "$status"
            done

            echo
            echo "Summary:"
            echo "  Total available targets: $(echo "$targets" | wc -l)"
            echo "  Default targets used by this script: ${#default_targets[@]}"
            echo
            echo "Default script targets:"
            for target in "${default_targets[@]}"; do
                echo "  - $target"
            done
        else
            echo "Error: Failed to get target list from 'go tool dist list'"
            exit 1
        fi
    else
        echo "Error: Go is not installed or not in PATH"
        exit 1
    fi
}

# Parse parameters
while (("$#")); do
    case $1 in
    -h | --help)
        SHOW_HELP=true
        shift
        ;;
    --show-targets)
        SHOW_TARGETS=true
        shift
        ;;
    -v | --verbose)
        VERBOSE=true
        shift
        ;;
    -j)
        if [[ ${2-} =~ ^[0-9]+$ ]]; then
            JOBS=$2
            shift 2
        else
            echo "Error: Missing number of jobs after -j"
            exit 1
        fi
        ;;
    -o)
        if [ -n "${2-}" ]; then
            OUTDIR=$2
            shift 2
        else
            echo "Error: Missing output directory after -o"
            exit 1
        fi
        ;;
    -e)
        if [ -n "${2-}" ]; then
            ENTRYPOINT=$2
            shift 2
        else
            echo "Error: Missing entry point file after -e"
            exit 1
        fi
        ;;
    *)
        break
        ;;
    esac
done

# Handle help and show-targets early
if [ "$SHOW_HELP" = true ]; then
    show_help
    exit 0
fi

if [ "$SHOW_TARGETS" = true ]; then
    show_targets
    exit 0
fi

# Validate Go installation
if ! command -v go >/dev/null 2>&1; then
    echo "Error: Go is not installed or not in PATH"
    echo "Please install Go from https://golang.org/dl/"
    echo "Or if running on Windows, use the PowerShell script: scripts\\build.ps1"
    exit 1
fi

# Validate entry point exists
if [ ! -f "$ENTRYPOINT" ]; then
    echo "Error: Entry point file '$ENTRYPOINT' does not exist"
    exit 1
fi

# Store remaining arguments as an array to preserve argument boundaries
GO_BUILD_FLAGS_ARRAY=("$@")

# Apply default ldflags if no custom ldflags were provided
HAS_CUSTOM_LDFLAGS=false
for arg in "${GO_BUILD_FLAGS_ARRAY[@]}"; do
    if [[ "$arg" == "-ldflags" ]]; then
        HAS_CUSTOM_LDFLAGS=true
        break
    fi
done

if [[ "$HAS_CUSTOM_LDFLAGS" == false ]] && [[ -n "$DEFAULT_LDFLAGS" ]]; then
    # Add default ldflags at the beginning
    GO_BUILD_FLAGS_ARRAY=("-ldflags" "$DEFAULT_LDFLAGS" "${GO_BUILD_FLAGS_ARRAY[@]}")
fi

# Verbose output
if [ "$VERBOSE" = true ]; then
    echo "Build Configuration:"
    echo "  Entry Point: $ENTRYPOINT"
    echo "  Output Dir:  $OUTDIR"
    echo "  Parallel Jobs: $JOBS"
    if [ ${#GO_BUILD_FLAGS_ARRAY[@]} -gt 0 ]; then
        echo "  Go Build Flags: ${GO_BUILD_FLAGS_ARRAY[*]}"
    else
        echo "  Go Build Flags: none"
    fi
    echo "  Targets: ${#OS[@]}×${#ARCH[@]} = $((${#OS[@]} * ${#ARCH[@]})) total"
    echo
fi

rm -rf "$OUTDIR"
mkdir -p "$OUTDIR"

# Get build start time
BUILD_START=$(date +%s)

# Compose all targets in an array
TARGETS=()
for os in "${OS[@]}"; do
    for arch in "${ARCH[@]}"; do
        BIN="articulate-parser-$os-$arch"
        [[ "$os" == "windows" ]] && BIN="$BIN.exe"
        TARGETS+=("$BIN|$os|$arch")
    done
done

# Show targets info if verbose
if [ "$VERBOSE" = true ]; then
    echo "Building targets:"
    for target in "${TARGETS[@]}"; do
        BIN="${target%%|*}"
        echo "  - $BIN"
    done
    echo
fi

# Print pending statuses and save line numbers
for idx in "${!TARGETS[@]}"; do
    BIN="${TARGETS[$idx]%%|*}"
    printf "[ ] %-35s ... pending\n" "$BIN"
done

# Make sure output isn't buffered
export PYTHONUNBUFFERED=1

# Function to update a line in-place (1-based index)
update_status() {
    local idx=$1
    local symbol=$2
    local msg=$3
    # Move cursor up to the correct line
    printf "\0337"                                  # Save cursor position
    printf "\033[%dA" $((${#TARGETS[@]} - idx + 1)) # Move up
    printf "\r\033[K[%s] %-35s\n" "$symbol" "$msg"  # Clear & update line
    printf "\0338"                                  # Restore cursor position
}

for idx in "${!TARGETS[@]}"; do
    while (($(jobs -rp | wc -l) >= JOBS)); do sleep 0.2; done
    (
        IFS='|' read -r BIN os arch <<<"${TARGETS[$idx]}"
        update_status $((idx + 1)) '>' "$BIN ... building"

        # Prepare build command as an array to properly handle arguments with spaces
        build_cmd=(go build)
        if [ "$VERBOSE" = true ]; then
            build_cmd+=(-v)
        fi
        build_cmd+=("${GO_BUILD_FLAGS_ARRAY[@]}" -o "$OUTDIR/$BIN" "$ENTRYPOINT")

        if CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" "${build_cmd[@]}" 2>"$OUTDIR/$BIN.log"; then
            update_status $((idx + 1)) '✔' "$BIN done"
            rm -f "$OUTDIR/$BIN.log"
        else
            update_status $((idx + 1)) '✖' "$BIN FAILED (see $OUTDIR/$BIN.log)"
        fi
    ) &
done

wait

# Calculate build time
BUILD_END=$(date +%s)
BUILD_DURATION=$((BUILD_END - BUILD_START))

echo -e "\nAll builds completed in ${BUILD_DURATION}s. Find them in $OUTDIR/"

# Show build summary if verbose
if [ "$VERBOSE" = true ]; then
    echo
    echo "Build Summary:"
    echo "=============="
    success_count=0
    total_size=0

    for target in "${TARGETS[@]}"; do
        BIN="${target%%|*}"
        if [ -f "$OUTDIR/$BIN" ]; then
            success_count=$((success_count + 1))
            size=$(stat -f%z "$OUTDIR/$BIN" 2>/dev/null || stat -c%s "$OUTDIR/$BIN" 2>/dev/null || echo "0")
            total_size=$((total_size + size))
            rm -f "$OUTDIR/$BIN.log"
            printf "  ✔ %-42s %s\n" "$OUTDIR/$BIN" "$(numfmt --to=iec-i --suffix=B $size 2>/dev/null || echo "${size} bytes")"
        else
            printf "  ✖ %-42s %s\n" "$OUTDIR/$BIN" "FAILED"
        fi
    done

    echo "  ────────────────────────────────────────────────"
    printf "  Total: %d/%d successful, %s total size\n" "$success_count" "${#TARGETS[@]}" "$(numfmt --to=iec-i --suffix=B $total_size 2>/dev/null || echo "${total_size} bytes")"
fi

# Clean up environment variables to avoid contaminating future builds
unset GOOS GOARCH CGO_ENABLED
