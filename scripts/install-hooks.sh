#!/bin/bash

# Git Hooks Installation Script for GoBase
# This script installs pre-commit hooks to ensure code quality

set -e

echo "🔧 Installing Git hooks for GoBase..."

# Get the repository root
REPO_ROOT=$(git rev-parse --show-toplevel)
HOOKS_DIR="$REPO_ROOT/.git/hooks"
SCRIPTS_DIR="$REPO_ROOT/scripts"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Create the pre-commit hook
cat > "$HOOKS_DIR/pre-commit" << 'EOF'
#!/bin/bash

# Pre-commit hook for GoBase
# This hook runs linting and security checks before allowing commits

set -e

echo "🔍 Running pre-commit checks..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if we have any Go files in the staged changes
staged_go_files=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' | wc -l)

if [ "$staged_go_files" -eq 0 ]; then
    print_warning "No Go files staged for commit, skipping Go-specific checks"
    exit 0
fi

print_status "Found $staged_go_files Go file(s) staged for commit"

# 1. Check if go.mod and go.sum are in sync
echo "🔧 Checking Go modules..."
if ! go mod verify > /dev/null 2>&1; then
    print_error "go mod verify failed"
    echo "Run 'go mod tidy' to fix module dependencies"
    exit 1
fi
print_status "Go modules verified"

# 2. Run go vet
echo "🔍 Running go vet..."
if ! go vet ./... > /dev/null 2>&1; then
    print_error "go vet failed"
    echo "Fix the issues reported by 'go vet ./...' before committing"
    exit 1
fi
print_status "go vet passed"

# 3. Run tests
echo "🧪 Running tests..."
if ! go test ./... > /dev/null 2>&1; then
    print_error "Tests failed"
    echo "Fix failing tests before committing"
    exit 1
fi
print_status "All tests passed"

# 4. Check if golangci-lint is available
if ! command -v golangci-lint &> /dev/null; then
    print_warning "golangci-lint not found, attempting to install..."
    if ! go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; then
        print_error "Failed to install golangci-lint"
        echo "Please install golangci-lint manually:"
        echo "go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        exit 1
    fi
    print_status "golangci-lint installed"
fi

# 5. Run linting
echo "🔍 Running golangci-lint..."
if ! golangci-lint run --timeout=5m > /dev/null 2>&1; then
    print_error "Linting failed"
    echo "Fix linting issues before committing:"
    echo "golangci-lint run"
    exit 1
fi
print_status "Linting passed"

# 6. Check if gosec is available
if ! command -v gosec &> /dev/null; then
    print_warning "gosec not found, attempting to install..."
    if ! go install github.com/securego/gosec/v2/cmd/gosec@latest; then
        print_error "Failed to install gosec"
        echo "Please install gosec manually:"
        echo "go install github.com/securego/gosec/v2/cmd/gosec@latest"
        exit 1
    fi
    print_status "gosec installed"
fi

# 7. Run security scan
echo "🔒 Running security scan..."
if ! gosec ./... > /dev/null 2>&1; then
    print_error "Security scan failed"
    echo "Fix security issues before committing:"
    echo "gosec ./..."
    exit 1
fi
print_status "Security scan passed"

# 8. Check for binary files being committed
echo "🗂️  Checking for binary files..."
staged_binaries=$(git diff --cached --name-only | grep -E "(gobase|gobase-cli|\.exe|\.bin)$" | wc -l)
if [ "$staged_binaries" -gt 0 ]; then
    print_error "Binary files detected in staged changes"
    echo "The following binary files are staged:"
    git diff --cached --name-only | grep -E "(gobase|gobase-cli|\.exe|\.bin)$"
    echo "Remove them with: git reset HEAD <file>"
    exit 1
fi
print_status "No binary files in commit"

# 9. Format check (optional - can auto-fix)
echo "🎨 Checking code formatting..."
unformatted=$(gofmt -l $(git diff --cached --name-only --diff-filter=ACM | grep '\.go$') 2>/dev/null)
if [ -n "$unformatted" ]; then
    print_warning "The following files are not properly formatted:"
    echo "$unformatted"
    echo ""
    read -p "Auto-format these files? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "$unformatted" | xargs gofmt -w
        echo "$unformatted" | xargs git add
        print_status "Files formatted and re-staged"
    else
        print_error "Commit aborted due to formatting issues"
        echo "Run 'gofmt -w .' to fix formatting"
        exit 1
    fi
fi
print_status "Code formatting verified"

print_status "All pre-commit checks passed! 🎉"
echo ""
EOF

# Make the hook executable
chmod +x "$HOOKS_DIR/pre-commit"

print_status "Pre-commit hook installed successfully!"
print_info "The hook will now run automatically before each commit"
print_info "It checks: go modules, go vet, tests, linting, security, binaries, and formatting"

echo ""
echo "To bypass the hook (not recommended), use: git commit --no-verify"
echo "To test the hook manually, run: .git/hooks/pre-commit"
echo ""
print_status "Setup complete! Your commits will now be automatically validated."
