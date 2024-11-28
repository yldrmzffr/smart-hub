package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	_ "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	pb "smart-hub/gen/proto/smart_feature/v1"
	_ "smart-hub/internal/application/service"
	"smart-hub/internal/domain/models"
	"testing"
	"time"
)

type mockSmartFeatureService struct {
	mock.Mock
}

func (m *mockSmartFeatureService) Create(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error) {
	args := m.Called(ctx, feature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureService) GetByID(ctx context.Context, id string) (*models.SmartFeature, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureService) GetWithModelID(ctx context.Context, modelID string) ([]*models.SmartFeature, error) {
	args := m.Called(ctx, modelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureService) GetAll(ctx context.Context) ([]*models.SmartFeature, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureService) Update(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error) {
	args := m.Called(ctx, feature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockSmartFeatureMapper struct {
	mock.Mock
}

func (m *mockSmartFeatureMapper) ToProto(feature *models.SmartFeature) (*pb.SmartFeature, error) {
	args := m.Called(feature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureMapper) ToDomain(req *pb.CreateSmartFeatureRequest) (*models.SmartFeature, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureMapper) ToDomainUpdate(req *pb.UpdateSmartFeatureRequest) (*models.SmartFeature, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SmartFeature), args.Error(1)
}

func (m *mockSmartFeatureMapper) ToCreateResponse(feature *models.SmartFeature) (*pb.CreateSmartFeatureResponse, error) {
	args := m.Called(feature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.CreateSmartFeatureResponse), args.Error(1)
}

func (m *mockSmartFeatureMapper) ToGetResponse(feature *models.SmartFeature) (*pb.GetSmartFeatureResponse, error) {
	args := m.Called(feature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetSmartFeatureResponse), args.Error(1)
}

func (m *mockSmartFeatureMapper) ToListResponse(features []*models.SmartFeature) (*pb.GetFeaturesByModelIDResponse, error) {
	args := m.Called(features)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetFeaturesByModelIDResponse), args.Error(1)
}

func (m *mockSmartFeatureMapper) ToUpdateResponse(feature *models.SmartFeature) (*pb.UpdateSmartFeatureResponse, error) {
	args := m.Called(feature)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.UpdateSmartFeatureResponse), args.Error(1)
}

func TestCreateSmartFeature_Success(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	params, _ := structpb.NewStruct(map[string]interface{}{
		"key": "value",
	})

	req := &pb.CreateSmartFeatureRequest{
		Feature: &pb.CreateSmartFeatureInput{
			ModelId:       uuid.New().String(),
			Name:          "Test Feature",
			Description:   "Test Description",
			Protocol:      pb.ProtocolType_REST,
			InterfacePath: "/test",
			Parameters:    params,
		},
	}

	domainFeature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.MustParse(req.Feature.ModelId),
		Name:          req.Feature.Name,
		Description:   req.Feature.Description,
		Protocol:      models.RestProtocol,
		InterfacePath: "/test",
		Parameters:    params.AsMap(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	protoFeature := &pb.SmartFeature{
		Id:            domainFeature.ID.String(),
		ModelId:       domainFeature.ModelID.String(),
		Name:          domainFeature.Name,
		Description:   domainFeature.Description,
		Protocol:      pb.ProtocolType_REST,
		InterfacePath: domainFeature.InterfacePath,
		Parameters:    params,
	}

	mockMapper.On("ToDomain", req).Return(domainFeature, nil)
	mockService.On("Create", mock.Anything, domainFeature).Return(domainFeature, nil)
	mockMapper.On("ToProto", domainFeature).Return(protoFeature, nil)

	resp, err := handler.CreateSmartFeature(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, protoFeature, resp.Feature)
	mockMapper.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetSmartFeature_Success(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	featureID := uuid.New()
	req := &pb.GetSmartFeatureRequest{
		Id: featureID.String(),
	}

	domainFeature := &models.SmartFeature{
		ID:            featureID,
		Name:          "Test Feature",
		Description:   "Test Description",
		Protocol:      models.RestProtocol,
		InterfacePath: "/test",
	}

	protoFeature := &pb.SmartFeature{
		Id:            featureID.String(),
		Name:          "Test Feature",
		Description:   "Test Description",
		Protocol:      pb.ProtocolType_REST,
		InterfacePath: "/test",
	}

	mockService.On("GetByID", mock.Anything, featureID.String()).Return(domainFeature, nil)
	mockMapper.On("ToProto", domainFeature).Return(protoFeature, nil)

	resp, err := handler.GetSmartFeature(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, protoFeature, resp.Feature)
	mockService.AssertExpectations(t)
	mockMapper.AssertExpectations(t)
}

func TestGetFeaturesByModelID_Success(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	modelID := uuid.New().String()
	req := &pb.GetFeaturesByModelIDRequest{
		ModelId: modelID,
	}

	domainFeatures := []*models.SmartFeature{
		{
			ID:      uuid.New(),
			ModelID: uuid.MustParse(modelID),
			Name:    "Feature 1",
		},
		{
			ID:      uuid.New(),
			ModelID: uuid.MustParse(modelID),
			Name:    "Feature 2",
		},
	}

	protoResponse := &pb.GetFeaturesByModelIDResponse{
		Features: []*pb.SmartFeature{
			{
				Id:      domainFeatures[0].ID.String(),
				ModelId: modelID,
				Name:    "Feature 1",
			},
			{
				Id:      domainFeatures[1].ID.String(),
				ModelId: modelID,
				Name:    "Feature 2",
			},
		},
	}

	mockService.On("GetWithModelID", mock.Anything, modelID).Return(domainFeatures, nil)
	mockMapper.On("ToListResponse", domainFeatures).Return(protoResponse, nil)

	resp, err := handler.GetFeaturesByModelID(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, protoResponse.Features, resp.Features)
	mockService.AssertExpectations(t)
	mockMapper.AssertExpectations(t)
}

func TestUpdateSmartFeature_Success(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	featureID := uuid.New()
	req := &pb.UpdateSmartFeatureRequest{
		Feature: &pb.UpdateSmartFeatureInput{
			Id:            featureID.String(),
			Name:          "Updated Feature",
			Description:   "Updated Description",
			Protocol:      pb.ProtocolType_REST,
			InterfacePath: "/updated",
		},
	}

	domainFeature := &models.SmartFeature{
		ID:            featureID,
		Name:          req.Feature.Name,
		Description:   req.Feature.Description,
		Protocol:      models.RestProtocol,
		InterfacePath: req.Feature.InterfacePath,
	}

	protoFeature := &pb.SmartFeature{
		Id:            featureID.String(),
		Name:          req.Feature.Name,
		Description:   req.Feature.Description,
		Protocol:      pb.ProtocolType_REST,
		InterfacePath: req.Feature.InterfacePath,
	}

	mockMapper.On("ToDomainUpdate", req).Return(domainFeature, nil)
	mockService.On("Update", mock.Anything, domainFeature).Return(domainFeature, nil)
	mockMapper.On("ToProto", domainFeature).Return(protoFeature, nil)

	resp, err := handler.UpdateSmartFeature(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, protoFeature, resp.Feature)
	mockMapper.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestDeleteSmartFeature_Success(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	featureID := uuid.New()
	req := &pb.DeleteSmartFeatureRequest{
		Id: featureID.String(),
	}

	mockService.On("Delete", mock.Anything, featureID.String()).Return(nil)

	resp, err := handler.DeleteSmartFeature(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockService.AssertExpectations(t)
}

func TestCreateSmartFeature_ValidationError(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	req := &pb.CreateSmartFeatureRequest{
		Feature: &pb.CreateSmartFeatureInput{
			Name: "Test", // Missing required fields
		},
	}

	mockMapper.On("ToDomain", req).Return(&models.SmartFeature{}, nil)

	resp, err := handler.CreateSmartFeature(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestCreateSmartFeature_ServiceError(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	req := &pb.CreateSmartFeatureRequest{
		Feature: &pb.CreateSmartFeatureInput{
			ModelId:       uuid.New().String(),
			Name:          "Test Feature",
			Description:   "Test Description", // Added required field
			Protocol:      pb.ProtocolType_REST,
			InterfacePath: "/test", // Added required field
		},
	}

	domainFeature := &models.SmartFeature{
		ModelID:       uuid.MustParse(req.Feature.ModelId),
		Name:          req.Feature.Name,
		Description:   req.Feature.Description,
		Protocol:      models.RestProtocol,
		InterfacePath: req.Feature.InterfacePath,
	}

	mockMapper.On("ToDomain", req).Return(domainFeature, nil)
	mockService.On("Create", mock.Anything, domainFeature).Return(nil, assert.AnError)

	resp, err := handler.CreateSmartFeature(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	mockMapper.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestGetSmartFeature_InvalidID(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	req := &pb.GetSmartFeatureRequest{
		Id: "invalid-uuid",
	}

	resp, err := handler.GetSmartFeature(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestGetSmartFeature_ServiceError(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	featureID := uuid.New()
	req := &pb.GetSmartFeatureRequest{
		Id: featureID.String(),
	}

	mockService.On("GetByID", mock.Anything, featureID.String()).Return(nil, assert.AnError)

	resp, err := handler.GetSmartFeature(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGetFeaturesByModelID_InvalidModelID(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	req := &pb.GetFeaturesByModelIDRequest{
		ModelId: "invalid-uuid",
	}

	resp, err := handler.GetFeaturesByModelID(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestGetFeaturesByModelID_ServiceError(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	modelID := uuid.New().String()
	req := &pb.GetFeaturesByModelIDRequest{
		ModelId: modelID,
	}

	mockService.On("GetWithModelID", mock.Anything, modelID).Return(nil, assert.AnError)

	resp, err := handler.GetFeaturesByModelID(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestUpdateSmartFeature_ValidationError(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	req := &pb.UpdateSmartFeatureRequest{
		Feature: &pb.UpdateSmartFeatureInput{
			Id: "invalid-uuid",
		},
	}

	mockMapper.On("ToDomainUpdate", req).Return(&models.SmartFeature{}, nil)

	resp, err := handler.UpdateSmartFeature(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUpdateSmartFeature_ServiceError(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	featureID := uuid.New()
	req := &pb.UpdateSmartFeatureRequest{
		Feature: &pb.UpdateSmartFeatureInput{
			Id:            featureID.String(),
			Name:          "Updated Feature",
			Description:   "Updated Description",
			Protocol:      pb.ProtocolType_REST,
			InterfacePath: "/test",
		},
	}

	domainFeature := &models.SmartFeature{
		ID:            featureID,
		Name:          req.Feature.Name,
		Description:   req.Feature.Description,
		Protocol:      models.RestProtocol,
		InterfacePath: "/test",
	}

	mockMapper.On("ToDomainUpdate", req).Return(domainFeature, nil)
	mockService.On("Update", mock.Anything, domainFeature).Return(nil, assert.AnError)

	resp, err := handler.UpdateSmartFeature(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	mockMapper.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

func TestDeleteSmartFeature_InvalidID(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	req := &pb.DeleteSmartFeatureRequest{
		Id: "invalid-uuid",
	}

	resp, err := handler.DeleteSmartFeature(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestDeleteSmartFeature_ServiceError(t *testing.T) {
	mockService := &mockSmartFeatureService{}
	mockMapper := &mockSmartFeatureMapper{}
	handler := NewSmartFeatureHandler(mockService, mockMapper)

	featureID := uuid.New()
	req := &pb.DeleteSmartFeatureRequest{
		Id: featureID.String(),
	}

	mockService.On("Delete", mock.Anything, featureID.String()).Return(assert.AnError)

	resp, err := handler.DeleteSmartFeature(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}
