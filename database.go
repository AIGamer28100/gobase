package gobase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig holds the configuration for database connections
type DatabaseConfig struct {
	Type     string
	Host     string
	User     string
	Password string
	Name     string
	Port     int
}

// Connection represents a database connection that can be either GORM or MongoDB
type Connection struct {
	Type        string
	GormDB      *gorm.DB
	MongoDB     *mongo.Database
	MongoClient *mongo.Client
}

// GetDB returns the underlying database connection
func (c *Connection) GetDB() interface{} {
	if c.Type == mongoDBType {
		return c.MongoDB
	}
	return c.GormDB
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.Type == mongoDBType && c.MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return c.MongoClient.Disconnect(ctx)
	}

	if c.GormDB != nil {
		sqlDB, err := c.GormDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return nil
}

// LoadConfig loads database configuration from environment variables
func LoadConfig() (*DatabaseConfig, error) {
	// Try to load .env file if it exists
	_ = godotenv.Load()

	config := &DatabaseConfig{
		Type:     os.Getenv("DB_TYPE"),
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	// Parse port
	portStr := os.Getenv("DB_PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid DB_PORT value: %s", portStr)
		}
		config.Port = port
	}

	// Validate required fields
	if config.Type == "" {
		return nil, errors.New("DB_TYPE is required")
	}

	if config.Name == "" {
		return nil, errors.New("DB_NAME is required")
	}

	// Set defaults based on database type
	switch config.Type {
	case "postgres":
		if config.Host == "" {
			config.Host = "localhost"
		}
		if config.Port == 0 {
			config.Port = 5432
		}
		if config.User == "" {
			return nil, errors.New("DB_USER is required for PostgreSQL")
		}
	case "sqlite":
		// For SQLite, DB_NAME is the file path
		// No other fields are required
	case mongoDBType:
		if config.Host == "" {
			config.Host = "localhost"
		}
		if config.Port == 0 {
			config.Port = 27017
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %s. Supported types: postgres, sqlite, mongodb", config.Type)
	}

	return config, nil
}

// InitDB establishes a connection to the database using configuration from environment variables.
// This function follows the Single Responsibility Principle by only
// handling database connection initialization.
func InitDB() (*Connection, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return InitDBWithConfig(config)
}

// InitDBWithConfig establishes a connection using the provided configuration
func InitDBWithConfig(config *DatabaseConfig) (*Connection, error) {
	switch config.Type {
	case "postgres":
		return initPostgreSQL(config)
	case "sqlite":
		return initSQLite(config)
	case mongoDBType:
		return initMongoDB(config)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}
}

// initPostgreSQL initializes a PostgreSQL connection
func initPostgreSQL(config *DatabaseConfig) (*Connection, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		config.Host, config.User, config.Password, config.Name, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	return &Connection{
		Type:   "postgres",
		GormDB: db,
	}, nil
}

// initSQLite initializes a SQLite connection
func initSQLite(config *DatabaseConfig) (*Connection, error) {
	db, err := gorm.Open(sqlite.Open(config.Name), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	return &Connection{
		Type:   "sqlite",
		GormDB: db,
	}, nil
}

// initMongoDB initializes a MongoDB connection
func initMongoDB(config *DatabaseConfig) (*Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%d", config.Host, config.Port)
	if config.User != "" && config.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d", config.User, config.Password, config.Host, config.Port)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(config.Name)

	return &Connection{
		Type:        mongoDBType,
		MongoDB:     database,
		MongoClient: client,
	}, nil
}
