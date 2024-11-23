package models

import (
	"time"

	"github.com/google/uuid"
)

type ProtocolType string

const (
	RestProtocol      ProtocolType = "rest"
	GrpcProtocol      ProtocolType = "grpc"
	MqttProtocol      ProtocolType = "mqtt"
	WebsocketProtocol ProtocolType = "websocket"
)

type SmartFeature struct {
	ID            uuid.UUID              `json:"id" db:"id"`
	ModelID       uuid.UUID              `json:"model_id" db:"model_id" validate:"required"`
	Name          string                 `json:"name" db:"name" validate:"required,min=2,max=255"`
	Description   string                 `json:"description" db:"description" validate:"required"`
	Protocol      ProtocolType           `json:"protocol" db:"protocol" validate:"required,oneof=rest grpc mqtt websocket"`
	InterfacePath string                 `json:"interface_path" db:"interface_path" validate:"required"`
	Parameters    map[string]interface{} `json:"parameters,omitempty" db:"parameters"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`
}
