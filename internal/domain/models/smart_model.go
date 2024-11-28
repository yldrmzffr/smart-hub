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
	ID           uuid.UUID              `json:"id" db:"id" validate:"omitempty,uuid"`
	Name         string                 `json:"name" db:"name" validate:"required,min=2,max=255"`
	Description  string                 `json:"description" db:"description" validate:"required,max=1000"`
	Type         ModelType              `json:"type" db:"type" validate:"required,lowercase,oneof=device service"`
	Category     ModelCategory          `json:"category" db:"category" validate:"required,lowercase,oneof=wearable camera weather entertainment"`
	Manufacturer string                 `json:"manufacturer,omitempty" db:"manufacturer" validate:"omitempty,max=255"`
	ModelNumber  string                 `json:"model_number,omitempty" db:"model_number" validate:"omitempty,max=50,alphanum"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata" validate:"omitempty,dive,keys,required,endkeys"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at" validate:"omitempty,ltefield=UpdatedAt"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at" validate:"omitempty"`
}
