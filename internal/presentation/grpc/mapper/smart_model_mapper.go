package mapper

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "smart-hub/gen/proto/smart_model/v1"
	"smart-hub/internal/domain/models"
	"time"
)

type SmartModelMapper interface {
	ToProto(*models.SmartModel) (*pb.SmartModel, error)
	ToProtoList([]*models.SmartModel) ([]*pb.SmartModel, error)
	ToDomain(*pb.CreateSmartModelRequest) (*models.SmartModel, error)
	ToDomainUpdate(*pb.UpdateSmartModelRequest) (*models.SmartModel, error)
	ToCreateResponse(*models.SmartModel) (*pb.CreateSmartModelResponse, error)
	ToGetResponse(*models.SmartModel) (*pb.GetSmartModelResponse, error)
	ToListResponse([]*models.SmartModel) (*pb.ListSmartModelsResponse, error)
	ToUpdateResponse(*models.SmartModel) (*pb.UpdateSmartModelResponse, error)
}

type smartModelMapper struct{}

func NewSmartModelMapper() SmartModelMapper {
	return &smartModelMapper{}
}

func (m *smartModelMapper) ToProto(model *models.SmartModel) (*pb.SmartModel, error) {
	if model == nil {
		return nil, nil
	}

	metadata, err := structpb.NewStruct(model.Metadata)
	if err != nil {
		return nil, err
	}

	return &pb.SmartModel{
		Id:           model.ID.String(),
		Name:         model.Name,
		Description:  model.Description,
		Type:         mapDomainTypeToProto(model.Type),
		Category:     mapDomainCategoryToProto(model.Category),
		Manufacturer: model.Manufacturer,
		ModelNumber:  model.ModelNumber,
		Metadata:     metadata,
		CreatedAt:    timestamppb.New(model.CreatedAt),
		UpdatedAt:    timestamppb.New(model.UpdatedAt),
	}, nil
}

func (m *smartModelMapper) ToProtoList(models []*models.SmartModel) ([]*pb.SmartModel, error) {
	protoModels := make([]*pb.SmartModel, len(models))
	for i, model := range models {
		protoModel, err := m.ToProto(model)
		if err != nil {
			return nil, err
		}
		protoModels[i] = protoModel
	}

	return protoModels, nil
}

func (m *smartModelMapper) ToDomainUpdate(req *pb.UpdateSmartModelRequest) (*models.SmartModel, error) {
	if req == nil || req.Model == nil {
		return nil, nil
	}

	metadata := make(map[string]interface{})
	if req.Model.Metadata != nil {
		metadata = req.Model.Metadata.AsMap()
	}

	id, err := uuid.Parse(req.Model.Id)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &models.SmartModel{
		ID:           id,
		Name:         req.Model.Name,
		Description:  req.Model.Description,
		Type:         mapProtoTypeToDomain(req.Model.Type),
		Category:     mapProtoCategoryToDomain(req.Model.Category),
		Manufacturer: req.Model.Manufacturer,
		ModelNumber:  req.Model.ModelNumber,
		Metadata:     metadata,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (m *smartModelMapper) ToDomain(req *pb.CreateSmartModelRequest) (*models.SmartModel, error) {
	if req == nil || req.Model == nil {
		return nil, nil
	}

	metadata := make(map[string]interface{})
	if req.Model.Metadata != nil {
		metadata = req.Model.Metadata.AsMap()
	}

	id := uuid.New()

	now := time.Now()

	return &models.SmartModel{
		ID:           id,
		Name:         req.Model.Name,
		Description:  req.Model.Description,
		Type:         mapProtoTypeToDomain(req.Model.Type),
		Category:     mapProtoCategoryToDomain(req.Model.Category),
		Manufacturer: req.Model.Manufacturer,
		ModelNumber:  req.Model.ModelNumber,
		Metadata:     metadata,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (m *smartModelMapper) ToCreateResponse(model *models.SmartModel) (*pb.CreateSmartModelResponse, error) {
	protoModel, err := m.ToProto(model)
	if err != nil {
		return nil, err
	}

	return &pb.CreateSmartModelResponse{
		Model: protoModel,
	}, nil
}

func (m *smartModelMapper) ToGetResponse(model *models.SmartModel) (*pb.GetSmartModelResponse, error) {
	protoModel, err := m.ToProto(model)
	if err != nil {
		return nil, err
	}

	return &pb.GetSmartModelResponse{
		Model: protoModel,
	}, nil
}

func (m *smartModelMapper) ToListResponse(models []*models.SmartModel) (*pb.ListSmartModelsResponse, error) {
	protoModels := make([]*pb.SmartModel, len(models))
	for i, model := range models {
		protoModel, err := m.ToProto(model)
		if err != nil {
			return nil, err
		}
		protoModels[i] = protoModel
	}

	return &pb.ListSmartModelsResponse{
		Models: protoModels,
	}, nil
}

func (m *smartModelMapper) ToUpdateResponse(model *models.SmartModel) (*pb.UpdateSmartModelResponse, error) {
	protoModel, err := m.ToProto(model)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateSmartModelResponse{
		Model: protoModel,
	}, nil
}

func mapProtoTypeToDomain(t pb.ModelType) models.ModelType {
	switch t {
	case pb.ModelType_DEVICE:
		return models.DeviceType
	case pb.ModelType_SERVICE:
		return models.ServiceType
	default:
		return models.DeviceType
	}
}

func mapDomainTypeToProto(t models.ModelType) pb.ModelType {
	switch t {
	case models.DeviceType:
		return pb.ModelType_DEVICE
	case models.ServiceType:
		return pb.ModelType_SERVICE
	default:
		return pb.ModelType_DEVICE
	}
}

func mapProtoCategoryToDomain(c pb.ModelCategory) models.ModelCategory {
	switch c {
	case pb.ModelCategory_WEARABLE:
		return models.WearableCategory
	case pb.ModelCategory_CAMERA:
		return models.CameraCategory
	case pb.ModelCategory_WEATHER:
		return models.WeatherCategory
	case pb.ModelCategory_ENTERTAINMENT:
		return models.EntertainmentCategory
	default:
		return models.WearableCategory
	}
}

func mapDomainCategoryToProto(c models.ModelCategory) pb.ModelCategory {
	switch c {
	case models.WearableCategory:
		return pb.ModelCategory_WEARABLE
	case models.CameraCategory:
		return pb.ModelCategory_CAMERA
	case models.WeatherCategory:
		return pb.ModelCategory_WEATHER
	case models.EntertainmentCategory:
		return pb.ModelCategory_ENTERTAINMENT
	default:
		return pb.ModelCategory_WEARABLE
	}
}
