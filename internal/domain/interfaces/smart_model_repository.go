package interfaces

import (
	"context"
	"smart-hub/internal/domain/models"
)

type SmartModelRepository interface {
	Create(ctx context.Context, model *models.SmartModel) (*models.SmartModel, error)
	GetByID(ctx context.Context, id int) (*models.SmartModel, error)
	GetWithType(ctx context.Context, modelType string) ([]*models.SmartModel, error)
	GetAll(ctx context.Context) ([]*models.SmartModel, error)
	Update(ctx context.Context, model *models.SmartModel) (*models.SmartModel, error)
	Delete(ctx context.Context, id int) error
}
