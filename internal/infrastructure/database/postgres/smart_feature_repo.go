package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"smart-hub/internal/common/database"
	"smart-hub/internal/domain/models"
)

type PGSmartFeatureRepository struct {
	db *pgxpool.Pool
}

func NewPGSmartFeatureRepository(db database.Database) *PGSmartFeatureRepository {
	pgDB, ok := db.(*database.PostgresDB)
	if !ok {
		return nil
	}

	return &PGSmartFeatureRepository{
		db: pgDB.GetPool(),
	}
}

func (r *PGSmartFeatureRepository) Create(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error) {
	query := `
		INSERT INTO smart_features (id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)	
		RETURNING id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query, feature.ID, feature.ModelID, feature.Name, feature.Description, feature.Protocol, feature.InterfacePath, feature.Parameters, feature.CreatedAt, feature.UpdatedAt)

	var createdFeature models.SmartFeature
	err := row.Scan(
		&createdFeature.ID,
		&createdFeature.ModelID,
		&createdFeature.Name,
		&createdFeature.Description,
		&createdFeature.Protocol,
		&createdFeature.InterfacePath,
		&createdFeature.Parameters,
		&createdFeature.CreatedAt,
		&createdFeature.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &createdFeature, nil
}

func (r *PGSmartFeatureRepository) GetByID(ctx context.Context, id string) (*models.SmartFeature, error) {
	query := `
		SELECT id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at
		FROM smart_features
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var feature models.SmartFeature
	err := row.Scan(
		&feature.ID,
		&feature.ModelID,
		&feature.Name,
		&feature.Description,
		&feature.Protocol,
		&feature.InterfacePath,
		&feature.Parameters,
		&feature.CreatedAt,
		&feature.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &feature, nil
}

func (r *PGSmartFeatureRepository) GetWithModelID(ctx context.Context, modelID string) ([]*models.SmartFeature, error) {
	query := `
		SELECT id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at
		FROM smart_features
		WHERE model_id = $1
	`

	rows, err := r.db.Query(ctx, query, modelID)
	if err != nil {
		return nil, err
	}

	var features []*models.SmartFeature
	for rows.Next() {
		var feature models.SmartFeature
		err = rows.Scan(
			&feature.ID,
			&feature.ModelID,
			&feature.Name,
			&feature.Description,
			&feature.Protocol,
			&feature.InterfacePath,
			&feature.Parameters,
			&feature.CreatedAt,
			&feature.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		features = append(features, &feature)
	}

	return features, nil
}

func (r *PGSmartFeatureRepository) GetAll(ctx context.Context) ([]*models.SmartFeature, error) {
	query := `
		SELECT id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at
		FROM smart_features
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var features []*models.SmartFeature

	for rows.Next() {
		var feature models.SmartFeature
		err = rows.Scan(
			&feature.ID,
			&feature.ModelID,
			&feature.Name,
			&feature.Description,
			&feature.Protocol,
			&feature.InterfacePath,
			&feature.Parameters,
			&feature.CreatedAt,
			&feature.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		features = append(features, &feature)
	}

	return features, nil
}

func (r *PGSmartFeatureRepository) Update(ctx context.Context, feature *models.SmartFeature) (*models.SmartFeature, error) {
	query := `
		UPDATE smart_features
		SET name = $2, description = $3, protocol = $4, interface_path = $5, parameters = $6, updated_at = $7
		WHERE id = $1
		RETURNING id, model_id, name, description, protocol, interface_path, parameters, created_at, updated_at
	`

	updatedFeature := models.SmartFeature{}

	err := r.db.QueryRow(ctx, query,
		feature.ID,
		feature.Name,
		feature.Description,
		feature.Protocol,
		feature.InterfacePath,
		feature.Parameters,
		feature.UpdatedAt,
	).Scan(
		&updatedFeature.ID,
		&updatedFeature.ModelID,
		&updatedFeature.Name,
		&updatedFeature.Description,
		&updatedFeature.Protocol,
		&updatedFeature.InterfacePath,
		&updatedFeature.Parameters,
		&updatedFeature.CreatedAt,
		&updatedFeature.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updatedFeature, nil
}

func (r *PGSmartFeatureRepository) Delete(ctx context.Context, id string) error {
	// TODO: Implement soft delete
	query := `
		DELETE FROM smart_features
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
