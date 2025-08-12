package gobase

import (
	"errors"
	"reflect"
	"time"

	"gorm.io/gorm"
)

// BaseModel provides standard fields that can be embedded in user models.
// These fields can be overridden by defining fields with the same name
// in the embedding struct, following Go's embedding rules.
// This design follows the Open/Closed Principle - open for extension
// (by embedding) but closed for modification.
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// GetID returns the ID of the model. This method can be overridden
// by embedding structs to provide custom ID logic.
func (bm *BaseModel) GetID() interface{} {
	return bm.ID
}

// SetID sets the ID of the model. This method can be overridden
// by embedding structs to provide custom ID logic.
func (bm *BaseModel) SetID(id interface{}) {
	if id, ok := id.(uint); ok {
		bm.ID = id
	}
}

// IsDeleted checks if the model is soft-deleted.
func (bm *BaseModel) IsDeleted() bool {
	return bm.DeletedAt.Valid
}

// ModelRegistry keeps track of all registered models that embed BaseModel
type ModelRegistry struct {
	models []interface{}
}

// Global model registry
var globalModelRegistry = &ModelRegistry{
	models: []interface{}{
		&User{}, // Register the default User model
	},
}

// RegisterModel adds a model to the global registry for auto-migration
func RegisterModel(model interface{}) {
	globalModelRegistry.models = append(globalModelRegistry.models, model)
}

// GetRegisteredModels returns all registered models
func GetRegisteredModels() []interface{} {
	return globalModelRegistry.models
}

// ValidateBaseModel checks if a model properly embeds BaseModel
func ValidateBaseModel(model interface{}) error {
	if model == nil {
		return errors.New("model cannot be nil")
	}

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		return errors.New("model must be a struct")
	}

	// Check if BaseModel is embedded
	hasBaseModel := false
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Anonymous && field.Type.Name() == "BaseModel" {
			hasBaseModel = true
			break
		}
	}

	if !hasBaseModel {
		return errors.New("model must embed gobase.BaseModel")
	}

	return nil
}

// IsModelRegistered checks if a model type is registered
func IsModelRegistered(model interface{}) bool {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for _, registeredModel := range globalModelRegistry.models {
		registeredType := reflect.TypeOf(registeredModel)
		if registeredType.Kind() == reflect.Ptr {
			registeredType = registeredType.Elem()
		}

		if modelType == registeredType {
			return true
		}
	}
	return false
}
