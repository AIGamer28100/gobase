# Version Management Guide

GoBase uses a Poetry-like version management system that makes it easy to bump versions and manage releases.

## Current Version System

The version is maintained in `cmd/gobase/version.go` and can be easily updated using the provided tools.

## Available Commands

### Show Current Version
```bash
make version                    # Show current version with build info
./build/gobase -version        # Show version from built binary
```

### Bump Version (Poetry-style)
```bash
# Patch version (0.0.1 -> 0.0.2)
make version-patch

# Minor version (0.1.0 -> 0.2.0)  
make version-minor

# Major version (1.0.0 -> 2.0.0)
make version-major

# Prerelease versions
make version-alpha              # 1.0.0 -> 1.0.0-alpha
make version-beta               # 1.0.0 -> 1.0.0-beta  
make version-rc                 # 1.0.0 -> 1.0.0-rc
```

### Manual Version Script
You can also use the version script directly:
```bash
# Show help
./scripts/version.sh

# Bump versions
./scripts/version.sh patch
./scripts/version.sh minor
./scripts/version.sh major
./scripts/version.sh prerelease alpha
./scripts/version.sh prerelease beta
./scripts/version.sh prerelease rc
```

## Release Process

### Quick Release
```bash
make release                   # Full automated release process
```

This will:
1. Clean build artifacts
2. Run tests
3. Run linter
4. Build for all platforms
5. Create git tag
6. Show next steps

### Manual Release Steps
```bash
# 1. Bump version
make version-patch             # or minor/major

# 2. Build and test
make clean test build-all

# 3. Commit changes
git add .
git commit -m "Bump version to v$(make version | head -1 | cut -d' ' -f3)"

# 4. Create and push tag
make tag
git push origin main --tags
```

## GitHub Actions Integration

The workflows automatically:
- Use the version from `cmd/gobase/version.go`
- Inject build date and git commit into binaries
- Create releases when tags are pushed
- Run tests on version bumps

## Version File Structure

The version is maintained in `cmd/gobase/version.go`:

```go
var (
    Version     = "v0.0.1-alpha"
    Name        = "GoBase CLI"
    Description = "A Django-inspired ORM and database toolkit for Go"
    BuildDate   = "unknown"  // Injected at build time
    GitCommit   = "unknown"  // Injected at build time
)
```

## Examples

### Development Workflow
```bash
# Start with alpha version
make version-alpha              # v0.1.0-alpha

# Make changes, test
make dev-run                   # Build and test quickly

# Ready for beta
make version-beta              # v0.1.0-beta

# More testing...

# Ready for release candidate  
make version-rc                # v0.1.0-rc

# Final release
make version-patch             # v0.1.0 (removes prerelease)
make release                   # Full release
```

### Major Release
```bash
# Current: v0.5.2
make version-major             # v1.0.0
make release
```

This system provides the same convenience as Poetry's version management while being native to Go and integrating seamlessly with the build system.
