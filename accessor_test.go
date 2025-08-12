package gobase

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestModel for testing
type TestModel struct {
	BaseModel
	Name string `json:"name"`
}

// CustomIDModel for testing field overriding
type CustomIDModel struct {
	BaseModel
	ID   string `gorm:"primarykey;size:36" json:"id"`
	Name string `json:"name"`
}

func (c *CustomIDModel) GetID() interface{} {
	return c.ID
}

func (c *CustomIDModel) SetID(id interface{}) {
	if id, ok := id.(string); ok {
		c.ID = id
	}
}

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) *Connection {
	// Create temporary SQLite database for testing
	config := &DatabaseConfig{
		Type: "sqlite",
		Name: ":memory:",
	}

	connection, err := InitDBWithConfig(config)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	return connection
}

// TestDatabaseConfiguration tests the database configuration loading
func TestDatabaseConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		expectedDB  string
	}{
		{
			name: "Valid SQLite config",
			envVars: map[string]string{
				"DB_TYPE": "sqlite",
				"DB_NAME": "test.db",
			},
			expectError: false,
			expectedDB:  "sqlite",
		},
		{
			name: "Valid PostgreSQL config",
			envVars: map[string]string{
				"DB_TYPE":     "postgres",
				"DB_HOST":     "localhost",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_PORT":     "5432",
			},
			expectError: false,
			expectedDB:  "postgres",
		},
		{
			name: "Missing DB_TYPE",
			envVars: map[string]string{
				"DB_NAME": "test.db",
			},
			expectError: true,
		},
		{
			name: "Invalid DB_TYPE",
			envVars: map[string]string{
				"DB_TYPE": "invalid",
				"DB_NAME": "test.db",
			},
			expectError: true,
		},
		{
			name: "Missing DB_NAME",
			envVars: map[string]string{
				"DB_TYPE": "sqlite",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config, err := LoadConfig()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if config.Type != tt.expectedDB {
					t.Errorf("Expected DB type %s, got %s", tt.expectedDB, config.Type)
				}
			}
		})
	}
}

