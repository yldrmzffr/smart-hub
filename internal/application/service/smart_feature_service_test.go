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

type mockSmartFeatureRepo struct {
	mock.Mock
}

func (m *mockSmartFeatureRepo) Create(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error) {
	args := m.Called(ctx, feature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureRepo) GetByID(ctx context.Context, id string) (*models.SmartFeature, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureRepo) GetWithModelID(ctx context.Context, modelID string) ([]*models.SmartFeature, error) {
	args := m.Called(ctx, modelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureRepo) GetAll(ctx context.Context) ([]*models.SmartFeature, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureRepo) Update(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error) {
	args := m.Called(ctx, feature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSmartFeatureService_Create(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	now := time.Now()

	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Feature Name",
		Description:   "Feature Description",
		Protocol:      "rest",
		InterfacePath: "/test",
		Parameters:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	mockRepo.On("Create", mock.Anything, feature).Return(feature, nil)

	createdFeature, err := service.Create(context.Background(), feature)

	assert.NoError(t, err)
	assert.NotNil(t, createdFeature)
	assert.Equal(t, feature.Name, createdFeature.Name)
	assert.Equal(t, feature.Description, createdFeature.Description)

	mockRepo.AssertExpectations(t)
}

func TestSmartFeatureService_Create_Error(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	now := time.Now()

	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Feature Name",
		Description:   "Feature Description",
		Protocol:      "rest",
		InterfacePath: "/test",
		Parameters:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	mockRepo.On("Create", mock.Anything, feature).Return(nil, assert.AnError)

	createdFeature, err := service.Create(context.Background(), feature)

	assert.Error(t, err)
	assert.Nil(t, createdFeature)

	mockRepo.AssertExpectations(t)
}

func TestSmartFeatureService_GetByID(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	now := time.Now()

	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Feature Name",
		Description:   "Feature Description",
		Protocol:      "rest",
		InterfacePath: "/test",
		Parameters:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	mockRepo.On("GetByID", mock.Anything, feature.ID.String()).Return(feature, nil)

	result, err := service.GetByID(context.Background(), feature.ID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, feature.Name, result.Name)
	assert.Equal(t, feature.Description, result.Description)

	mockRepo.AssertExpectations(t)
}

func TestSmartFeatureService_GetByID_Error(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	now := time.Now()

	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Feature Name",
		Description:   "Feature Description",
		Protocol:      "rest",
		InterfacePath: "/test",
		Parameters:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	mockRepo.On("GetByID", mock.Anything, feature.ID.String()).Return(nil, assert.AnError)

	result, err := service.GetByID(context.Background(), feature.ID.String())

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestSmartFeatureService_GetWithModelID(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	now := time.Now()

	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Feature Name",
		Description:   "Feature Description",
		Protocol:      "rest",
		InterfacePath: "/test",
		Parameters:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	mockRepo.On("GetWithModelID", mock.Anything, feature.ModelID.String()).Return([]*models.SmartFeature{feature}, nil)

	result, err := service.GetWithModelID(context.Background(), feature.ModelID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, feature.Name, result[0].Name)
	assert.Equal(t, feature.Description, result[0].Description)

	mockRepo.AssertExpectations(t)
}

func TestSmartFeatureService_GetWithModelID_Error(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	now := time.Now()

	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Feature Name",
		Description:   "Feature Description",
		Protocol:      "rest",
		InterfacePath: "/test",
		Parameters:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	mockRepo.On("GetWithModelID", mock.Anything, feature.ModelID.String()).Return(nil, assert.AnError)

	result, err := service.GetWithModelID(context.Background(), feature.ModelID.String())

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestSmartFeatureService_GetAll(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	now := time.Now()

	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Feature Name",
		Description:   "Feature Description",
		Protocol:      "rest",
		InterfacePath: "/test",
		Parameters:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	mockRepo.On("GetAll", mock.Anything).Return([]*models.SmartFeature{feature}, nil)

	result, err := service.GetAll(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, feature.Name, result[0].Name)
	assert.Equal(t, feature.Description, result[0].Description)

	mockRepo.AssertExpectations(t)
}

func TestSmartFeatureService_GetAll_Error(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	mockRepo.On("GetAll", mock.Anything).Return(nil, assert.AnError)

	result, err := service.GetAll(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestSmartFeatureService_Update(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	now := time.Now()

	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Feature Name",
		Description:   "Feature Description",
		Protocol:      "rest",
		InterfacePath: "/test",
		Parameters:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	mockRepo.On("Update", mock.Anything, feature).Return(feature, nil)

	updatedFeature, err := service.Update(context.Background(), feature)

	assert.NoError(t, err)
	assert.NotNil(t, updatedFeature)
	assert.Equal(t, feature.Name, updatedFeature.Name)
	assert.Equal(t, feature.Description, updatedFeature.Description)

	mockRepo.AssertExpectations(t)
}

func TestSmartFeatureService_Update_Error(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	now := time.Now()

	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Feature Name",
		Description:   "Feature Description",
		Protocol:      "rest",
		InterfacePath: "/test",
		Parameters:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	mockRepo.On("Update", mock.Anything, feature).Return(nil, assert.AnError)

	updatedFeature, err := service.Update(context.Background(), feature)

	assert.Error(t, err)
	assert.Nil(t, updatedFeature)

	mockRepo.AssertExpectations(t)
}

func TestSmartFeatureService_Delete(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "test-id").Return(nil)

	err := service.Delete(context.Background(), "test-id")

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestSmartFeatureService_Delete_Error(t *testing.T) {
	mockRepo := new(mockSmartFeatureRepo)
	service := NewSmartFeatureService(mockRepo)

	mockRepo.On("Delete", mock.Anything, "test-id").Return(assert.AnError)

	err := service.Delete(context.Background(), "test-id")

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
