package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"smart-hub/internal/common/logger"
	"sync"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}

func ValidateStruct(s interface{}) error {
	err := GetValidator().Struct(s)
	if err != nil {
		logger.Error("Validation error", "error", err)
	}

	return err
}

func ValidateUUID(s string) error {
	_, err := uuid.Parse(s)
	if err != nil {
		logger.Error("Invalid UUID", "error", err)
		return err
	}
	return nil
}