// TestModelValidation tests BaseModel validation
func TestModelValidation(t *testing.T) {
	accessor := NewAccessor(setupTestDB(t))

	tests := []struct {
		name        string
		model       interface{}
		expectError bool
	}{
		{
			name:        "Valid model with BaseModel",
			model:       &TestModel{Name: "Test"},
			expectError: false,
		},
		{
			name:        "Nil model",
			model:       nil,
			expectError: true,
		},
		{
			name:        "Model without BaseModel",
			model:       &struct{ Name string }{Name: "Test"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := accessor.ValidateModel(tt.model)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// TestAccessor_Create tests the Create method
func TestAccessor_Create(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Register test model
	RegisterModel(&TestModel{})

	err := accessor.Migrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	model := &TestModel{Name: "Test Model"}

	err = accessor.Create(model)
	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	if model.ID == 0 {
		t.Error("Expected ID to be set after creation")
	}

	if model.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set after creation")
	}
}

// TestAccessor_Get tests the Get method (Django-style)
func TestAccessor_Get(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Register test model
	RegisterModel(&TestModel{})

	err := accessor.Migrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create a model
	original := &TestModel{Name: "Original Model"}
	err = accessor.Create(original)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// Get it back using Django-style Get method
	read := &TestModel{}
	err = accessor.Get(read, original.ID)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}

	if read.Name != original.Name {
		t.Errorf("Expected name %s, got %s", original.Name, read.Name)
	}

	if read.ID != original.ID {
		t.Errorf("Expected ID %d, got %d", original.ID, read.ID)
	}
}

// TestAccessor_All tests the All method (Django-style)
func TestAccessor_All(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Register test model
	RegisterModel(&TestModel{})

	err := accessor.Migrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create multiple models
	models := []*TestModel{
		{Name: "Model 1"},
		{Name: "Model 2"},
		{Name: "Model 3"},
	}

	for _, model := range models {
		err = accessor.Create(model)
		if err != nil {
			t.Fatalf("Failed to create model: %v", err)
		}
	}

	// Get all using Django-style All method
	var found []TestModel
	err = accessor.All(&found)
	if err != nil {
		t.Errorf("All failed: %v", err)
	}

	if len(found) != 3 {
		t.Errorf("Expected 3 models, got %d", len(found))
	}
}

// TestAccessor_Filter tests the Filter method (Django-style)
func TestAccessor_Filter(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Register test model
	RegisterModel(&TestModel{})

	err := accessor.Migrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create test data
	models := []*TestModel{
		{Name: "Alpha"},
		{Name: "Beta"},
		{Name: "Alpha"},
	}

	for _, model := range models {
		err = accessor.Create(model)
		if err != nil {
			t.Fatalf("Failed to create model: %v", err)
		}
	}

	// Filter using Django-style Filter method
	var filtered []TestModel
	err = accessor.Filter(&filtered, map[string]interface{}{"name": "Alpha"})
	if err != nil {
		t.Errorf("Filter failed: %v", err)
	}

	if len(filtered) != 2 {
		t.Errorf("Expected 2 filtered models, got %d", len(filtered))
	}

	for _, model := range filtered {
		if model.Name != "Alpha" {
			t.Errorf("Expected name 'Alpha', got %s", model.Name)
		}
	}
}

// TestAccessor_Update tests the Update method
func TestAccessor_Update(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Register test model
	RegisterModel(&TestModel{})

	err := accessor.Migrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create a model
	model := &TestModel{Name: "Original Name"}
	err = accessor.Create(model)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	originalUpdatedAt := model.UpdatedAt

	// Small delay to ensure UpdatedAt changes
	time.Sleep(time.Millisecond * 10)

	// Update it
	model.Name = "Updated Name"
	err = accessor.Update(model)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// Get it back to verify
	read := &TestModel{}
	err = accessor.Get(read, model.ID)
	if err != nil {
		t.Fatalf("Failed to read updated model: %v", err)
	}

	if read.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", read.Name)
	}

	if !read.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

// TestAccessor_Delete tests the Delete method (soft delete)
func TestAccessor_Delete(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Register test model
	RegisterModel(&TestModel{})

	err := accessor.Migrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create a model
	model := &TestModel{Name: "To Delete"}
	err = accessor.Create(model)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// Delete it
	err = accessor.Delete(model)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Try to read it back (should fail with soft delete)
	read := &TestModel{}
	err = accessor.Get(read, model.ID)
	if err == nil {
		t.Error("Expected error when reading soft-deleted model")
	}

	// Verify it still exists with Unscoped
	var count int64
	connection.GormDB.Unscoped().Model(&TestModel{}).Where("id = ?", model.ID).Count(&count)
	if count != 1 {
		t.Error("Expected soft-deleted model to still exist in database")
	}
}

// TestMigrate tests manual and automatic model migration
func TestMigrate(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Test manual migration with explicit models
	err := accessor.Migrate(&TestModel{}, &User{})
	if err != nil {
		t.Errorf("Manual Migrate failed: %v", err)
	}

	// Verify tables were created
	hasTestTable := connection.GormDB.Migrator().HasTable(&TestModel{})
	hasUserTable := connection.GormDB.Migrator().HasTable(&User{})

	if !hasTestTable {
		t.Error("TestModel table was not created")
	}

	if !hasUserTable {
		t.Error("User table was not created")
	}

	// Test migration with registered models
	globalModelRegistry.models = []interface{}{}
	RegisterModel(&TestModel{})

	// Create new connection for clean test
	connection2 := setupTestDB(t)
	accessor2 := NewAccessor(connection2)

	err = accessor2.Migrate()
	if err != nil {
		t.Errorf("Registered model Migrate failed: %v", err)
	}

	// Verify registered model table was created
	hasTestTable2 := connection2.GormDB.Migrator().HasTable(&TestModel{})
	if !hasTestTable2 {
		t.Error("RegisteredTestModel table was not created")
	}
} // TestCustomFieldOverride tests field overriding capability
func TestCustomFieldOverride(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Migrate specific model
	err := accessor.Migrate(&CustomIDModel{})
	if err != nil {
		t.Fatalf("Failed to migrate CustomIDModel: %v", err)
	}

	// Test creating a model with custom string ID
	customModel := &CustomIDModel{
		ID:   "custom-001",
		Name: "Custom Model",
	}
	err = accessor.Create(customModel)
	if err != nil {
		t.Fatalf("Failed to create model with custom ID: %v", err)
	}

	// Test reading the model
	var retrieved CustomIDModel
	err = accessor.Get(&retrieved, "custom-001")
	if err != nil {
		t.Fatalf("Failed to read model with custom ID: %v", err)
	}

	// Verify the custom ID field
	if retrieved.ID != "custom-001" {
		t.Errorf("Expected ID 'custom-001', got %s", retrieved.ID)
	}

	if retrieved.Name != "Custom Model" {
		t.Errorf("Expected name 'Custom Model', got %s", retrieved.Name)
	}
}

// TestPreload tests the JSON data preloading functionality
func TestPreload(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Migrate models
	err := accessor.Migrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create temporary JSON files
	testData := []map[string]interface{}{
		{"name": "Test Item 1"},
		{"name": "Test Item 2"},
		{"name": "Test Item 3"},
	}

	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	// Write to temporary file with proper name
	tmpDir, err := os.MkdirTemp("", "gobase_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test_models.json")
	err = os.WriteFile(tmpFile, jsonData, 0o644)
	if err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// Define model registry
	modelRegistry := map[string]interface{}{
		"test_models": &TestModel{},
	}

	// Test preload
	err = accessor.Preload(modelRegistry, tmpFile)
	if err != nil {
		t.Fatalf("Preload failed: %v", err)
	}

	// Verify data was loaded
	var models []TestModel
	err = accessor.All(&models)
	if err != nil {
		t.Fatalf("Failed to retrieve models: %v", err)
	}

	if len(models) != 3 {
		t.Errorf("Expected 3 models, got %d", len(models))
	}

	// Verify content
	expectedNames := []string{"Test Item 1", "Test Item 2", "Test Item 3"}
	for i, model := range models {
		if model.Name != expectedNames[i] {
			t.Errorf("Expected name '%s', got '%s'", expectedNames[i], model.Name)
		}
	}
}

// TestPreloadDuplicatePrevention tests that preload prevents duplicates
func TestPreloadDuplicatePrevention(t *testing.T) {
	connection := setupTestDB(t)
	accessor := NewAccessor(connection)

	// Migrate models
	err := accessor.Migrate(&TestModel{})
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Create test data with explicit IDs
	testData := []map[string]interface{}{
		{"id": 1, "name": "Test Item 1"},
		{"id": 2, "name": "Test Item 2"},
	}

	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	// Write to temporary file with proper name
	tmpDir, err := os.MkdirTemp("", "gobase_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test_models.json")
	err = os.WriteFile(tmpFile, jsonData, 0o644)
	if err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// Define model registry
	modelRegistry := map[string]interface{}{
		"test_models": &TestModel{},
	}

	// First preload
	err = accessor.Preload(modelRegistry, tmpFile)
	if err != nil {
		t.Fatalf("First preload failed: %v", err)
	}

	// Second preload (should update, not duplicate)
	err = accessor.Preload(modelRegistry, tmpFile)
	if err != nil {
		t.Fatalf("Second preload failed: %v", err)
	}

	// Verify no duplicates
	var models []TestModel
	err = accessor.All(&models)
	if err != nil {
		t.Fatalf("Failed to retrieve models: %v", err)
	}

	if len(models) != 2 {
		t.Errorf("Expected 2 models (no duplicates), got %d", len(models))
	}
}
