package service

import (
	"context"
	"smart-hub/internal/common/logger"
	"smart-hub/internal/domain/interfaces"
	"smart-hub/internal/domain/models"
)

type SmartFeatureService struct {
	repo interfaces.SmartFeatureRepository
}

func NewSmartFeatureService(repo interfaces.SmartFeatureRepository) *SmartFeatureService {
	return &SmartFeatureService{
		repo: repo,
	}
}

func (s *SmartFeatureService) Create(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error) {
	logger.Debug("Create smart feature", "feature", feature)
	return s.repo.Create(ctx, feature)
}

func (s *SmartFeatureService) GetByID(ctx context.Context, id string) (*models.SmartFeature, error) {
	logger.Debug("Get smart feature by ID", "id", id)
	return s.repo.GetByID(ctx, id)
}

func (s *SmartFeatureService) GetWithModelID(ctx context.Context, modelID string) ([]*models.SmartFeature, error) {
	logger.Debug("Get smart feature by model ID", "modelID", modelID)
	return s.repo.GetWithModelID(ctx, modelID)
}

func (s *SmartFeatureService) GetAll(ctx context.Context) ([]*models.SmartFeature, error) {
	logger.Debug("Get all smart features")
	return s.repo.GetAll(ctx)
}

func (s *SmartFeatureService) Update(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error) {
	logger.Debug("Update smart feature", "feature", feature)
	return s.repo.Update(ctx, feature)
}

func (s *SmartFeatureService) Delete(ctx context.Context, id string) error {
	logger.Debug("Delete smart feature", "id", id)
	return s.repo.Delete(ctx, id)
}
