package service

import (
	"context"
	"smart-hub/internal/common/logger"
	"smart-hub/internal/domain/interfaces"
	"smart-hub/internal/domain/models"
)

type SmartModelService struct {
	repo interfaces.SmartModelRepository
}

func NewSmartModelService(repo interfaces.SmartModelRepository) *SmartModelService {
	return &SmartModelService{
		repo: repo,
	}
}

func (s *SmartModelService) Create(ctx context.Context, model *models.SmartModel) (*models.SmartModel, error) {
	logger.Debug("Create smart model", "model", model)
	return s.repo.Create(ctx, model)
}

func (s *SmartModelService) GetByID(ctx context.Context, id string) (*models.SmartModel, error) {
	logger.Debug("Get smart model by ID", "id", id)
	return s.repo.GetByID(ctx, id)
}

func (s *SmartModelService) GetWithModelID(ctx context.Context, modelID string) ([]*models.SmartModel, error) {
	logger.Debug("Get smart model by model ID", "modelID", modelID)
	return s.repo.GetWithType(ctx, modelID)
}

func (s *SmartModelService) GetAll(ctx context.Context) ([]*models.SmartModel, error) {
	logger.Debug("Get all smart models")
	return s.repo.GetAll(ctx)
}

func (s *SmartModelService) Update(ctx context.Context, model *models.SmartModel) (*models.SmartModel, error) {
	logger.Debug("Update smart model", "model", model)
	return s.repo.Update(ctx, model)
}

func (s *SmartModelService) Delete(ctx context.Context, id string) error {
	logger.Debug("Delete smart model", "id", id)
	return s.repo.Delete(ctx, id)
}
