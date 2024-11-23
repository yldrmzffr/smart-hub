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
	return nil, nil
}

func (r *PGSmartModelRepository) GetByID(ctx context.Context, id int) (*models.SmartModel, error) {
	return nil, nil
}

func (r *PGSmartModelRepository) GetWithType(ctx context.Context, modelType string) ([]*models.SmartModel, error) {
	return nil, nil
}

func (r *PGSmartModelRepository) GetAll(ctx context.Context) ([]*models.SmartModel, error) {
	return nil, nil
}

func (r *PGSmartModelRepository) Update(ctx context.Context, model *models.SmartModel) (*models.SmartModel, error) {
	return nil, nil
}

func (r *PGSmartModelRepository) Delete(ctx context.Context, id int) error {
	return nil
}
