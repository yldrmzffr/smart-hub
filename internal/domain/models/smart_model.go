package models

import (
	"time"

	"github.com/google/uuid"
)

type ModelType string
type ModelCategory string

const (
	DeviceType  ModelType = "device"
	ServiceType ModelType = "service"

	WearableCategory      ModelCategory = "wearable"
	CameraCategory        ModelCategory = "camera"
	WeatherCategory       ModelCategory = "weather"
	EntertainmentCategory ModelCategory = "entertainment"
)

type SmartModel struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	Name         string                 `json:"name" db:"name" validate:"required,min=2,max=255"`
	Description  string                 `json:"description" db:"description" validate:"required"`
	Type         ModelType              `json:"type" db:"type" validate:"required,oneof=device service"`
	Category     ModelCategory          `json:"category" db:"category" validate:"required"`
	Manufacturer string                 `json:"manufacturer,omitempty" db:"manufacturer"`
	ModelNumber  string                 `json:"model_number,omitempty" db:"model_number"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
}
