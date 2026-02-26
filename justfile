# justfile

# Default command: show available recipes
default:
    @just --list

# Run unit tests using the golden files in testdata
test:
    go test -v ./...

# Build the binary
build:
    go build -o k8s-sync .

# Refresh an environment using local testdata (dry-run style)
refresh-local env="dev":
    go run . --env {{env}} --local ./testdata

# Refresh an environment using the real K8s cluster
refresh env:
    go run . --env {{env}}

# Clean up build artifacts
clean:
    rm -f k8s-sync

# Tidy up go modules
tidy:
    go mod tidy

# Format the code
fmt:
    go fmt ./...
