package gobase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

const (
	mongoDBType = "mongodb"
)

// Accessor implements both ModelAccessor and MigrationProvider interfaces.
// This follows the Interface Segregation Principle by implementing
// focused interfaces, and the Single Responsibility Principle by
// handling only data access operations.
type Accessor struct {
	connection *Connection
}

// NewAccessor creates a new Accessor instance with the provided database connection.
// This follows the Dependency Injection pattern, supporting the Dependency
// Inversion Principle.
func NewAccessor(connection *Connection) *Accessor {
	return &Accessor{connection: connection}
}

// ValidateModel checks if the model properly embeds BaseModel
func (a *Accessor) ValidateModel(model interface{}) error {
	return ValidateBaseModel(model)
}

// Create inserts a new record into the database.
// This method follows the Single Responsibility Principle by only
// handling record creation. Uses Django-style naming.
func (a *Accessor) Create(model interface{}) error {
	if err := a.ValidateModel(model); err != nil {
		return fmt.Errorf("model validation failed: %w", err)
	}

	// Only support GORM for now (SQLite/PostgreSQL)
	if a.connection.Type == mongoDBType {
		return errors.New("MongoDB support not yet implemented for Create operation")
	}

	result := a.connection.GormDB.Create(model)
	return result.Error
}

// Get retrieves a record by its ID and populates the provided model.
// This method follows Django-style naming (Get instead of ReadByID).
func (a *Accessor) Get(model interface{}, id interface{}) error {
	if err := a.ValidateModel(model); err != nil {
		return fmt.Errorf("model validation failed: %w", err)
	}

	if id == nil {
		return errors.New("id cannot be nil")
	}

	// Only support GORM for now (SQLite/PostgreSQL)
	if a.connection.Type == mongoDBType {
		return errors.New("MongoDB support not yet implemented for Get operation")
	}

	// Handle both numeric and string IDs properly
	result := a.connection.GormDB.Where("id = ?", id).First(model)
	return result.Error
}

// All retrieves all records and populates the provided slice.
// This method follows Django-style naming (All instead of FindAll).
func (a *Accessor) All(models interface{}) error {
	if models == nil {
		return errors.New("models cannot be nil")
	}

	// Validate that models is a slice
	rv := reflect.ValueOf(models)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return errors.New("models must be a pointer to a slice")
	}

	// Get the element type of the slice and validate it embeds BaseModel
	sliceType := rv.Elem().Type()
	elemType := sliceType.Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	// Create a temporary instance to validate
	tempModel := reflect.New(elemType).Interface()
	if err := a.ValidateModel(tempModel); err != nil {
		return fmt.Errorf("model validation failed: %w", err)
	}

	// Only support GORM for now (SQLite/PostgreSQL)
	if a.connection.Type == mongoDBType {
		return errors.New("MongoDB support not yet implemented for All operation")
	}

	result := a.connection.GormDB.Find(models)
	return result.Error
}

// Filter retrieves records based on conditions. Django-style filtering.
func (a *Accessor) Filter(models interface{}, conditions map[string]interface{}) error {
	if models == nil {
		return errors.New("models cannot be nil")
	}

	if len(conditions) == 0 {
		return a.All(models)
	}

	// Validate that models is a slice
	rv := reflect.ValueOf(models)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Slice {
		return errors.New("models must be a pointer to a slice")
	}

	// Only support GORM for now (SQLite/PostgreSQL)
	if a.connection.Type == mongoDBType {
		return errors.New("MongoDB support not yet implemented for Filter operation")
	}

	query := a.connection.GormDB
	for field, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	result := query.Find(models)
	return result.Error
}

// Update saves changes to an existing record.
// This method follows the Single Responsibility Principle by only
// handling record updates.
func (a *Accessor) Update(model interface{}) error {
	if err := a.ValidateModel(model); err != nil {
		return fmt.Errorf("model validation failed: %w", err)
	}

	// Only support GORM for now (SQLite/PostgreSQL)
	if a.connection.Type == mongoDBType {
		return errors.New("MongoDB support not yet implemented for Update operation")
	}

	result := a.connection.GormDB.Save(model)
	return result.Error
}

// Delete performs a soft delete on the record.
// This method follows the Single Responsibility Principle by only
// handling record deletion.
func (a *Accessor) Delete(model interface{}) error {
	if err := a.ValidateModel(model); err != nil {
		return fmt.Errorf("model validation failed: %w", err)
	}

	// Only support GORM for now (SQLite/PostgreSQL)
	if a.connection.Type == mongoDBType {
		return errors.New("MongoDB support not yet implemented for Delete operation")
	}

	result := a.connection.GormDB.Delete(model)
	return result.Error
}

// AutoMigrate automatically migrates the schema for all registered models.
// This method has been modified to auto-discover models that embed BaseModel.
// Migrate performs database schema migration for the provided models.
// If no models are provided, it will migrate all registered models.
// The default User model is only migrated if explicitly used or passed as an argument.
func (a *Accessor) Migrate(models ...interface{}) error {
	// Only support GORM for now (SQLite/PostgreSQL)
	if a.connection.Type == mongoDBType {
		return errors.New("MongoDB support not yet implemented for Migrate operation")
	}

	var modelsToMigrate []interface{}

	if len(models) > 0 {
		// Use explicitly provided models
		modelsToMigrate = models
	} else {
		// Use registered models, but exclude User unless it's being used
		registeredModels := GetRegisteredModels()
		for _, model := range registeredModels {
			// Check if this is the default User model
			if isDefaultUserModel(model) {
				// Only include if explicitly used (this is a simplified check)
				// In a real implementation, you might want to scan the code or use reflection
				// to determine if the User model is actually being used
				continue
			}
			modelsToMigrate = append(modelsToMigrate, model)
		}
	}

	if len(modelsToMigrate) == 0 {
		return errors.New("no models to migrate")
	}

	return a.connection.GormDB.AutoMigrate(modelsToMigrate...)
}

