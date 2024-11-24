package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "smart-hub/gen/proto/smart_model/v1"
	"smart-hub/internal/application/service"
	"smart-hub/internal/common/logger"
	"smart-hub/internal/presentation/grpc/mapper"
)

type SmartModelHandler struct {
	pb.UnimplementedSmartModelServiceServer
	service *service.SmartModelService
	mapper  mapper.SmartModelMapper
}

func NewSmartModelHandler(
	service *service.SmartModelService,
	mapper mapper.SmartModelMapper,
) *SmartModelHandler {
	return &SmartModelHandler{
		service: service,
		mapper:  mapper,
	}
}

func (h *SmartModelHandler) CreateSmartModel(ctx context.Context, req *pb.CreateSmartModelRequest) (*pb.CreateSmartModelResponse, error) {
	logger.Debug("Creating smart model", "request", req)

	smartModel, err := h.mapper.ToDomain(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request: model is required")
	}

	createdModel, err := h.service.Create(ctx, smartModel)
	if err != nil {
		logger.Error("Failed to create smart model", "error", err)
		return nil, status.Error(codes.Internal, "failed to create smart model")
	}

	protoModel, err := h.mapper.ToProto(createdModel)
	if err != nil {
		logger.Error("Failed to convert smart model to proto", "error", err)
		return nil, status.Error(codes.Internal, "failed to convert smart model to proto")
	}

	return &pb.CreateSmartModelResponse{
		Model: protoModel,
	}, nil
}
