package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/AIGamer28100/gobase"
	"golang.org/x/term"
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
		versionFlag        = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	if *versionFlag {
		fmt.Println(getVersionInfo())
		fmt.Println("Build Date:", BuildDate)
		fmt.Println("Git Commit:", GitCommit)
		return
	}

	// If no commands are provided, show help
	if !*migrateCmd && !*createSuperuserCmd && !*preloadCmd {
		printHelp()
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
	}

	if *createSuperuserCmd {
		handleCreateSuperuser(accessor, *username, *email, *password)
	}

	if *preloadCmd {
		handlePreload(accessor, *jsonFiles)
	}
}

func printHelp() {
	fmt.Println("GoBase CLI - Database management tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gobase [command] [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  -migrate              Run database migration")
	fmt.Println("  -createsuperuser      Create a superuser")
	fmt.Println("  -preload              Preload data from JSON files")
	fmt.Println("  -version              Show version information")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -username string      Username for the superuser")
	fmt.Println("  -email string         Email for the superuser")
	fmt.Println("  -password string      Password for the superuser")
	fmt.Println("  -files string         Comma-separated list of JSON files")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gobase -migrate")
	fmt.Println("  gobase -createsuperuser -username admin -email admin@example.com")
	fmt.Println("  gobase -preload -files articles.json,users.json")
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
		username = promptInput("Username: ")
	}

	// Get email if not provided
	if email == "" {
		email = promptInput("Email: ")
	}

	// Get password if not provided
	if password == "" {
		password = promptPassword("Password: ")
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

	fileList := strings.Split(files, ",")
	for i, file := range fileList {
		fileList[i] = strings.TrimSpace(file)
	}

	// Create a basic model registry for common models
	modelRegistry := map[string]interface{}{
		"users": &gobase.User{},
		"articles": &struct {
			gobase.BaseModel
			Title   string `json:"title"`
			Content string `json:"content"`
			Author  string `json:"author"`
			Status  string `json:"status"`
		}{},
	}

	err := accessor.Preload(modelRegistry, fileList...)
	if err != nil {
		log.Fatalf("Preload failed: %v", err)
	}

	fmt.Printf("✓ Successfully preloaded data from %d files\n", len(fileList))
}

func promptInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptPassword(prompt string) string {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	fmt.Println()
	return string(bytePassword)
}
