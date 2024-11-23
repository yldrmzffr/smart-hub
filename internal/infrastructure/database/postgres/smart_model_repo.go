package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"smart-hub/internal/common/database"
	"smart-hub/internal/domain/models"
)

type PGSmartModelRepository struct {
	db *pgxpool.Pool
}

func NewPGSmartModelRepository(db database.Database) *PGSmartModelRepository {
	pgDB, ok := db.(*database.PostgresDB)
	if !ok {
		return nil
	}

	return &PGSmartModelRepository{
		db: pgDB.GetPool(),
	}
}

func (r *PGSmartModelRepository) Create(ctx context.Context, model *models.SmartModel) (*models.SmartModel, error) {
	query := `
		INSERT INTO smart_models (id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query, model.ID, model.Name, model.Description, model.Type, model.Category, model.Manufacturer, model.ModelNumber, model.Metadata, model.CreatedAt, model.UpdatedAt)

	var createdModel models.SmartModel
	err := row.Scan(
		&createdModel.ID,
		&createdModel.Name,
		&createdModel.Description,
		&createdModel.Type,
		&createdModel.Category,
		&createdModel.Manufacturer,
		&createdModel.ModelNumber,
		&createdModel.Metadata,
		&createdModel.CreatedAt,
		&createdModel.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &createdModel, nil
}

func (r *PGSmartModelRepository) GetByID(ctx context.Context, id string) (*models.SmartModel, error) {
	query := `
  SELECT id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at
  FROM smart_models
  WHERE id = $1
 `

	row := r.db.QueryRow(ctx, query, id)

	var model models.SmartModel
	err := row.Scan(
		&model.ID,
		&model.Name,
		&model.Description,
		&model.Type,
		&model.Category,
		&model.Manufacturer,
		&model.ModelNumber,
		&model.Metadata,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *PGSmartModelRepository) GetWithType(ctx context.Context, modelType string) ([]*models.SmartModel, error) {
	query := `
	  SELECT id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at
	  FROM smart_models
	  WHERE type = $1
	`

	rows, err := r.db.Query(ctx, query, modelType)
	if err != nil {
		return nil, err
	}

	var smartModels []*models.SmartModel
	for rows.Next() {
		var model models.SmartModel
		err = rows.Scan(
			&model.ID,
			&model.Name,
			&model.Description,
			&model.Type,
			&model.Category,
			&model.Manufacturer,
			&model.ModelNumber,
			&model.Metadata,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		smartModels = append(smartModels, &model)
	}

	return smartModels, nil
}

func (r *PGSmartModelRepository) GetAll(ctx context.Context) ([]*models.SmartModel, error) {
	query := `
	  SELECT id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at
	  FROM smart_models
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var smartModels []*models.SmartModel
	for rows.Next() {
		var model models.SmartModel
		err = rows.Scan(
			&model.ID,
			&model.Name,
			&model.Description,
			&model.Type,
			&model.Category,
			&model.Manufacturer,
			&model.ModelNumber,
			&model.Metadata,
			&model.CreatedAt,
			&model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		smartModels = append(smartModels, &model)
	}

	return smartModels, nil
}

func (r *PGSmartModelRepository) Update(ctx context.Context, model *models.SmartModel) (*models.SmartModel, error) {
	query := `
		UPDATE smart_models
		SET name = $2, description = $3, type = $4, category = $5, manufacturer = $6, model_number = $7, metadata = $8, updated_at = $9
		WHERE id = $1
		RETURNING id, name, description, type, category, manufacturer, model_number, metadata, created_at, updated_at
	`

	updatedModel := models.SmartModel{}

	err := r.db.QueryRow(ctx, query, model.ID,
		model.Name,
		model.Description,
		model.Type,
		model.Category,
		model.Manufacturer,
		model.ModelNumber,
		model.Metadata,
		model.UpdatedAt,
	).Scan(
		&updatedModel.ID,
		&updatedModel.Name,
		&updatedModel.Description,
		&updatedModel.Type,
		&updatedModel.Category,
		&updatedModel.Manufacturer,
		&updatedModel.ModelNumber,
		&updatedModel.Metadata,
		&updatedModel.CreatedAt,
		&updatedModel.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updatedModel, nil
}

func (r *PGSmartModelRepository) Delete(ctx context.Context, id string) error {
	// TODO: Implement soft delete

	query := `
		DELETE FROM smart_models
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