// isDefaultUserModel checks if the given model is the default gobase.User model
func isDefaultUserModel(model interface{}) bool {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.PkgPath() == "gobase" && t.Name() == "User"
}

// Preload loads data from JSON files into the database
func (a *Accessor) Preload(modelRegistry map[string]interface{}, jsonFilePaths ...string) error {
	for _, filePath := range jsonFilePaths {
		err := a.preloadFromFile(modelRegistry, filePath)
		if err != nil {
			return fmt.Errorf("failed to preload from %s: %w", filePath, err)
		}
	}
	return nil
}

// preloadFromFile loads data from a single JSON file
func (a *Accessor) preloadFromFile(modelRegistry map[string]interface{}, filePath string) error {
	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Determine model type from filename
	filename := filepath.Base(filePath)
	modelName := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Find the model type in the registry
	modelTemplate, exists := modelRegistry[modelName]
	if !exists {
		return fmt.Errorf("no model registered for %s", modelName)
	}

	// Parse JSON array
	var jsonObjects []map[string]interface{}
	err = json.Unmarshal(data, &jsonObjects)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Process each object
	for _, objData := range jsonObjects {
		err = a.preloadObject(modelTemplate, objData)
		if err != nil {
			return fmt.Errorf("failed to preload object: %w", err)
		}
	}

	return nil
}

// preloadObject processes a single JSON object and creates/updates the database record
func (a *Accessor) preloadObject(modelTemplate interface{}, objData map[string]interface{}) error {
	// Create a new instance of the model
	modelType := reflect.TypeOf(modelTemplate)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	newModel := reflect.New(modelType).Interface()

	// Convert JSON data to the model struct
	jsonBytes, err := json.Marshal(objData)
	if err != nil {
		return fmt.Errorf("failed to marshal object data: %w", err)
	}

	err = json.Unmarshal(jsonBytes, newModel)
	if err != nil {
		return fmt.Errorf("failed to unmarshal to model: %w", err)
	}

	// Check if record already exists to prevent duplicates
	exists, err := a.recordExists(newModel, objData)
	if err != nil {
		return fmt.Errorf("failed to check if record exists: %w", err)
	}

	if exists {
		// Update existing record
		return a.Update(newModel)
	} else {
		// Create new record
		return a.Create(newModel)
	}
}

// recordExists checks if a record already exists based on unique fields
func (a *Accessor) recordExists(model interface{}, objData map[string]interface{}) (bool, error) {
	// For simplicity, we'll check by ID if provided, otherwise by unique fields
	if id, hasID := objData["id"]; hasID {
		tempModel := reflect.New(reflect.TypeOf(model).Elem()).Interface()
		err := a.Get(tempModel, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, err
		}
		// Copy the found ID to our model for updating
		a.copyID(tempModel, model)
		return true, nil
	}

	// If no ID, assume it's a new record for now
	// In a production system, you'd check other unique fields
	return false, nil
}

// copyID copies the ID from source to destination model
func (a *Accessor) copyID(source, dest interface{}) {
	srcVal := reflect.ValueOf(source).Elem()
	destVal := reflect.ValueOf(dest).Elem()

	// Look for ID field in BaseModel
	srcID := srcVal.FieldByName("ID")
	destID := destVal.FieldByName("ID")

	if srcID.IsValid() && destID.IsValid() && destID.CanSet() {
		destID.Set(srcID)
	}
}

// Preload is a convenience function for preloading data from JSON files
func Preload(accessor *Accessor, modelRegistry map[string]interface{}, jsonFilePaths ...string) error {
	return accessor.Preload(modelRegistry, jsonFilePaths...)
}

// Advanced query methods that extend functionality while maintaining SOLID principles

// FindWhere retrieves records based on a WHERE clause.
func (a *Accessor) FindWhere(models interface{}, condition string, args ...interface{}) error {
	if models == nil {
		return errors.New("models cannot be nil")
	}

	// Only support GORM for now (SQLite/PostgreSQL)
	if a.connection.Type == mongoDBType {
		return errors.New("MongoDB support not yet implemented for FindWhere operation")
	}

	result := a.connection.GormDB.Where(condition, args...).Find(models)
	return result.Error
}

// Count returns the number of records matching the given conditions.
func (a *Accessor) Count(model interface{}, condition string, args ...interface{}) (int64, error) {
	// Only support GORM for now (SQLite/PostgreSQL)
	if a.connection.Type == mongoDBType {
		return 0, errors.New("MongoDB support not yet implemented for Count operation")
	}

	var count int64
	result := a.connection.GormDB.Model(model).Where(condition, args...).Count(&count)
	return count, result.Error
}

// Transaction executes a function within a database transaction.
// This follows the Single Responsibility Principle by handling only transaction management.
func (a *Accessor) Transaction(fn func(*Accessor) error) error {
	// Only support GORM for now (SQLite/PostgreSQL)
	if a.connection.Type == mongoDBType {
		return errors.New("MongoDB support not yet implemented for Transaction operation")
	}

	return a.connection.GormDB.Transaction(func(tx *gorm.DB) error {
		txConnection := &Connection{
			Type:   a.connection.Type,
			GormDB: tx,
		}
		txAccessor := NewAccessor(txConnection)
		return fn(txAccessor)
	})
}
