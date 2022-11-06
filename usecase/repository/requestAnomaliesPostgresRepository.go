package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/infrastructure/database"
	"go.uber.org/zap"
)

type RequestsAnomaliesRepository interface {
	UpsertAnomaly(ctx context.Context, a *entity.RequestAnomaly) error
	UpsertAnomalies(ctx context.Context, as []*entity.RequestAnomaly) error
	GetAnomaliesByRootIds(ctx context.Context, ids []string, cases []int) ([]entity.RequestAnomaly, error)
	SetCustomRequestAnomaly(ctx context.Context, rootId string, isAnomaly bool) error
}

type requestsAnomaliesPostgresRepository struct {
	db *database.Postgres
}

func NewRequestsAnomaliesPostgresRepository(db *database.Postgres) *requestsAnomaliesPostgresRepository {
	return &requestsAnomaliesPostgresRepository{db: db}
}

const upsertAnomalyStmt = `
    INSERT INTO
        requests_anomalies(root_id, is_anomaly, anomaly_cases, net_probability)
    VALUES ($1, $2, $3, $4)
    ON CONFLICT (root_id) DO UPDATE SET
        is_anomaly = $2,
        anomaly_cases = $3,
        net_probability = $4;
`

func (r *requestsAnomaliesPostgresRepository) UpsertAnomaly(ctx context.Context, a *entity.RequestAnomaly) error {
	if _, err := r.db.Pool.Exec(ctx, upsertAnomalyStmt, a.RootID, a.IsAnomaly, a.AnomalyCases); err != nil {
		return fmt.Errorf("failed to upsert request anomaly: %w", err)
	}

	return nil
}

func (r *requestsAnomaliesPostgresRepository) UpsertAnomalies(ctx context.Context, as []*entity.RequestAnomaly) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}

	defer func() {
		if err = tx.Rollback(context.Background()); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			zap.S().With("context", "RequestsAnomaliesPostgresRepository").Errorw("TX rollback error", "error", err.Error())
		}
	}()

	batch := &pgx.Batch{}

	for _, a := range as {
		batch.Queue(
			upsertAnomalyStmt,
			a.RootID,
			a.IsAnomaly,
			a.AnomalyCases,
			a.NetProbability,
		)
	}

	batchResult := tx.SendBatch(ctx, batch)

	if err := batchResult.Close(); err != nil {
		return fmt.Errorf("failed to close requests anomalies batch result:%w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit requests anomalies batches: %w", err)
	}

	return nil
}

const getAnomaliesByRootIdsStmt = `
SELECT
	a.root_id AS root_id,
	(ca.is_anomaly IS NULL AND a.is_anomaly) OR (ca.is_anomaly IS NOT NULL AND ca.is_anomaly) AS is_anomaly,
	(
		CASE
        	WHEN ca.is_anomaly IS NOT NULL AND NOT ca.is_anomaly THEN '{}'::int[]
    	ELSE
        	a.anomaly_cases
    	END
    ) as anomaly_cases,
    ca.root_id IS NOT NULL AS is_custom,
    a.net_probability AS net_probability
FROM requests_anomalies AS a
LEFT JOIN custom_requests_anomalies AS ca ON
	ca.root_id = a.root_id
    WHERE
        a.root_id = ANY($1) AND
        (array_length($2::int[], 1) IS NULL OR a.anomaly_cases && ($2));
`

func (r *requestsAnomaliesPostgresRepository) GetAnomaliesByRootIds(
	ctx context.Context,
	ids []string,
	cases []int,
) ([]entity.RequestAnomaly, error) {
	anomalies := make([]entity.RequestAnomaly, 0, len(ids))

	rows, err := r.db.Pool.Query(ctx, getAnomaliesByRootIdsStmt, ids, cases)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return anomalies, nil
		}

		return nil, fmt.Errorf("failed to get anomalies: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		anomaly := entity.RequestAnomaly{}

		if err := rows.Scan(
			&anomaly.RootID,
			&anomaly.IsAnomaly,
			&anomaly.AnomalyCases,
			&anomaly.IsCustom,
			&anomaly.NetProbability,
		); err != nil {
			return nil, fmt.Errorf("failed to scan anomalies: %w", err)
		}

		anomalies = append(anomalies, anomaly)
	}

	return anomalies, nil
}

const setCustomRequestAnomalyStmt = `
    INSERT INTO
        custom_requests_anomalies(root_id, is_anomaly)
    VALUES ($1, $2)
    ON CONFLICT (root_id) DO UPDATE SET
        is_anomaly = $2;
`

func (r *requestsAnomaliesPostgresRepository) SetCustomRequestAnomaly(ctx context.Context, rootId string, isAnomaly bool) error {
	if _, err := r.db.Pool.Exec(ctx, setCustomRequestAnomalyStmt, rootId, isAnomaly); err != nil {
		return fmt.Errorf("failed to upsert custom request anomaly: %w", err)
	}

	return nil
}
