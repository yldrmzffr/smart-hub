package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "smart-hub/gen/proto/smart_model/v1"
	"smart-hub/internal/domain/models"
	"testing"
	"time"
)

type mockSmartModelService struct {
	mock.Mock
}

func (m *mockSmartModelService) Create(ctx context.Context, model *models.SmartModel) (*models.SmartModel, error) {
	args := m.Called(ctx, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartModel), args.Error(1)
}

func (m *mockSmartModelService) GetByID(ctx context.Context, id string) (*models.SmartModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartModel), args.Error(1)
}

func (m *mockSmartModelService) GetAll(ctx context.Context) ([]*models.SmartModel, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SmartModel), args.Error(1)
}

func (m *mockSmartModelService) Update(ctx context.Context, model *models.SmartModel) (*models.SmartModel, error) {
	args := m.Called(ctx, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartModel), args.Error(1)
}

func (m *mockSmartModelService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockSmartModelMapper struct {
	mock.Mock
}

func (m *mockSmartModelMapper) ToProto(model *models.SmartModel) (*pb.SmartModel, error) {
	args := m.Called(model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.SmartModel), args.Error(1)
}

func (m *mockSmartModelMapper) ToProtoList(models []*models.SmartModel) ([]*pb.SmartModel, error) {
	args := m.Called(models)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*pb.SmartModel), args.Error(1)
}

func (m *mockSmartModelMapper) ToDomain(req *pb.CreateSmartModelRequest) (*models.SmartModel, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartModel), args.Error(1)
}

func (m *mockSmartModelMapper) ToDomainUpdate(req *pb.UpdateSmartModelRequest) (*models.SmartModel, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartModel), args.Error(1)
}

func (m *mockSmartModelMapper) ToCreateResponse(model *models.SmartModel) (*pb.CreateSmartModelResponse, error) {
	args := m.Called(model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.CreateSmartModelResponse), args.Error(1)
}

func (m *mockSmartModelMapper) ToGetResponse(model *models.SmartModel) (*pb.GetSmartModelResponse, error) {
	args := m.Called(model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetSmartModelResponse), args.Error(1)
}

func (m *mockSmartModelMapper) ToListResponse(models []*models.SmartModel) (*pb.ListSmartModelsResponse, error) {
	args := m.Called(models)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.ListSmartModelsResponse), args.Error(1)
}

func (m *mockSmartModelMapper) ToUpdateResponse(model *models.SmartModel) (*pb.UpdateSmartModelResponse, error) {
	args := m.Called(model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.UpdateSmartModelResponse), args.Error(1)
}

