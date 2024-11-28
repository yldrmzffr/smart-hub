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
	ID            uuid.UUID              `json:"id" db:"id" validate:"omitempty,uuid"`
	ModelID       uuid.UUID              `json:"model_id" db:"model_id" validate:"uuid"`
	Name          string                 `json:"name" db:"name" validate:"required,min=2,max=255"`
	Description   string                 `json:"description" db:"description" validate:"required,max=1000"`
	Protocol      ProtocolType           `json:"protocol" db:"protocol" validate:"required,oneof=rest grpc mqtt websocket"`
	InterfacePath string                 `json:"interface_path" db:"interface_path" validate:"required,startswith=/"`
	Parameters    map[string]interface{} `json:"parameters,omitempty" db:"parameters" validate:"omitempty,dive,keys,required,endkeys"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at" validate:"omitempty,ltefield=UpdatedAt"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at" validate:"omitempty"`
}
