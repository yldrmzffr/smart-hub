package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	pb "smart-hub/gen/proto/smart_model/v1"
	"smart-hub/internal/application/service"
	"smart-hub/internal/infrastructure/database/postgres"
	"smart-hub/internal/presentation/grpc/handler"
	"smart-hub/internal/presentation/grpc/mapper"
	"testing"
)

func TestSmartModelIntegration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)

	repo := postgres.NewPGSmartModelRepository(db)
	svc := service.NewSmartModelService(repo)
	modelMapper := mapper.NewSmartModelMapper()
	handler := handler.NewSmartModelHandler(svc, modelMapper)

	ctx := context.Background()

	t.Run("Full CRUD Flow", func(t *testing.T) {
		metadata, err := structpb.NewStruct(map[string]interface{}{
			"version": "1.0",
			"vendor":  "test",
		})
		require.NoError(t, err)

		createReq := &pb.CreateSmartModelRequest{
			Model: &pb.CreateSmartModelInput{
				Name:         "Test Integration Model",
				Description:  "Test Integration Description",
				Type:         pb.ModelType_DEVICE,
				Category:     pb.ModelCategory_WEARABLE,
				Manufacturer: "Test Manufacturer",
				ModelNumber:  "TEST123",
				Metadata:     metadata,
			},
		}

		createResp, err := handler.CreateSmartModel(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, createResp)
		assert.NotEmpty(t, createResp.Model.Id)
		assert.Equal(t, createReq.Model.Name, createResp.Model.Name)

		modelID := createResp.Model.Id

		getReq := &pb.GetSmartModelRequest{
			Id: modelID,
		}

		getResp, err := handler.GetSmartModel(ctx, getReq)
		require.NoError(t, err)
		require.NotNil(t, getResp)
		assert.Equal(t, createResp.Model.Name, getResp.Model.Name)

		updateReq := &pb.UpdateSmartModelRequest{
			Model: &pb.UpdateSmartModelInput{
				Id:           modelID,
				Name:         "Updated Integration Model",
				Description:  "Updated Integration Description",
				Type:         pb.ModelType_DEVICE,
				Category:     pb.ModelCategory_CAMERA,
				Manufacturer: "Updated Manufacturer",
				ModelNumber:  "UPDATE123",
				Metadata:     metadata,
			},
		}

		updateResp, err := handler.UpdateSmartModel(ctx, updateReq)
		require.NoError(t, err)
		require.NotNil(t, updateResp)
		assert.Equal(t, updateReq.Model.Name, updateResp.Model.Name)

		getUpdatedResp, err := handler.GetSmartModel(ctx, getReq)
		require.NoError(t, err)
		assert.Equal(t, updateReq.Model.Name, getUpdatedResp.Model.Name)

		deleteReq := &pb.DeleteSmartModelRequest{
			Id: modelID,
		}

		_, err = handler.DeleteSmartModel(ctx, deleteReq)
		require.NoError(t, err)

		_, err = handler.GetSmartModel(ctx, getReq)
		require.Error(t, err)
	})

	t.Run("Error Cases", func(t *testing.T) {
		_, err := handler.GetSmartModel(ctx, &pb.GetSmartModelRequest{
			Id: uuid.New().String(),
		})
		require.Error(t, err)

		_, err = handler.CreateSmartModel(ctx, &pb.CreateSmartModelRequest{
			Model: &pb.CreateSmartModelInput{
				Name: "",
			},
		})
		require.Error(t, err)

		_, err = handler.UpdateSmartModel(ctx, &pb.UpdateSmartModelRequest{
			Model: &pb.UpdateSmartModelInput{
				Id:   uuid.New().String(),
				Name: "Updated Name",
			},
		})
		require.Error(t, err)
	})
}
