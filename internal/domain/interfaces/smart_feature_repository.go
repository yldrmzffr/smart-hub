package interfaces

import (
	"context"
	"smart-hub/internal/domain/models"
)

type SmartFeatureRepository interface {
	Create(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error)
	GetByID(ctx context.Context, id int) (*models.SmartFeature, error)
	GetWithModelID(ctx context.Context, modelID int) ([]*models.SmartFeature, error)
	GetAll(ctx context.Context) ([]*models.SmartFeature, error)
	Update(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error)
	Delete(ctx context.Context, id int) error
}
