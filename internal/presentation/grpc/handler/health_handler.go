package handler

import (
	"context"
	pb "smart-hub/gen/proto/health/v1"
	"smart-hub/internal/common/database"
	"smart-hub/internal/common/logger"
)

type HealthHandler struct {
	db database.Database
}

func NewHealthHandler(db database.Database) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

func (h *HealthHandler) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	logger.Debug("Health check requested", "service", req.Service)

	err := h.db.Ping(ctx)
	if err != nil {
		logger.Error("Database health check failed", "error", err)
		return &pb.HealthCheckResponse{
			Status: pb.HealthCheckResponse_SERVING_STATUS_NOT_SERVING,
		}, nil
	}

	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING_STATUS_SERVING,
	}, nil
}
