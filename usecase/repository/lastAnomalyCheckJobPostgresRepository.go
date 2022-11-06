package repository

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/infrastructure/database"
)

type LastAnomalyCheckJobRepository interface {
	GetLastTimestamp(ctx context.Context) (time.Time, bool, error)
	SetJobAsActive(ctx context.Context) (time.Time, error)
	SetJobAsInactive(ctx context.Context, newLastJobTime *time.Time) (time.Time, error)
}

type lastAnomalyCheckJobPostgresRepository struct {
	db *database.Postgres
}

func NewLastAnomalyCheckJobPostgresRepository(db *database.Postgres) *lastAnomalyCheckJobPostgresRepository {
	return &lastAnomalyCheckJobPostgresRepository{db: db}
}

func (r *lastAnomalyCheckJobPostgresRepository) GetLastTimestamp(ctx context.Context) (time.Time, bool, error) {
	row := r.db.Pool.QueryRow(ctx, `SELECT ts, is_active FROM last_anomaly_check_job LIMIT 1;`)

	t := time.Time{}
	isActive := false

	if err := row.Scan(&t, &isActive); err != nil {
		return time.Time{}, false, fmt.Errorf("failed to scan last anomaly check job time: %w", err)
	}

	return t, isActive, nil
}

const setIsActiveAndGetTimeStmt = `
UPDATE
    last_anomaly_check_job
SET
    is_active = true
WHERE
    name = 'anomaly_job'
RETURNING
    ts;
`

func (r *lastAnomalyCheckJobPostgresRepository) SetJobAsActive(ctx context.Context) (time.Time, error) {
	row := r.db.Pool.QueryRow(ctx, setIsActiveAndGetTimeStmt)

	t := time.Time{}

	if err := row.Scan(&t); err != nil {
		return time.Time{}, fmt.Errorf("failed to set last anomaly job as active: %w", err)
	}

	return t, nil
}

const setJobAsInactiveWithTimeStmt = `
UPDATE
    last_anomaly_check_job
SET
    is_active = false,
    ts = COALESCE($1, ts)
WHERE
    name = 'anomaly_job'
RETURNING
    ts;
`

func (r *lastAnomalyCheckJobPostgresRepository) SetJobAsInactive(ctx context.Context, newLastJobTime *time.Time) (time.Time, error) {
	row := r.db.Pool.QueryRow(ctx, setJobAsInactiveWithTimeStmt, newLastJobTime)

	t := time.Time{}

	if err := row.Scan(&t); err != nil {
		return time.Time{}, fmt.Errorf("failed to set last anomaly check job as inactive: %w", err)
	}

	return t, nil
}
