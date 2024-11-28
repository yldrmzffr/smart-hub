package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "smart-hub/gen/proto/smart_feature/v1"
	"smart-hub/internal/application/interfaces"
	"smart-hub/internal/common/logger"
	"smart-hub/internal/common/validation"
	"smart-hub/internal/presentation/grpc/mapper"
)

type SmartFeatureHandler struct {
	pb.UnimplementedSmartFeatureServiceServer
	service interfaces.SmartFeatureService
	mapper  mapper.SmartFeatureMapper
}

func NewSmartFeatureHandler(
	service interfaces.SmartFeatureService,
	mapper mapper.SmartFeatureMapper,
) *SmartFeatureHandler {
	return &SmartFeatureHandler{
		service: service,
		mapper:  mapper,
	}
}

func (h *SmartFeatureHandler) CreateSmartFeature(ctx context.Context, req *pb.CreateSmartFeatureRequest) (*pb.CreateSmartFeatureResponse, error) {
	logger.Debug("Creating smart feature", "request", req)

	smartFeature, err := h.mapper.ToDomain(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request: feature is required")
	}

	if err := validation.ValidateUUID(smartFeature.ModelID.String()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validation.ValidateStruct(smartFeature); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	createdFeature, err := h.service.Create(ctx, smartFeature)
	if err != nil {
		logger.Error("Failed to create smart feature", "error", err)
		return nil, status.Error(codes.Internal, "failed to create smart feature")
	}

	protoFeature, err := h.mapper.ToProto(createdFeature)
	if err != nil {
		logger.Error("Failed to convert smart feature to proto", "error", err)
		return nil, status.Error(codes.Internal, "failed to convert smart feature to proto")
	}

	return &pb.CreateSmartFeatureResponse{
		Feature: protoFeature,
	}, nil
}

func (h *SmartFeatureHandler) GetSmartFeature(ctx context.Context, req *pb.GetSmartFeatureRequest) (*pb.GetSmartFeatureResponse, error) {
	logger.Debug("Getting smart feature", "request", req)

	if err := validation.ValidateUUID(req.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	smartFeature, err := h.service.GetByID(ctx, req.Id)
	if err != nil {
		logger.Error("Failed to get smart feature", "error", err)
		return nil, status.Error(codes.Internal, "failed to get smart feature")
	}

	protoFeature, err := h.mapper.ToProto(smartFeature)
	if err != nil {
		logger.Error("Failed to convert smart feature to proto", "error", err)
		return nil, status.Error(codes.Internal, "failed to convert smart feature to proto")
	}

	return &pb.GetSmartFeatureResponse{
		Feature: protoFeature,
	}, nil
}

func (h *SmartFeatureHandler) GetFeaturesByModelID(ctx context.Context, req *pb.GetFeaturesByModelIDRequest) (*pb.GetFeaturesByModelIDResponse, error) {
	logger.Debug("Getting smart features by model ID", "request", req)

	err := validation.ValidateUUID(req.ModelId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	smartFeatures, err := h.service.GetWithModelID(ctx, req.ModelId)
	if err != nil {
		logger.Error("Failed to get smart features by model ID", "error", err)
		return nil, status.Error(codes.Internal, "failed to get smart features by model ID")
	}

	protoFeatures, err := h.mapper.ToListResponse(smartFeatures)
	if err != nil {
		logger.Error("Failed to convert smart features to proto", "error", err)
		return nil, status.Error(codes.Internal, "failed to convert smart features to proto")
	}

	return &pb.GetFeaturesByModelIDResponse{
		Features: protoFeatures.Features,
	}, nil
}

func (h *SmartFeatureHandler) UpdateSmartFeature(ctx context.Context, req *pb.UpdateSmartFeatureRequest) (*pb.UpdateSmartFeatureResponse, error) {
	logger.Debug("Updating smart feature", "request", req)

	smartFeature, err := h.mapper.ToDomainUpdate(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request: feature is required")
	}

	if err := validation.ValidateStruct(smartFeature); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	updatedFeature, err := h.service.Update(ctx, smartFeature)
	if err != nil {
		logger.Error("Failed to update smart feature", "error", err)
		return nil, status.Error(codes.Internal, "failed to update smart feature")
	}

	protoFeature, err := h.mapper.ToProto(updatedFeature)
	if err != nil {
		logger.Error("Failed to convert smart feature to proto", "error", err)
		return nil, status.Error(codes.Internal, "failed to convert smart feature to proto")
	}

	return &pb.UpdateSmartFeatureResponse{
		Feature: protoFeature,
	}, nil
}

func (h *SmartFeatureHandler) DeleteSmartFeature(ctx context.Context, req *pb.DeleteSmartFeatureRequest) (*pb.DeleteSmartFeatureResponse, error) {
	logger.Debug("Deleting smart feature", "request", req)

	err := validation.ValidateUUID(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = h.service.Delete(ctx, req.Id)
	if err != nil {
		logger.Error("Failed to delete smart feature", "error", err)
		return nil, status.Error(codes.Internal, "failed to delete smart feature")
	}

	return &pb.DeleteSmartFeatureResponse{}, nil
}
