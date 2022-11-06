package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"gitlab.com/techpotion/leadershack2022/api/infrastructure/database"
)

type MoscowRegionsRepository interface {
	GetUniqueRegions(ctx context.Context) ([]string, error)
}

type moscowRegionsPostgresRepository struct {
	db *database.Postgres
}

func NewMoscowRegionsPostgresRepository(db *database.Postgres) *moscowRegionsPostgresRepository {
	return &moscowRegionsPostgresRepository{db}
}

const maxHoods = 125

func (r *moscowRegionsPostgresRepository) GetUniqueRegions(ctx context.Context) ([]string, error) {
	hoods := make([]string, 0, maxHoods)

	rows, err := r.db.Pool.Query(
		ctx,
		`SELECT DISTINCT(hood) FROM requests WHERE hood IS NOT NULL;`,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return hoods, nil
		}

		return nil, fmt.Errorf("failed to get unqiue regions: %w", err)
	}

	defer rows.Close()

	var hood string

	for rows.Next() {
		if err = rows.Scan(&hood); err != nil {
			return nil, fmt.Errorf("failed to scan regions: %w", err)
		}

		hoods = append(hoods, hood)
	}

	return hoods, nil
}
