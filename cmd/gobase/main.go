package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/AIGamer28100/gobase"
)

func main() {
	// Define command-line flags
	var (
		migrateCmd         = flag.Bool("migrate", false, "Run database migration")
		createSuperuserCmd = flag.Bool("createsuperuser", false, "Create a superuser")
		preloadCmd         = flag.Bool("preload", false, "Preload data from JSON files")
		username           = flag.String("username", "", "Username for the superuser")
		email              = flag.String("email", "", "Email for the superuser")
		password           = flag.String("password", "", "Password for the superuser (if not provided, will be prompted)")
		jsonFiles          = flag.String("files", "", "Comma-separated list of JSON files for preloading")
		version            = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *version {
		fmt.Println("GoBase CLI v1.0.0")
		fmt.Println("A Django-inspired ORM and database toolkit for Go")
		return
	}

	fmt.Println("=== GoBase CLI ===")

	// Initialize database connection using environment variables
	connection, err := gobase.InitDB()
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		fmt.Println("Note: Make sure to create a .env file with database configuration")
		fmt.Println("Example .env file:")
		fmt.Println("DB_TYPE=sqlite")
		fmt.Println("DB_NAME=test.db")
		os.Exit(1)
	}
	defer connection.Close()

	fmt.Printf("✓ Connected to %s database\n", connection.Type)

	// Create accessor
	accessor := gobase.NewAccessor(connection)

	// Handle commands
	if *migrateCmd {
		handleMigrate(accessor)
		return
	}

	if *createSuperuserCmd {
		handleCreateSuperuser(accessor, *username, *email, *password)
		return
	}

	if *preloadCmd {
		handlePreload(accessor, *jsonFiles)
		return
	}

	// If no specific command, show usage
	showUsage()
}

func showUsage() {
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  -migrate              Run database migration")
	fmt.Println("  -createsuperuser      Create a superuser")
	fmt.Println("  -preload              Preload data from JSON files")
	fmt.Println("  -version              Show version information")
	fmt.Println()
	fmt.Println("Use -help for detailed flag information")
}

func handleMigrate(accessor *gobase.Accessor) {
	fmt.Println()
	fmt.Println("=== Running Database Migration ===")
	fmt.Println()

	// Migrate with automatic model registration (includes User model)
	err := accessor.Migrate()
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("✓ Database migration completed successfully")
}

func handleCreateSuperuser(accessor *gobase.Accessor, username, email, password string) {
	fmt.Println()
	fmt.Println("=== Creating Superuser ===")
	fmt.Println()

	// Get username if not provided
	if username == "" {
		fmt.Print("Username: ")
		fmt.Scanln(&username)
	}

	// Get email if not provided
	if email == "" {
		fmt.Print("Email: ")
		fmt.Scanln(&email)
	}

	// Get password if not provided
	if password == "" {
		fmt.Print("Password: ")
		fmt.Scanln(&password)
	}

	// Create superuser
	err := gobase.CreateSuperuser(accessor, username, email, password)
	if err != nil {
		log.Fatalf("Failed to create superuser: %v", err)
	}

	fmt.Printf("✓ Superuser '%s' created successfully!\n", username)
}

func handlePreload(accessor *gobase.Accessor, files string) {
	fmt.Println()
	fmt.Println("=== Preloading Data ===")
	fmt.Println()

	if files == "" {
		log.Fatal("No files specified. Use -files flag to specify JSON files")
	}

	// For simplicity, we'll require the user to define their own model registry
	// in their application. This CLI version will show an example.
	fmt.Println("Note: This CLI tool requires you to define model registries in your application.")
	fmt.Println("Please use the gobase package directly in your Go application for preloading.")
	fmt.Println()
	fmt.Println("Example usage in your Go code:")
	fmt.Println(`
modelRegistry := map[string]interface{}{
    "users": &YourUserModel{},
    "articles": &YourArticleModel{},
}

err := accessor.Preload(modelRegistry, "data.json")
`)
}
