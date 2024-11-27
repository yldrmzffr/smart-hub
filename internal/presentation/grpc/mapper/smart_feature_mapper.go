package mapper

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "smart-hub/gen/proto/smart_feature/v1"
	"smart-hub/internal/domain/models"
	"time"
)

type SmartFeatureMapper interface {
	ToProto(*models.SmartFeature) (*pb.SmartFeature, error)
	ToDomain(*pb.CreateSmartFeatureRequest) (*models.SmartFeature, error)
	ToDomainUpdate(*pb.UpdateSmartFeatureRequest) (*models.SmartFeature, error)
	ToCreateResponse(*models.SmartFeature) (*pb.CreateSmartFeatureResponse, error)
	ToGetResponse(*models.SmartFeature) (*pb.GetSmartFeatureResponse, error)
	ToListResponse([]*models.SmartFeature) (*pb.GetFeaturesByModelIDResponse, error)
	ToUpdateResponse(*models.SmartFeature) (*pb.UpdateSmartFeatureResponse, error)
}

type smartFeatureMapper struct{}

func NewSmartFeatureMapper() SmartFeatureMapper {
	return &smartFeatureMapper{}
}

func (m *smartFeatureMapper) ToProto(model *models.SmartFeature) (*pb.SmartFeature, error) {
	if model == nil {
		return nil, nil
	}

	parameters, err := structpb.NewStruct(model.Parameters)
	if err != nil {
		return nil, err
	}

	return &pb.SmartFeature{
		Id:            model.ID.String(),
		ModelId:       model.ModelID.String(),
		Name:          model.Name,
		Description:   model.Description,
		Protocol:      mapDomainProtocolToProto(model.Protocol),
		InterfacePath: model.InterfacePath,
		Parameters:    parameters,
		CreatedAt:     timestamppb.New(model.CreatedAt),
		UpdatedAt:     timestamppb.New(model.UpdatedAt),
	}, nil
}

func (m *smartFeatureMapper) ToDomain(req *pb.CreateSmartFeatureRequest) (*models.SmartFeature, error) {

	parameters := make(map[string]interface{})
	if req.Feature.Parameters != nil {
		parameters = req.Feature.Parameters.AsMap()
	}

	modelId, err := uuid.Parse(req.Feature.ModelId)
	if err != nil {
		return nil, err
	}

	id := uuid.New()
	now := time.Now()

	return &models.SmartFeature{
		ID:            id,
		Name:          req.Feature.Name,
		ModelID:       modelId,
		Description:   req.Feature.Description,
		Protocol:      mapProtoProtocolToDomain(req.Feature.Protocol),
		InterfacePath: req.Feature.InterfacePath,
		Parameters:    parameters,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (m *smartFeatureMapper) ToDomainUpdate(req *pb.UpdateSmartFeatureRequest) (*models.SmartFeature, error) {
	if req == nil || req.Feature == nil {
		return nil, nil
	}

	parameters := make(map[string]interface{})
	if req.Feature.Parameters != nil {
		parameters = req.Feature.Parameters.AsMap()
	}

	id, err := uuid.Parse(req.Feature.Id)
	if err != nil {
		return nil, err
	}

	return &models.SmartFeature{
		ID:            id,
		Name:          req.Feature.Name,
		Description:   req.Feature.Description,
		Protocol:      mapProtoProtocolToDomain(req.Feature.Protocol),
		InterfacePath: req.Feature.InterfacePath,
		Parameters:    parameters,
	}, nil
}

func (m *smartFeatureMapper) ToCreateResponse(model *models.SmartFeature) (*pb.CreateSmartFeatureResponse, error) {
	protoModel, err := m.ToProto(model)
	if err != nil {
		return nil, err
	}

	return &pb.CreateSmartFeatureResponse{
		Feature: protoModel,
	}, nil
}

func (m *smartFeatureMapper) ToGetResponse(model *models.SmartFeature) (*pb.GetSmartFeatureResponse, error) {
	protoModel, err := m.ToProto(model)
	if err != nil {
		return nil, err
	}

	return &pb.GetSmartFeatureResponse{
		Feature: protoModel,
	}, nil
}

func (m *smartFeatureMapper) ToListResponse(models []*models.SmartFeature) (*pb.GetFeaturesByModelIDResponse, error) {
	protoModels := make([]*pb.SmartFeature, 0, len(models))
	for _, model := range models {
		protoModel, err := m.ToProto(model)
		if err != nil {
			return nil, err
		}
		protoModels = append(protoModels, protoModel)
	}

	return &pb.GetFeaturesByModelIDResponse{
		Features: protoModels,
	}, nil
}

func (m *smartFeatureMapper) ToUpdateResponse(model *models.SmartFeature) (*pb.UpdateSmartFeatureResponse, error) {
	protoModel, err := m.ToProto(model)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateSmartFeatureResponse{
		Feature: protoModel,
	}, nil
}

func mapProtoProtocolToDomain(p pb.ProtocolType) models.ProtocolType {
	switch p {
	case pb.ProtocolType_REST:
		return models.RestProtocol
	case pb.ProtocolType_GRPC:
		return models.GrpcProtocol
	case pb.ProtocolType_MQTT:
		return models.MqttProtocol
	case pb.ProtocolType_WEBSOCKET:
		return models.WebsocketProtocol
	default:
		return models.RestProtocol
	}
}

func mapDomainProtocolToProto(p models.ProtocolType) pb.ProtocolType {
	switch p {
	case models.RestProtocol:
		return pb.ProtocolType_REST
	case models.GrpcProtocol:
		return pb.ProtocolType_GRPC
	case models.MqttProtocol:
		return pb.ProtocolType_MQTT
	case models.WebsocketProtocol:
		return pb.ProtocolType_WEBSOCKET
	default:
		return pb.ProtocolType_REST
	}
}
