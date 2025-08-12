package gobase

// ModelAccessor defines the contract for CRUD operations on database models using Django-style naming.
// This interface follows the Interface Segregation Principle by only including
// the essential CRUD operations that all model accessors should support.
type ModelAccessor interface {
	Create(model interface{}) error
	Get(model interface{}, id interface{}) error // Django-style: Get instead of ReadByID
	All(models interface{}) error                // Django-style: All instead of FindAll
	Update(model interface{}) error
	Delete(model interface{}) error
	Filter(models interface{}, conditions map[string]interface{}) error // Django-style filtering
}

// MigrationProvider defines the contract for database schema migration operations.
// This is separated from ModelAccessor following the Interface Segregation Principle
// as not all accessors may need migration capabilities.
type MigrationProvider interface {
	Migrate(models ...interface{}) error // Renamed from AutoMigrate and supports manual model registration
}

// DatabaseConnection defines the contract for database connectivity.
// This abstraction follows the Dependency Inversion Principle by allowing
// different database implementations while maintaining the same interface.
type DatabaseConnection interface {
	GetDB() interface{}
	Close() error
}

// ModelValidator defines the contract for validating models.
type ModelValidator interface {
	ValidateModel(model interface{}) error
}
