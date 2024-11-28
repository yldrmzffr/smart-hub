package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "smart-hub/gen/proto/smart_model/v1"
	"smart-hub/internal/application/interfaces"
	"smart-hub/internal/common/logger"
	"smart-hub/internal/common/validation"
	"smart-hub/internal/presentation/grpc/mapper"
)

type SmartModelHandler struct {
	pb.UnimplementedSmartModelServiceServer
	service interfaces.SmartModelService
	mapper  mapper.SmartModelMapper
}

func NewSmartModelHandler(
	service interfaces.SmartModelService,
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

	if err := validation.ValidateStruct(smartModel); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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

func (h *SmartModelHandler) GetSmartModel(ctx context.Context, req *pb.GetSmartModelRequest) (*pb.GetSmartModelResponse, error) {
	logger.Debug("Getting smart model", "request", req)

	err := validation.ValidateUUID(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	smartModel, err := h.service.GetByID(ctx, req.Id)
	if err != nil {
		logger.Error("Failed to get smart model", "error", err)
		return nil, status.Error(codes.Internal, "failed to get smart model")
	}

	protoModel, err := h.mapper.ToProto(smartModel)
	if err != nil {
		logger.Error("Failed to convert smart model to proto", "error", err)
		return nil, status.Error(codes.Internal, "failed to convert smart model to proto")
	}

	return &pb.GetSmartModelResponse{
		Model: protoModel,
	}, nil
}

func (h *SmartModelHandler) ListSmartModels(ctx context.Context, req *pb.ListSmartModelsRequest) (*pb.ListSmartModelsResponse, error) {
	logger.Debug("Listing smart models", "request", req)

	models, err := h.service.GetAll(ctx)
	if err != nil {
		logger.Error("Failed to list smart models", "error", err)
		return nil, status.Error(codes.Internal, "failed to list smart models")
	}

	protoModels, err := h.mapper.ToProtoList(models)
	if err != nil {
		logger.Error("Failed to convert smart models to proto", "error", err)
		return nil, status.Error(codes.Internal, "failed to convert smart models to proto")
	}

	return &pb.ListSmartModelsResponse{
		Models: protoModels,
	}, nil
}

func (h *SmartModelHandler) UpdateSmartModel(ctx context.Context, req *pb.UpdateSmartModelRequest) (*pb.UpdateSmartModelResponse, error) {
	logger.Debug("Updating smart model", "request", req)

	smartModel, err := h.mapper.ToDomainUpdate(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request: model is required")
	}

	err = validation.ValidateStruct(smartModel)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	updatedModel, err := h.service.Update(ctx, smartModel)
	if err != nil {
		logger.Error("Failed to update smart model", "error", err)
		return nil, status.Error(codes.Internal, "failed to update smart model")
	}

	protoModel, err := h.mapper.ToProto(updatedModel)
	if err != nil {
		logger.Error("Failed to convert smart model to proto", "error", err)
		return nil, status.Error(codes.Internal, "failed to convert smart model to proto")
	}

	return &pb.UpdateSmartModelResponse{
		Model: protoModel,
	}, nil
}

func (h *SmartModelHandler) DeleteSmartModel(ctx context.Context, req *pb.DeleteSmartModelRequest) (*pb.DeleteSmartModelResponse, error) {
	logger.Debug("Deleting smart model", "request", req)

	err := validation.ValidateUUID(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = h.service.Delete(ctx, req.Id)
	if err != nil {
		logger.Error("Failed to delete smart model", "error", err)
		return nil, status.Error(codes.Internal, "failed to delete smart model")
	}

	return &pb.DeleteSmartModelResponse{}, nil
}
