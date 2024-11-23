package interfaces

import (
	"context"
	"smart-hub/internal/domain/models"
)

type SmartFeatureService interface {
	Create(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error)
	GetByID(ctx context.Context, id string) (*models.SmartFeature, error)
	GetWithModelID(ctx context.Context, modelID string) ([]*models.SmartFeature, error)
	GetAll(ctx context.Context) ([]*models.SmartFeature, error)
	Update(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error)
	Delete(ctx context.Context, id string) error
}
