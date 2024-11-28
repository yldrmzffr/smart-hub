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

type mockFeatureDB struct {
	pgxmock.PgxPoolIface
}

func (m *mockFeatureDB) Begin(ctx context.Context) (pgx.Tx, error) {
	return m.PgxPoolIface.Begin(ctx)
}

func (m *mockFeatureDB) Close() {}

func (m *mockFeatureDB) GetPool() database.PgxPool {
	return m
}

func TestPGSmartFeatureRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()
	now := time.Now()
	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Test Feature",
		Description:   "Test Description",
		Protocol:      models.RestProtocol,
		InterfacePath: "/test",
		Parameters:    map[string]interface{}{"key": "value"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	rows := pgxmock.NewRows([]string{
		"id", "model_id", "name", "description", "protocol",
		"interface_path", "parameters", "created_at", "updated_at",
	}).AddRow(
		feature.ID, feature.ModelID, feature.Name, feature.Description,
		feature.Protocol, feature.InterfacePath, feature.Parameters,
		feature.CreatedAt, feature.UpdatedAt,
	)

	const expectedSQL = `INSERT INTO smart_features (id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(
			feature.ID, feature.ModelID, feature.Name, feature.Description,
			feature.Protocol, feature.InterfacePath, feature.Parameters,
			feature.CreatedAt, feature.UpdatedAt,
		).
		WillReturnRows(rows)

	result, err := repo.Create(ctx, feature)
	assert.NoError(t, err)
	assert.Equal(t, feature.ID, result.ID)
	assert.Equal(t, feature.Name, result.Name)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartFeatureRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()
	now := time.Now()
	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Test Feature",
		Description:   "Test Description",
		Protocol:      models.RestProtocol,
		InterfacePath: "/test",
		Parameters:    map[string]interface{}{"key": "value"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	rows := pgxmock.NewRows([]string{
		"id", "model_id", "name", "description", "protocol",
		"interface_path", "parameters", "created_at", "updated_at",
	}).AddRow(
		feature.ID, feature.ModelID, feature.Name, feature.Description,
		feature.Protocol, feature.InterfacePath, feature.Parameters,
		feature.CreatedAt, feature.UpdatedAt,
	)

	const expectedSQL = `SELECT id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at FROM smart_features WHERE id = $1`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(feature.ID.String()).
		WillReturnRows(rows)

	result, err := repo.GetByID(ctx, feature.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, feature.ID, result.ID)
	assert.Equal(t, feature.Name, result.Name)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartFeatureRepository_GetWithModelID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()
	now := time.Now()
	modelID := uuid.New()

	features := []*models.SmartFeature{
		{
			ID:            uuid.New(),
			ModelID:       modelID,
			Name:          "Test Feature 1",
			Description:   "Test Description 1",
			Protocol:      models.RestProtocol,
			InterfacePath: "/test1",
			Parameters:    map[string]interface{}{"key": "value1"},
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            uuid.New(),
			ModelID:       modelID,
			Name:          "Test Feature 2",
			Description:   "Test Description 2",
			Protocol:      models.RestProtocol,
			InterfacePath: "/test2",
			Parameters:    map[string]interface{}{"key": "value2"},
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	rows := pgxmock.NewRows([]string{
		"id", "model_id", "name", "description", "protocol",
		"interface_path", "parameters", "created_at", "updated_at",
	})

	for _, f := range features {
		rows.AddRow(
			f.ID, f.ModelID, f.Name, f.Description,
			f.Protocol, f.InterfacePath, f.Parameters,
			f.CreatedAt, f.UpdatedAt,
		)
	}

	const expectedSQL = `SELECT id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at FROM smart_features WHERE model_id = $1`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(modelID.String()).
		WillReturnRows(rows)

	result, err := repo.GetWithModelID(ctx, modelID.String())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, features[0].ID, result[0].ID)
	assert.Equal(t, features[1].ID, result[1].ID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartFeatureRepository_GetAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()
	now := time.Now()

	features := []*models.SmartFeature{
		{
			ID:            uuid.New(),
			ModelID:       uuid.New(),
			Name:          "Test Feature 1",
			Description:   "Test Description 1",
			Protocol:      models.RestProtocol,
			InterfacePath: "/test1",
			Parameters:    map[string]interface{}{"key": "value1"},
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            uuid.New(),
			ModelID:       uuid.New(),
			Name:          "Test Feature 2",
			Description:   "Test Description 2",
			Protocol:      models.RestProtocol,
			InterfacePath: "/test2",
			Parameters:    map[string]interface{}{"key": "value2"},
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	rows := pgxmock.NewRows([]string{
		"id", "model_id", "name", "description", "protocol",
		"interface_path", "parameters", "created_at", "updated_at",
	})

	for _, f := range features {
		rows.AddRow(
			f.ID, f.ModelID, f.Name, f.Description,
			f.Protocol, f.InterfacePath, f.Parameters,
			f.CreatedAt, f.UpdatedAt,
		)
	}

	const expectedSQL = `SELECT id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at FROM smart_features`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WillReturnRows(rows)

	result, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, features[0].ID, result[0].ID)
	assert.Equal(t, features[1].ID, result[1].ID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartFeatureRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()
	now := time.Now()
	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Updated Feature",
		Description:   "Updated Description",
		Protocol:      models.RestProtocol,
		InterfacePath: "/updated",
		Parameters:    map[string]interface{}{"key": "updated_value"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	rows := pgxmock.NewRows([]string{
		"id", "model_id", "name", "description", "protocol",
		"interface_path", "parameters", "created_at", "updated_at",
	}).AddRow(
		feature.ID, feature.ModelID, feature.Name, feature.Description,
		feature.Protocol, feature.InterfacePath, feature.Parameters,
		feature.CreatedAt, feature.UpdatedAt,
	)

	const expectedSQL = `UPDATE smart_features SET name = $2, description = $3, protocol = $4, interface_path = $5, parameters = $6, updated_at = $7 WHERE id = $1 RETURNING id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(
			feature.ID, feature.Name, feature.Description,
			feature.Protocol, feature.InterfacePath, feature.Parameters,
			feature.UpdatedAt,
		).
		WillReturnRows(rows)

	result, err := repo.Update(ctx, feature)
	assert.NoError(t, err)
	assert.Equal(t, feature.ID, result.ID)
	assert.Equal(t, feature.Name, result.Name)
	assert.Equal(t, feature.Description, result.Description)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartFeatureRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()
	featureID := uuid.New()

	const expectedSQL = `DELETE FROM smart_features WHERE id = $1`

	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
		WithArgs(featureID.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(ctx, featureID.String())
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartFeatureRepository_Create_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()
	now := time.Now()
	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Test Feature",
		Description:   "Test Description",
		Protocol:      models.RestProtocol,
		InterfacePath: "/test",
		Parameters:    map[string]interface{}{"key": "value"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	const expectedSQL = `INSERT INTO smart_features (id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(
			feature.ID, feature.ModelID, feature.Name, feature.Description,
			feature.Protocol, feature.InterfacePath, feature.Parameters,
			feature.CreatedAt, feature.UpdatedAt,
		).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.Create(ctx, feature)
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartFeatureRepository_GetByID_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()
	featureID := uuid.New()

	const expectedSQL = `SELECT id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at FROM smart_features WHERE id = $1`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(featureID.String()).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.GetByID(ctx, featureID.String())
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartFeatureRepository_GetWithModelID_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()
	modelID := uuid.New()

	const expectedSQL = `SELECT id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at FROM smart_features WHERE model_id = $1`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(modelID.String()).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.GetWithModelID(ctx, modelID.String())
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartFeatureRepository_GetAll_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()

	const expectedSQL = `SELECT id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at FROM smart_features`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.GetAll(ctx)
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartFeatureRepository_Update_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()
	now := time.Now()
	feature := &models.SmartFeature{
		ID:            uuid.New(),
		ModelID:       uuid.New(),
		Name:          "Updated Feature",
		Description:   "Updated Description",
		Protocol:      models.RestProtocol,
		InterfacePath: "/updated",
		Parameters:    map[string]interface{}{"key": "updated_value"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	const expectedSQL = `UPDATE smart_features SET name = $2, description = $3, protocol = $4, interface_path = $5, parameters = $6, updated_at = $7 WHERE id = $1 RETURNING id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at`

	mock.ExpectQuery(regexp.QuoteMeta(expectedSQL)).
		WithArgs(
			feature.ID, feature.Name, feature.Description,
			feature.Protocol, feature.InterfacePath, feature.Parameters,
			feature.UpdatedAt,
		).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.Update(ctx, feature)
	assert.Error(t, err)
	assert.Nil(t, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPGSmartFeatureRepository_Delete_Failed(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db := &mockFeatureDB{mock}
	repo := NewPGSmartFeatureRepository(db)

	ctx := context.Background()
	featureID := uuid.New()

	const expectedSQL = `DELETE FROM smart_features WHERE id = $1`

	mock.ExpectExec(regexp.QuoteMeta(expectedSQL)).
		WithArgs(featureID.String()).
		WillReturnError(pgx.ErrNoRows)

	err = repo.Delete(ctx, featureID.String())
	assert.Error(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
