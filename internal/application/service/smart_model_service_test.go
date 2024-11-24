package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"smart-hub/internal/domain/models"
	"testing"
	"time"
)

type mockSmartModelRepo struct {
	mock.Mock
}

func (m *mockSmartModelRepo) Create(ctx context.Context, model *models.SmartModel) (*models.SmartModel, error) {
	args := m.Called(ctx, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartModel), args.Error(1)
}

func (m *mockSmartModelRepo) GetByID(ctx context.Context, id string) (*models.SmartModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartModel), args.Error(1)
}

func (m *mockSmartModelRepo) GetWithType(ctx context.Context, modelType models.ModelType) ([]*models.SmartModel, error) {
	args := m.Called(ctx, modelType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SmartModel), args.Error(1)
}

func (m *mockSmartModelRepo) GetAll(ctx context.Context) ([]*models.SmartModel, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SmartModel), args.Error(1)
}

func (m *mockSmartModelRepo) Update(ctx context.Context, model *models.SmartModel) (*models.SmartModel, error) {
	args := m.Called(ctx, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartModel), args.Error(1)
}

func (m *mockSmartModelRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSmartModelService_Create(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	now := time.Now()
	testModel := &models.SmartModel{
		ID:          uuid.New(),
		Name:        "Test Device",
		Description: "Test Description",
		Type:        models.DeviceType,
		Category:    models.WearableCategory,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("Create", mock.Anything, testModel).Return(testModel, nil)

	result, err := service.Create(context.Background(), testModel)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testModel.Name, result.Name)
	assert.Equal(t, testModel.Description, result.Description)

	mockRepo.AssertExpectations(t)
}

func TestSmartModelService_Create_Error(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	now := time.Now()
	testModel := &models.SmartModel{
		ID:          uuid.New(),
		Name:        "Test Device",
		Description: "Test Description",
		Type:        models.DeviceType,
		Category:    models.WearableCategory,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("Create", mock.Anything, testModel).Return(nil, assert.AnError)

	result, err := service.Create(context.Background(), testModel)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestSmartModelService_GetByID(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	now := time.Now()
	testModel := &models.SmartModel{
		ID:          uuid.New(),
		Name:        "Test Device",
		Description: "Test Description",
		Type:        models.DeviceType,
		Category:    models.WearableCategory,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("GetByID", mock.Anything, testModel.ID.String()).Return(testModel, nil)

	result, err := service.GetByID(context.Background(), testModel.ID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testModel.Name, result.Name)
	assert.Equal(t, testModel.Description, result.Description)

	mockRepo.AssertExpectations(t)
}

func TestSmartModelService_GetByID_Error(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	now := time.Now()
	testModel := &models.SmartModel{
		ID:          uuid.New(),
		Name:        "Test Device",
		Description: "Test Description",
		Type:        models.DeviceType,
		Category:    models.WearableCategory,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("GetByID", mock.Anything, testModel.ID.String()).Return(nil, assert.AnError)

	result, err := service.GetByID(context.Background(), testModel.ID.String())

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestSmartModelService_GetWithType(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	now := time.Now()
	testModel := &models.SmartModel{
		ID:          uuid.New(),
		Name:        "Test Device",
		Description: "Test Description",
		Type:        models.DeviceType,
		Category:    models.WearableCategory,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("GetWithType", mock.Anything, models.DeviceType).Return([]*models.SmartModel{testModel}, nil)

	result, err := service.GetWithType(context.Background(), models.DeviceType)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testModel.Name, result[0].Name)
	assert.Equal(t, testModel.Description, result[0].Description)

	mockRepo.AssertExpectations(t)
}

func TestSmartModelService_GetWithType_Error(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	mockRepo.On("GetWithType", mock.Anything, models.DeviceType).Return(nil, assert.AnError)

	result, err := service.GetWithType(context.Background(), models.DeviceType)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestSmartModelService_GetAll(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	now := time.Now()
	testModel := &models.SmartModel{
		ID:          uuid.New(),
		Name:        "Test Device",
		Description: "Test Description",
		Type:        models.DeviceType,
		Category:    models.WearableCategory,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("GetAll", mock.Anything).Return([]*models.SmartModel{testModel}, nil)

	result, err := service.GetAll(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testModel.Name, result[0].Name)
	assert.Equal(t, testModel.Description, result[0].Description)

	mockRepo.AssertExpectations(t)
}

func TestSmartModelService_GetAll_Error(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	mockRepo.On("GetAll", mock.Anything).Return(nil, assert.AnError)

	result, err := service.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestSmartModelService_Update(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	now := time.Now()
	testModel := &models.SmartModel{
		ID:          uuid.New(),
		Name:        "Test Device",
		Description: "Test Description",
		Type:        models.DeviceType,
		Category:    models.WearableCategory,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("Update", mock.Anything, testModel).Return(testModel, nil)

	result, err := service.Update(context.Background(), testModel)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testModel.Name, result.Name)
	assert.Equal(t, testModel.Description, result.Description)

	mockRepo.AssertExpectations(t)
}

func TestSmartModelService_Update_Error(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	now := time.Now()
	testModel := &models.SmartModel{
		ID:          uuid.New(),
		Name:        "Test Device",
		Description: "Test Description",
		Type:        models.DeviceType,
		Category:    models.WearableCategory,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("Update", mock.Anything, testModel).Return(nil, assert.AnError)

	result, err := service.Update(context.Background(), testModel)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestSmartModelService_Delete(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	testID := uuid.New()

	mockRepo.On("Delete", mock.Anything, testID.String()).Return(nil)

	err := service.Delete(context.Background(), testID.String())

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSmartModelService_Delete_Error(t *testing.T) {
	mockRepo := new(mockSmartModelRepo)
	service := NewSmartModelService(mockRepo)

	testID := uuid.New()

	mockRepo.On("Delete", mock.Anything, testID.String()).Return(assert.AnError)

	err := service.Delete(context.Background(), testID.String())

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
