package main

import (
	"fmt"
	"log"
	"time"

	"gobase"
)

// User model that embeds BaseModel
// This demonstrates how to use the gobase package with standard fields
type User struct {
	gobase.BaseModel
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"size:100;uniqueIndex" json:"email"`
	Age      int    `json:"age"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// Product model that overrides the ID field to use string instead of uint
// This demonstrates the field overriding capability of BaseModel
type Product struct {
	ID        string    `gorm:"primarykey;size:36" json:"id"` // Override with string ID
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `gorm:"size:200;not null" json:"name"`
	Price     float64   `gorm:"precision:10;scale:2" json:"price"`
	Category  string    `gorm:"size:100" json:"category"`
}

// Custom ID methods for Product to work with the overridden ID field
func (p *Product) GetID() interface{} {
	return p.ID
}

func (p *Product) SetID(id interface{}) {
	if id, ok := id.(string); ok {
		p.ID = id
	}
}

func main() {
	// Database connection string for PostgreSQL
	// Replace with your actual database credentials
	dsn := "host=localhost user=postgres password=yourpassword dbname=testdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	// Initialize database connection
	fmt.Println("Initializing database connection...")
	db, err := gobase.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create accessor instance - this follows the Dependency Injection pattern
	accessor := gobase.NewAccessor(db)

	// Auto-migrate the schema
	fmt.Println("Running database migrations...")
	err = accessor.AutoMigrate(&User{}, &Product{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Demonstrate CRUD operations with User model
	demonstrateUserCRUD(accessor)

	// Demonstrate CRUD operations with Product model (overridden ID field)
	demonstrateProductCRUD(accessor)

	// Demonstrate advanced features
	demonstrateAdvancedFeatures(accessor)
}

func demonstrateUserCRUD(accessor *gobase.Accessor) {
	fmt.Println("\n=== User CRUD Operations ===")

	// Create operation
	user := &User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Age:      30,
		IsActive: true,
	}

	fmt.Println("Creating user...")
	err := accessor.Create(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return
	}
	fmt.Printf("Created user with ID: %d\n", user.ID)

	// Read operation
	fmt.Println("Reading user by ID...")
	readUser := &User{}
	err = accessor.ReadByID(readUser, user.ID)
	if err != nil {
		log.Printf("Error reading user: %v", err)
		return
	}
	fmt.Printf("Read user: %+v\n", readUser)

	// Update operation
	fmt.Println("Updating user...")
	readUser.Name = "John Smith"
	readUser.Age = 31
	err = accessor.Update(readUser)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return
	}
	fmt.Printf("Updated user: %+v\n", readUser)

	// Create another user for FindAll demonstration
	user2 := &User{
		Name:     "Jane Doe",
		Email:    "jane.doe@example.com",
		Age:      28,
		IsActive: true,
	}
	err = accessor.Create(user2)
	if err != nil {
		log.Printf("Error creating second user: %v", err)
		return
	}

	// FindAll operation
	fmt.Println("Finding all users...")
	var users []User
	err = accessor.FindAll(&users)
	if err != nil {
		log.Printf("Error finding all users: %v", err)
		return
	}
	fmt.Printf("Found %d users:\n", len(users))
	for i, u := range users {
		fmt.Printf("  %d. %+v\n", i+1, u)
	}

	// Delete operation (soft delete)
	fmt.Println("Deleting user...")
	err = accessor.Delete(readUser)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return
	}
	fmt.Println("User deleted successfully")

	// Verify soft delete - user should not appear in normal queries
	fmt.Println("Verifying soft delete...")
	var activeUsers []User
	err = accessor.FindAll(&activeUsers)
	if err != nil {
		log.Printf("Error finding active users: %v", err)
		return
	}
	fmt.Printf("Active users after deletion: %d\n", len(activeUsers))
}

func demonstrateProductCRUD(accessor *gobase.Accessor) {
	fmt.Println("\n=== Product CRUD Operations (Custom String ID) ===")

	// Create operation with custom string ID
	product := &Product{
		ID:       "prod-001",
		Name:     "Laptop",
		Price:    999.99,
		Category: "Electronics",
	}

	fmt.Println("Creating product...")
	err := accessor.Create(product)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		return
	}
	fmt.Printf("Created product with ID: %s\n", product.ID)

	// Read operation
	fmt.Println("Reading product by ID...")
	readProduct := &Product{}
	err = accessor.ReadByID(readProduct, "prod-001")
	if err != nil {
		log.Printf("Error reading product: %v", err)
		return
	}
	fmt.Printf("Read product: %+v\n", readProduct)

	// Update operation
	fmt.Println("Updating product...")
	readProduct.Price = 899.99
	readProduct.Category = "Computers"
	err = accessor.Update(readProduct)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		return
	}
	fmt.Printf("Updated product: %+v\n", readProduct)

	// Create another product
	product2 := &Product{
		ID:       "prod-002",
		Name:     "Mouse",
		Price:    29.99,
		Category: "Accessories",
	}
	err = accessor.Create(product2)
	if err != nil {
		log.Printf("Error creating second product: %v", err)
		return
	}

	// FindAll operation
	fmt.Println("Finding all products...")
	var products []Product
	err = accessor.FindAll(&products)
	if err != nil {
		log.Printf("Error finding all products: %v", err)
		return
	}
	fmt.Printf("Found %d products:\n", len(products))
	for i, p := range products {
		fmt.Printf("  %d. %+v\n", i+1, p)
	}
}

func demonstrateAdvancedFeatures(accessor *gobase.Accessor) {
	fmt.Println("\n=== Advanced Features ===")

	// Conditional queries
	fmt.Println("Finding users with age > 25...")
	var adultUsers []User
	err := accessor.FindWhere(&adultUsers, "age > ?", 25)
	if err != nil {
		log.Printf("Error finding adult users: %v", err)
		return
	}
	fmt.Printf("Found %d adult users\n", len(adultUsers))

	// Count operation
	fmt.Println("Counting active users...")
	count, err := accessor.Count(&User{}, "is_active = ?", true)
	if err != nil {
		log.Printf("Error counting active users: %v", err)
		return
	}
	fmt.Printf("Total active users: %d\n", count)

	// Transaction example
	fmt.Println("Demonstrating transaction...")
	err = accessor.Transaction(func(txAccessor *gobase.Accessor) error {
		// Create a new user within transaction
		user := &User{
			Name:     "Transaction User",
			Email:    "transaction@example.com",
			Age:      25,
			IsActive: true,
		}

		err := txAccessor.Create(user)
		if err != nil {
			return fmt.Errorf("failed to create user in transaction: %w", err)
		}

		// Create a related product within same transaction
		product := &Product{
			ID:       "tx-prod-001",
			Name:     "Transaction Product",
			Price:    199.99,
			Category: "Test",
		}

		err = txAccessor.Create(product)
		if err != nil {
			return fmt.Errorf("failed to create product in transaction: %w", err)
		}

		fmt.Println("Both user and product created successfully in transaction")
		return nil
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
	} else {
		fmt.Println("Transaction completed successfully")
	}
}
