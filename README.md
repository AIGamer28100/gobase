# GoBase - Django-Inspired ORM for Go

GoBase is a powerful, Django-inspired ORM library for Go that provides elegant database operations with automatic migrations, user management, and CLI tools.

[![go report card](https://goreportcard.com/badge/github.com/AIGamer28100/gobase "go report card")](https://goreportcard.com/report/github.com/AIGamer28100/gobase)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Tests](https://github.com/AIGamer28100/gobase/workflows/Tests/badge.svg?branch=main "test status")](https://github.com/AIGamer28100/gobase/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/AIGamer28100/gobase.svg)](https://pkg.go.dev/github.com/AIGamer28100/gobase)

## Features

- **Django-style CRUD Operations**: Create, Get, All, Filter, Update, Delete
- **Automatic Database Migrations**: Support for SQLite, PostgreSQL, and MongoDB
- **Built-in User Management**: Ready-to-use User model with authentication
- **JSON Data Preloading**: Load initial data from JSON files with duplicate prevention
- **CLI Tools**: Database migration, superuser creation, and data preloading commands
- **Soft Deletes**: Non-destructive record deletion
- **Model Validation**: Ensure proper model structure
- **Password Encryption**: Built-in bcrypt password hashing

## Installation

```bash
go get github.com/AIGamer28100/gobase
```

## Quick Start

### 1. Basic Setup

```go
package main

import (
    "log"
    "github.com/AIGamer28100/gobase"
)

// Define your model
type Article struct {
    gobase.BaseModel
    Title   string `json:"title"`
    Content string `json:"content"`
    Author  string `json:"author"`
    Status  string `json:"status" gorm:"default:draft"`
}

func main() {
    // Initialize database connection
    connection, err := gobase.InitDB()
    if err != nil {
        log.Fatal(err)
    }
    defer connection.Close()

    // Create accessor for database operations
    accessor := gobase.NewAccessor(connection)

    // Migrate database schema
    err = accessor.Migrate(&Article{})
    if err != nil {
        log.Fatal(err)
    }

    // Create a new article
    article := &Article{
        Title:   "Getting Started with GoBase",
        Content: "GoBase makes database operations simple...",
        Author:  "Developer",
        Status:  "published",
    }
    
    err = accessor.Create(article)
    if err != nil {
        log.Fatal(err)
    }
}
```

### 2. Environment Configuration

Create a `.env` file in your project root:

```env
DB_TYPE=sqlite
DB_NAME=myapp.db

# For PostgreSQL:
# DB_TYPE=postgres
# DB_HOST=localhost
# DB_PORT=5432
# DB_USER=username
# DB_PASSWORD=password
# DB_NAME=mydb
```

### 3. CRUD Operations

```go
// Create
article := &Article{Title: "New Article", Content: "Content here"}
err := accessor.Create(article)

// Get by ID
var article Article
err := accessor.Get(&article, 1)

// Get all records
var articles []Article
err := accessor.All(&articles)

// Filter records
var publishedArticles []Article
err := accessor.Filter(&publishedArticles, "status = ?", "published")

// Update
article.Status = "archived"
err := accessor.Update(&article)

// Delete (soft delete)
err := accessor.Delete(&article)
```

### 4. User Management

```go
// Create a superuser
user := &gobase.User{
    Username: "admin",
    Email:    "admin@example.com",
}
err := user.SetPassword("SecurePassword123!")
if err != nil {
    log.Fatal(err)
}

err = gobase.CreateSuperuser(accessor, user)
if err != nil {
    log.Fatal(err)
}

// Verify password
isValid := user.CheckPassword("SecurePassword123!")
```

## CLI Tools

GoBase provides a powerful CLI tool for common database operations.

### Installation

```bash
go install github.com/AIGamer28100/gobase/cmd/gobase@latest
```

### Usage

```bash
# Run database migration
gobase -migrate

# Create a superuser
gobase -createsuperuser -username admin -email admin@example.com

# Preload data from JSON files
gobase -preload -files articles.json,users.json

# Get help
gobase -help
```

### JSON Preloading

Create JSON files with data to preload. The filename (without extension) should match your model name:

**articles.json**:
```json
[
  {
    "id": 1,
    "title": "Sample Article",
    "content": "This is sample content...",
    "author": "Author Name",
    "status": "published"
  }
]
```

## Advanced Features

### Custom Models

```go
type CustomModel struct {
    gobase.BaseModel
    ID   string `gorm:"primarykey;size:36" json:"id"` // Override default ID
    Name string `json:"name"`
}
```

### Model Registry for Preloading

```go
modelRegistry := map[string]interface{}{
    "articles": &Article{},
    "users":    &gobase.User{},
}

err := accessor.Preload(modelRegistry, "data.json")
```

### Database Configuration

```go
config := gobase.DatabaseConfig{
    Type: "postgres",
    Host: "localhost",
    Port: 5432,
    User: "username",
    Password: "password",
    Database: "mydb",
}

connection, err := gobase.InitDBWithConfig(config)
```

## Testing

Run the test suite:

```bash
cd gobase
go test -v
```

## Database Support

- **SQLite**: Default, perfect for development and small applications
- **PostgreSQL**: Production-ready relational database
- **MongoDB**: Document database support (coming soon)

## Development

### Setup Development Environment

1. Clone the repository:
```bash
git clone https://github.com/AIGamer28100/gobase.git
cd gobase
```

2. Install development dependencies:
```bash
go mod download
```

3. Install Git hooks for code quality (recommended):
```bash
make install-hooks
```

The Git hooks will automatically run before each commit to ensure:
- ✅ Go modules are in sync
- ✅ Code passes `go vet`
- ✅ All tests pass
- ✅ Code passes linting (golangci-lint)
- ✅ Code passes security scan (gosec)
- ✅ No binary files are committed
- ✅ Code is properly formatted

### Development Commands

```bash
# Build the project
make build

# Run tests
make test

# Run linter
make lint

# Run security scan
make security

# Run all checks (like CI)
make test lint security

# Install Git hooks
make install-hooks

# Clean build artifacts
make clean
```

### Manual Hook Installation

If you prefer to install hooks manually:
```bash
./scripts/install-hooks.sh
```

To bypass hooks for a commit (not recommended):
```bash
git commit --no-verify
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Add tests for your changes
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Examples

Check the [example application](../gobase-example/) for a complete implementation showing all features.