func TestCreateSmartModel_Success(t *testing.T) {
	mockService := &mockSmartModelService{}
	mockMapper := &mockSmartModelMapper{}
	handler := NewSmartModelHandler(mockService, mockMapper)

	req := &pb.CreateSmartModelRequest{
		Model: &pb.CreateSmartModelInput{
			Name:         "Test Model",
			Description:  "Test Description",
			Type:         pb.ModelType_DEVICE,
			Category:     pb.ModelCategory_WEARABLE,
			Manufacturer: "Test Manufacturer",
			ModelNumber:  "TEST123",
		},
	}

	domainModel := &models.SmartModel{
		ID:           uuid.New(),
		Name:         req.Model.Name,
		Description:  req.Model.Description,
		Type:         models.DeviceType,
		Category:     models.WearableCategory,
		Manufacturer: req.Model.Manufacturer,
		ModelNumber:  req.Model.ModelNumber,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	protoModel := &pb.SmartModel{
		Id:           domainModel.ID.String(),
		Name:         domainModel.Name,
		Description:  domainModel.Description,
		Type:         pb.ModelType_DEVICE,
		Category:     pb.ModelCategory_WEARABLE,
		Manufacturer: domainModel.Manufacturer,
		ModelNumber:  domainModel.ModelNumber,
	}

	mockMapper.On("ToDomain", req).Return(domainModel, nil)
	mockService.On("Create", mock.Anything, domainModel).Return(domainModel, nil)
	mockMapper.On("ToProto", domainModel).Return(protoModel, nil)

	resp, err := handler.CreateSmartModel(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, protoModel, resp.Model)
	mockMapper.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetSmartModel_Success(t *testing.T) {
	mockService := &mockSmartModelService{}
	mockMapper := &mockSmartModelMapper{}
	handler := NewSmartModelHandler(mockService, mockMapper)

	modelID := uuid.New()
	req := &pb.GetSmartModelRequest{
		Id: modelID.String(),
	}

	domainModel := &models.SmartModel{
		ID:          modelID,
		Name:        "Test Model",
		Description: "Test Description",
		Type:        models.DeviceType,
		Category:    models.WearableCategory,
	}

	protoModel := &pb.SmartModel{
		Id:          modelID.String(),
		Name:        "Test Model",
		Description: "Test Description",
		Type:        pb.ModelType_DEVICE,
		Category:    pb.ModelCategory_WEARABLE,
	}

	mockService.On("GetByID", mock.Anything, modelID.String()).Return(domainModel, nil)
	mockMapper.On("ToProto", domainModel).Return(protoModel, nil)

	resp, err := handler.GetSmartModel(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, protoModel, resp.Model)
	mockService.AssertExpectations(t)
	mockMapper.AssertExpectations(t)
}

func TestListSmartModels_Success(t *testing.T) {
	mockService := &mockSmartModelService{}
	mockMapper := &mockSmartModelMapper{}
	handler := NewSmartModelHandler(mockService, mockMapper)

	req := &pb.ListSmartModelsRequest{}

	domainModels := []*models.SmartModel{
		{
			ID:   uuid.New(),
			Name: "Model 1",
			Type: models.DeviceType,
		},
		{
			ID:   uuid.New(),
			Name: "Model 2",
			Type: models.ServiceType,
		},
	}

	protoModels := []*pb.SmartModel{
		{
			Id:   domainModels[0].ID.String(),
			Name: "Model 1",
			Type: pb.ModelType_DEVICE,
		},
		{
			Id:   domainModels[1].ID.String(),
			Name: "Model 2",
			Type: pb.ModelType_SERVICE,
		},
	}

	mockService.On("GetAll", mock.Anything).Return(domainModels, nil)
	mockMapper.On("ToProtoList", domainModels).Return(protoModels, nil)

	resp, err := handler.ListSmartModels(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, protoModels, resp.Models)
	mockService.AssertExpectations(t)
	mockMapper.AssertExpectations(t)
}

func TestUpdateSmartModel_Success(t *testing.T) {
	mockService := &mockSmartModelService{}
	mockMapper := &mockSmartModelMapper{}
	handler := NewSmartModelHandler(mockService, mockMapper)

	modelID := uuid.New()
	req := &pb.UpdateSmartModelRequest{
		Model: &pb.UpdateSmartModelInput{
			Id:          modelID.String(),
			Name:        "Updated Model",
			Description: "Updated Description",
			Type:        pb.ModelType_DEVICE,
			Category:    pb.ModelCategory_WEARABLE,
		},
	}

	domainModel := &models.SmartModel{
		ID:          modelID,
		Name:        req.Model.Name,
		Description: req.Model.Description,
		Type:        models.DeviceType,
		Category:    models.WearableCategory,
	}

	protoModel := &pb.SmartModel{
		Id:          modelID.String(),
		Name:        req.Model.Name,
		Description: req.Model.Description,
		Type:        pb.ModelType_DEVICE,
		Category:    pb.ModelCategory_WEARABLE,
	}

	mockMapper.On("ToDomainUpdate", req).Return(domainModel, nil)
	mockService.On("Update", mock.Anything, domainModel).Return(domainModel, nil)
	mockMapper.On("ToProto", domainModel).Return(protoModel, nil)

	resp, err := handler.UpdateSmartModel(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, protoModel, resp.Model)
	mockMapper.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestDeleteSmartModel_Success(t *testing.T) {
	mockService := &mockSmartModelService{}
	mockMapper := &mockSmartModelMapper{}
	handler := NewSmartModelHandler(mockService, mockMapper)

	modelID := uuid.New()
	req := &pb.DeleteSmartModelRequest{
		Id: modelID.String(),
	}

	mockService.On("Delete", mock.Anything, modelID.String()).Return(nil)

	resp, err := handler.DeleteSmartModel(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockService.AssertExpectations(t)
}

func TestCreateSmartModel_ValidationError(t *testing.T) {
	mockService := new(mockSmartModelService)
	mockMapper := new(mockSmartModelMapper)
	handler := NewSmartModelHandler(mockService, mockMapper)

	req := &pb.CreateSmartModelRequest{
		Model: &pb.CreateSmartModelInput{
			Name: "Test", // Missing required fields
		},
	}

	domainModel := &models.SmartModel{
		Name: req.Model.Name,
	}

	mockMapper.On("ToDomain", req).Return(domainModel, nil)

	resp, err := handler.CreateSmartModel(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestCreateSmartModel_ServiceError(t *testing.T) {
	mockService := new(mockSmartModelService)
	mockMapper := new(mockSmartModelMapper)
	handler := NewSmartModelHandler(mockService, mockMapper)

	req := &pb.CreateSmartModelRequest{
		Model: &pb.CreateSmartModelInput{
			Name:         "Test Model",
			Description:  "Test Description",
			Type:         pb.ModelType_DEVICE,
			Category:     pb.ModelCategory_WEARABLE,
			Manufacturer: "Test Manufacturer",
			ModelNumber:  "TEST123",
		},
	}

	domainModel := &models.SmartModel{
		Name:         req.Model.Name,
		Description:  req.Model.Description,
		Type:         models.DeviceType,
		Category:     models.WearableCategory,
		Manufacturer: req.Model.Manufacturer,
		ModelNumber:  req.Model.ModelNumber,
	}

	mockMapper.On("ToDomain", req).Return(domainModel, nil)
	mockService.On("Create", mock.Anything, domainModel).Return(nil, assert.AnError)

	resp, err := handler.CreateSmartModel(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGetSmartModel_InvalidID(t *testing.T) {
	mockService := new(mockSmartModelService)
	mockMapper := new(mockSmartModelMapper)
	handler := NewSmartModelHandler(mockService, mockMapper)

	req := &pb.GetSmartModelRequest{
		Id: "invalid-uuid",
	}

	resp, err := handler.GetSmartModel(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestGetSmartModel_ServiceError(t *testing.T) {
	mockService := new(mockSmartModelService)
	mockMapper := new(mockSmartModelMapper)
	handler := NewSmartModelHandler(mockService, mockMapper)

	modelID := uuid.New()
	req := &pb.GetSmartModelRequest{
		Id: modelID.String(),
	}

	mockService.On("GetByID", mock.Anything, modelID.String()).Return(nil, assert.AnError)

	resp, err := handler.GetSmartModel(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestListSmartModels_ServiceError(t *testing.T) {
	mockService := new(mockSmartModelService)
	mockMapper := new(mockSmartModelMapper)
	handler := NewSmartModelHandler(mockService, mockMapper)

	mockService.On("GetAll", mock.Anything).Return(nil, assert.AnError)

	resp, err := handler.ListSmartModels(context.Background(), &pb.ListSmartModelsRequest{})

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestUpdateSmartModel_ValidationError(t *testing.T) {
	mockService := new(mockSmartModelService)
	mockMapper := new(mockSmartModelMapper)
	handler := NewSmartModelHandler(mockService, mockMapper)

	req := &pb.UpdateSmartModelRequest{
		Model: &pb.UpdateSmartModelInput{
			Id: "invalid-uuid",
		},
	}

	mockMapper.On("ToDomainUpdate", req).Return(&models.SmartModel{}, nil)

	resp, err := handler.UpdateSmartModel(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUpdateSmartModel_ServiceError(t *testing.T) {
	mockService := new(mockSmartModelService)
	mockMapper := new(mockSmartModelMapper)
	handler := NewSmartModelHandler(mockService, mockMapper)

	modelID := uuid.New()
	req := &pb.UpdateSmartModelRequest{
		Model: &pb.UpdateSmartModelInput{
			Id:          modelID.String(),
			Name:        "Updated Model",
			Description: "Updated Description",
			Type:        pb.ModelType_DEVICE,
			Category:    pb.ModelCategory_WEARABLE,
		},
	}

	domainModel := &models.SmartModel{
		ID:          modelID,
		Name:        req.Model.Name,
		Description: req.Model.Description,
		Type:        models.DeviceType,
		Category:    models.WearableCategory,
	}

	mockMapper.On("ToDomainUpdate", req).Return(domainModel, nil)
	mockService.On("Update", mock.Anything, domainModel).Return(nil, assert.AnError)

	resp, err := handler.UpdateSmartModel(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestDeleteSmartModel_InvalidID(t *testing.T) {
	mockService := new(mockSmartModelService)
	mockMapper := new(mockSmartModelMapper)
	handler := NewSmartModelHandler(mockService, mockMapper)

	req := &pb.DeleteSmartModelRequest{
		Id: "invalid-uuid",
	}

	resp, err := handler.DeleteSmartModel(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestDeleteSmartModel_ServiceError(t *testing.T) {
	mockService := new(mockSmartModelService)
	mockMapper := new(mockSmartModelMapper)
	handler := NewSmartModelHandler(mockService, mockMapper)

	modelID := uuid.New()
	req := &pb.DeleteSmartModelRequest{
		Id: modelID.String(),
	}

	mockService.On("Delete", mock.Anything, modelID.String()).Return(assert.AnError)

	resp, err := handler.DeleteSmartModel(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}
