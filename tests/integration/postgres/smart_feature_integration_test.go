package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	pbFeature "smart-hub/gen/proto/smart_feature/v1"
	pbModel "smart-hub/gen/proto/smart_model/v1"
	"smart-hub/internal/application/service"
	"smart-hub/internal/infrastructure/database/postgres"
	"smart-hub/internal/presentation/grpc/handler"
	"smart-hub/internal/presentation/grpc/mapper"
	"testing"
)

func TestSmartFeatureIntegration(t *testing.T) {
	db := SetupTestDB(t)
	defer CleanupTestDB(t, db)

	modelRepo := postgres.NewPGSmartModelRepository(db)
	modelSvc := service.NewSmartModelService(modelRepo)
	modelMapper := mapper.NewSmartModelMapper()
	modelHandler := handler.NewSmartModelHandler(modelSvc, modelMapper)

	featureRepo := postgres.NewPGSmartFeatureRepository(db)
	featureSvc := service.NewSmartFeatureService(featureRepo)
	featureMapper := mapper.NewSmartFeatureMapper()
	featureHandler := handler.NewSmartFeatureHandler(featureSvc, featureMapper)

	ctx := context.Background()

	t.Run("Full CRUD Flow", func(t *testing.T) {
		// First create a model
		modelMetadata, err := structpb.NewStruct(map[string]interface{}{
			"version": "1.0",
		})
		require.NoError(t, err)

		createModelReq := &pbModel.CreateSmartModelRequest{
			Model: &pbModel.CreateSmartModelInput{
				Name:         "Test Model for Feature",
				Description:  "Test Model Description",
				Type:         pbModel.ModelType_DEVICE,
				Category:     pbModel.ModelCategory_WEARABLE,
				Manufacturer: "Test Manufacturer",
				ModelNumber:  "TEST123",
				Metadata:     modelMetadata,
			},
		}

		modelResp, err := modelHandler.CreateSmartModel(ctx, createModelReq)
		require.NoError(t, err)
		require.NotNil(t, modelResp)

		modelID := modelResp.Model.Id

		// Now create feature with valid model ID
		parameters, err := structpb.NewStruct(map[string]interface{}{
			"param1": "value1",
			"param2": 123,
		})
		require.NoError(t, err)

		createReq := &pbFeature.CreateSmartFeatureRequest{
			Feature: &pbFeature.CreateSmartFeatureInput{
				ModelId:       modelID,
				Name:          "Test Integration Feature",
				Description:   "Test Integration Description",
				Protocol:      pbFeature.ProtocolType_REST,
				InterfacePath: "/test/path",
				Parameters:    parameters,
			},
		}

		createResp, err := featureHandler.CreateSmartFeature(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, createResp)
		assert.NotEmpty(t, createResp.Feature.Id)
		assert.Equal(t, createReq.Feature.Name, createResp.Feature.Name)

		featureID := createResp.Feature.Id

		getReq := &pbFeature.GetSmartFeatureRequest{
			Id: featureID,
		}

		getResp, err := featureHandler.GetSmartFeature(ctx, getReq)
		require.NoError(t, err)
		require.NotNil(t, getResp)
		assert.Equal(t, createResp.Feature.Name, getResp.Feature.Name)

		getModelReq := &pbFeature.GetFeaturesByModelIDRequest{
			ModelId: modelID,
		}

		getModelResp, err := featureHandler.GetFeaturesByModelID(ctx, getModelReq)
		require.NoError(t, err)
		require.NotNil(t, getModelResp)
		assert.Len(t, getModelResp.Features, 1)
		assert.Equal(t, featureID, getModelResp.Features[0].Id)

		updateParameters, err := structpb.NewStruct(map[string]interface{}{
			"param1": "updated_value",
			"param2": 456,
		})
		require.NoError(t, err)

		updateReq := &pbFeature.UpdateSmartFeatureRequest{
			Feature: &pbFeature.UpdateSmartFeatureInput{
				Id:            featureID,
				Name:          "Updated Integration Feature",
				Description:   "Updated Integration Description",
				Protocol:      pbFeature.ProtocolType_GRPC,
				InterfacePath: "/updated/path",
				Parameters:    updateParameters,
			},
		}

		updateResp, err := featureHandler.UpdateSmartFeature(ctx, updateReq)
		require.NoError(t, err)
		require.NotNil(t, updateResp)
		assert.Equal(t, updateReq.Feature.Name, updateResp.Feature.Name)

		getUpdatedResp, err := featureHandler.GetSmartFeature(ctx, getReq)
		require.NoError(t, err)
		assert.Equal(t, updateReq.Feature.Name, getUpdatedResp.Feature.Name)
		assert.Equal(t, updateReq.Feature.Protocol, getUpdatedResp.Feature.Protocol)

		deleteReq := &pbFeature.DeleteSmartFeatureRequest{
			Id: featureID,
		}

		_, err = featureHandler.DeleteSmartFeature(ctx, deleteReq)
		require.NoError(t, err)

		_, err = featureHandler.GetSmartFeature(ctx, getReq)
		require.Error(t, err)
	})

	t.Run("Error Cases", func(t *testing.T) {
		_, err := featureHandler.GetSmartFeature(ctx, &pbFeature.GetSmartFeatureRequest{
			Id: uuid.New().String(),
		})
		require.Error(t, err)

		_, err = featureHandler.CreateSmartFeature(ctx, &pbFeature.CreateSmartFeatureRequest{
			Feature: &pbFeature.CreateSmartFeatureInput{
				ModelId: uuid.New().String(),
				Name:    "Test Feature",
			},
		})
		require.Error(t, err)

		_, err = featureHandler.UpdateSmartFeature(ctx, &pbFeature.UpdateSmartFeatureRequest{
			Feature: &pbFeature.UpdateSmartFeatureInput{
				Id:   uuid.New().String(),
				Name: "Updated Name",
			},
		})
		require.Error(t, err)
	})
}
