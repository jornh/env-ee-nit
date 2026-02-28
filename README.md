

Companion for <https://github.com/dhth/envee/> to init (and update) the versions.toml file

Queries a K8s namespace or local yaml manifests for deployments and extracts image names and tags. The `versions.toml` file is updated accordingly

See the justfile or executable help:
```bash
just -l
Available recipes:
    build                   # Build the binary
    clean                   # Clean up build artifacts
    default                 # Default command: show available recipes
    fmt                     # Format the code
    refresh env             # Refresh an environment using the real K8s cluster
    refresh-local env="dev" # Refresh an environment using local testdata (dry-run style)
    test                    # Run unit tests using the golden files in testdata
    tidy                    # Tidy up go modules

# ------------
./env-ee-nit --help
Usage of ./env-ee-nit:
  -env string
        The environment/namespace to refresh
  -local string
        Path to local k8s yaml files (optional)
  -versions string
        Path to output envee versions file (default "versions.toml")
```

