package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"smart-hub/internal/common/database"
	"smart-hub/internal/domain/models"
	"testing"
	"time"
)

type mockModelDB struct {
	pgxmock.PgxPoolIface
}

func (m *mockModelDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return m.PgxPoolIface.Begin(ctx)
}

func (m *mockModelDB) Close() {}

func (m *mockModelDB) GetPool() database.PgxPool {
	return m
}

func TestPGSmartModelRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()
	now := time.Now()
	model := &models.SmartModel{
		ID:           uuid.New(),
		Name:         "Test Model",
		Description:  "Test Description",
		Type:         models.DeviceType,
		Category:     models.WearableCategory,
		Manufacturer: "Test Manufacturer",
		ModelNumber:  "TEST123",
		Metadata:     map[string]interface{}{"key": "value"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	rows := pgxmock.NewRows([]string{
		"id", "name", "description", "type", "category",
		"manufacturer", "model_number", "metadata", "created_at", "updated_at",
	}).AddRow(
		model.ID, model.Name, model.Description, model.Type, model.Category,
		model.Manufacturer, model.ModelNumber, model.Metadata,
		model.CreatedAt, model.UpdatedAt,
	)

	const expectedSQL = `INSERT INTO smart_models (id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(
			model.ID, model.Name, model.Description, model.Type, model.Category,
			model.Manufacturer, model.ModelNumber, model.Metadata,
			model.CreatedAt, model.UpdatedAt,
		).
		WillReturnRows(rows)

	result, err := repo.Create(ctx, model)
	assert.NoError(t, err)
	assert.Equal(t, model.ID, result.ID)
	assert.Equal(t, model.Name, result.Name)
	assert.Equal(t, model.Description, result.Description)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartModelRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()
	now := time.Now()
	model := &models.SmartModel{
		ID:           uuid.New(),
		Name:         "Test Model",
		Description:  "Test Description",
		Type:         models.DeviceType,
		Category:     models.WearableCategory,
		Manufacturer: "Test Manufacturer",
		ModelNumber:  "TEST123",
		Metadata:     map[string]interface{}{"key": "value"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	rows := pgxmock.NewRows([]string{
		"id", "name", "description", "type", "category",
		"manufacturer", "model_number", "metadata", "created_at", "updated_at",
	}).AddRow(
		model.ID, model.Name, model.Description, model.Type, model.Category,
		model.Manufacturer, model.ModelNumber, model.Metadata,
		model.CreatedAt, model.UpdatedAt,
	)

	const expectedSQL = `SELECT id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at FROM smart_models WHERE id = $1`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(model.ID.String()).
		WillReturnRows(rows)

	result, err := repo.GetByID(ctx, model.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, model.ID, result.ID)
	assert.Equal(t, model.Name, result.Name)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartModelRepository_GetWithType(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()
	now := time.Now()

	testModels := []*models.SmartModel{
		{
			ID:           uuid.New(),
			Name:         "Test Model 1",
			Description:  "Test Description 1",
			Type:         models.DeviceType,
			Category:     models.WearableCategory,
			Manufacturer: "Test Manufacturer 1",
			ModelNumber:  "TEST123",
			Metadata:     map[string]interface{}{"key": "value1"},
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New(),
			Name:         "Test Model 2",
			Description:  "Test Description 2",
			Type:         models.DeviceType,
			Category:     models.CameraCategory,
			Manufacturer: "Test Manufacturer 2",
			ModelNumber:  "TEST456",
			Metadata:     map[string]interface{}{"key": "value2"},
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	rows := pgxmock.NewRows([]string{
		"id", "name", "description", "type", "category",
		"manufacturer", "model_number", "metadata", "created_at", "updated_at",
	})

	for _, m := range testModels {
		rows.AddRow(
			m.ID, m.Name, m.Description, m.Type, m.Category,
			m.Manufacturer, m.ModelNumber, m.Metadata,
			m.CreatedAt, m.UpdatedAt,
		)
	}

	const expectedSQL = `SELECT id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at FROM smart_models WHERE type = $1`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(models.DeviceType).
		WillReturnRows(rows)

	result, err := repo.GetWithType(ctx, models.DeviceType)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, testModels[0].ID, result[0].ID)
	assert.Equal(t, testModels[1].ID, result[1].ID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartModelRepository_GetAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()
	now := time.Now()

	testModels := []*models.SmartModel{
		{
			ID:           uuid.New(),
			Name:         "Test Model 1",
			Description:  "Test Description 1",
			Type:         models.DeviceType,
			Category:     models.WearableCategory,
			Manufacturer: "Test Manufacturer 1",
			ModelNumber:  "TEST123",
			Metadata:     map[string]interface{}{"key": "value1"},
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           uuid.New(),
			Name:         "Test Model 2",
			Description:  "Test Description 2",
			Type:         models.ServiceType,
			Category:     models.CameraCategory,
			Manufacturer: "Test Manufacturer 2",
			ModelNumber:  "TEST456",
			Metadata:     map[string]interface{}{"key": "value2"},
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	rows := pgxmock.NewRows([]string{
		"id", "name", "description", "type", "category",
		"manufacturer", "model_number", "metadata", "created_at", "updated_at",
	})

	for _, m := range testModels {
		rows.AddRow(
			m.ID, m.Name, m.Description, m.Type, m.Category,
			m.Manufacturer, m.ModelNumber, m.Metadata,
			m.CreatedAt, m.UpdatedAt,
		)
	}

	const expectedSQL = `SELECT id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at FROM smart_models`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WillReturnRows(rows)

	result, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, testModels[0].ID, result[0].ID)
	assert.Equal(t, testModels[1].ID, result[1].ID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartModelRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()
	now := time.Now()
	model := &models.SmartModel{
		ID:           uuid.New(),
		Name:         "Updated Model",
		Description:  "Updated Description",
		Type:         models.DeviceType,
		Category:     models.WearableCategory,
		Manufacturer: "Updated Manufacturer",
		ModelNumber:  "UPDATE123",
		Metadata:     map[string]interface{}{"key": "updated_value"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	rows := pgxmock.NewRows([]string{
		"id", "name", "description", "type", "category",
		"manufacturer", "model_number", "metadata", "created_at", "updated_at",
	}).AddRow(
		model.ID, model.Name, model.Description, model.Type, model.Category,
		model.Manufacturer, model.ModelNumber, model.Metadata,
		model.CreatedAt, model.UpdatedAt,
	)

	const expectedSQL = `UPDATE smart_models SET name = $2, description = $3, type = $4, category = $5, manufacturer = $6, model_number = $7, metadata = $8, updated_at = $9 WHERE id = $1 RETURNING id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(
			model.ID, model.Name, model.Description, model.Type, model.Category,
			model.Manufacturer, model.ModelNumber, model.Metadata, model.UpdatedAt,
		).
		WillReturnRows(rows)

	result, err := repo.Update(ctx, model)
	assert.NoError(t, err)
	assert.Equal(t, model.ID, result.ID)
	assert.Equal(t, model.Name, result.Name)
	assert.Equal(t, model.Description, result.Description)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartModelRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()
	modelID := uuid.New()

	const expectedSQL = `DELETE FROM smart_models WHERE id = $1`

	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
		WithArgs(modelID.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(ctx, modelID.String())
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartModelRepository_Create_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()
	now := time.Now()
	model := &models.SmartModel{
		ID:           uuid.New(),
		Name:         "Test Model",
		Description:  "Test Description",
		Type:         models.DeviceType,
		Category:     models.WearableCategory,
		Manufacturer: "Test Manufacturer",
		ModelNumber:  "TEST123",
		Metadata:     map[string]interface{}{"key": "value"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	const expectedSQL = `INSERT INTO smart_models (id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(
			model.ID, model.Name, model.Description, model.Type, model.Category,
			model.Manufacturer, model.ModelNumber, model.Metadata,
			model.CreatedAt, model.UpdatedAt,
		).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.Create(ctx, model)
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartModelRepository_GetByID_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()
	modelID := uuid.New()

	const expectedSQL = `SELECT id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at FROM smart_models WHERE id = $1`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(modelID.String()).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.GetByID(ctx, modelID.String())
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartModelRepository_GetWithType_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()

	const expectedSQL = `SELECT id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at FROM smart_models WHERE type = $1`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(models.DeviceType).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.GetWithType(ctx, models.DeviceType)
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartModelRepository_GetAll_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()

	const expectedSQL = `SELECT id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at FROM smart_models`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.GetAll(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartModelRepository_Update_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()
	now := time.Now()
	model := &models.SmartModel{
		ID:           uuid.New(),
		Name:         "Updated Model",
		Description:  "Updated Description",
		Type:         models.DeviceType,
		Category:     models.WearableCategory,
		Manufacturer: "Updated Manufacturer",
		ModelNumber:  "UPDATE123",
		Metadata:     map[string]interface{}{"key": "updated_value"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	const expectedSQL = `UPDATE smart_models SET name = $2, description = $3, type = $4, category = $5, manufacturer = $6, model_number = $7, metadata = $8, updated_at = $9 WHERE id = $1 RETURNING id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(
			model.ID, model.Name, model.Description, model.Type, model.Category,
			model.Manufacturer, model.ModelNumber, model.Metadata, model.UpdatedAt,
		).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.Update(ctx, model)
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartModelRepository_Delete_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockModelDB{mock}
	repo := NewPGSmartModelRepository(db)

	ctx := context.Background()
	modelID := uuid.New()

	const expectedSQL = `DELETE FROM smart_models WHERE id = $1`

	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
		WithArgs(modelID.String()).
		WillReturnError(pgx.ErrNoRows)

	err = repo.Delete(ctx, modelID.String())
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
