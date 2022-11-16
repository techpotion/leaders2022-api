package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	CountAnomalies(ctx context.Context, dateFrom, dateTo time.Time, region string) (int, float32, error)
	CountAnomaliesGroupped(ctx context.Context, dateFrom, dateTo time.Time, region string) ([]*entity.CountAnomaliesGroupped, error)
	GetAnomaliesAmountDynamics(ctx context.Context, dateFrom, dateTo time.Time, region string) ([]*entity.AnomaliesDynamics, error)
	CountAnomaliesByOwnerCompanies(
		ctx context.Context,
		dateFrom, dateTo time.Time,
		region string,
	) ([]*entity.AnomaliesByOwnerCompany, error)
	CountAnomaliesByServingCompanies(
		ctx context.Context,
		dateFrom, dateTo time.Time,
		region string,
	) ([]*entity.AnomaliesByServingCompany, error)
	CountAnomaliesByDeffectCategories(
		ctx context.Context,
		dateFrom, dateTo time.Time,
		region string,
	) ([]*entity.AnomaliesByDeffectCategory, error)
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

const countAnomaliesStmt = `
    with r as (
        select
            r.*,
            ra.is_anomaly,
            ra.net_probability,
            ra.anomaly_cases
        from requests as r
        left join requests_anomalies as ra on
            ra.root_id = r.root_id
        where
            r.hood = $1 AND
            r.closure_date BETWEEN $2 AND $3
    ),
    r_a as (
        select count(*) as count from r
        where
            r.is_anomaly
    ),
    r_not_a as (
        select count(*) as count from r
        where
            not r.is_anomaly
    )

    select
        r_a.count as count,
        r_a.count::float/(r_a.count + r_not_a.count) as perc
    from r_not_a
    left join r_a on true
	where
        r_a.count + r_not_a.count > 0;
`

func (r *requestsAnomaliesPostgresRepository) CountAnomalies(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) (count int, perc float32, err error) {
	row := r.db.Pool.QueryRow(ctx, countAnomaliesStmt, region, dateFrom, dateTo)

	if err = row.Scan(
		&count,
		&perc,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, 0.0, nil
		}

		return 0, 0.0, fmt.Errorf("failed to get anomalies count: %w", err)
	}

	return count, perc, nil
}

const countAnomaliesGrouppedStmt = `
    with anomalies_unset as (
        select
            unnest(ra.anomaly_cases) as cases
        from requests as r
        left join requests_anomalies as ra on
            ra.root_id = r.root_id
        where
            ra.is_anomaly AND
            r.hood = $1 AND
            r.closure_date BETWEEN $2 AND $3
    )

    select
    	cases                                                    as type,
    	count(*)                                                 as count,
    	count(*)::float / (select count(*) from anomalies_unset) as percent
    from
        anomalies_unset
    where (select count(*) from anomalies_unset) > 0
    group by type
    order by count desc;
`

const groupMapMinSize = 5

func (r *requestsAnomaliesPostgresRepository) CountAnomaliesGroupped(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.CountAnomaliesGroupped, error) {
	m := make([]*entity.CountAnomaliesGroupped, 0, groupMapMinSize)

	rows, err := r.db.Pool.Query(ctx, countAnomaliesGrouppedStmt, region, dateFrom, dateTo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return m, nil
		}

		return nil, fmt.Errorf("failed to count anomalies groupped: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		c := &entity.CountAnomaliesGroupped{}

		if err := rows.Scan(&c.Type, &c.Count, &c.Percent); err != nil {
			return nil, fmt.Errorf("failed to scan anomalies groupped: %w", err)
		}

		m = append(m, c)
	}

	return m, nil
}

const getAnomaliesAmountDynamicsStmt = `
with all_requests as (
    select DATE_TRUNC('day', r.closure_date) as d, 0 as count from requests as r
    inner join requests_anomalies as ra on
        ra.root_id = r.root_id
    where
        r.hood = $1 AND
        r.closure_date BETWEEN $2 AND $3
    group by DATE_TRUNC('day', r.closure_date)
),
anomaly_requests as (
    select DATE_TRUNC('day', r.closure_date) as d, count(*) from requests as r
    inner join requests_anomalies as ra on
        ra.root_id = r.root_id and
        ra.is_anomaly
    where
        r.hood = $1 AND
        r.closure_date BETWEEN $2 AND $3
    group by DATE_TRUNC('day', r.closure_date)
)

select d, max(count) from (
	select all_r.d, all_r.count from all_requests as all_r
	union
	select an_r.d, an_r.count from anomaly_requests as an_r
) as tmp
group by d
order by d
`

func (r *requestsAnomaliesPostgresRepository) GetAnomaliesAmountDynamics(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesDynamics, error) {
	ads := make([]*entity.AnomaliesDynamics, 0, groupMapMinSize)

	rows, err := r.db.Pool.Query(ctx, getAnomaliesAmountDynamicsStmt, region, dateFrom, dateTo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to count anomalies groupped: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		ad := &entity.AnomaliesDynamics{}

		if err := rows.Scan(&ad.Day, &ad.Amount); err != nil {
			return nil, fmt.Errorf("failed to scan anomalies groupped: %w", err)
		}

		ads = append(ads, ad)
	}

	return ads, nil
}

const countAnomaliesCountByOwnerCompaniesStmt = `
with all_requests as (
  select r.owner_company, count(*) from requests as r
  inner join requests_anomalies as ra on
    ra.root_id = r.root_id
  where
    r.hood = $1 AND
    r.closure_date BETWEEN $2 AND $3
  group by r.owner_company
),
anomaly_requests as (
  select r.owner_company, count(*) from requests as r
  inner join requests_anomalies as ra on
    ra.root_id = r.root_id and
    ra.is_anomaly
  where
    r.hood = $1 AND
    r.closure_date BETWEEN $2 AND $3
  group by r.owner_company
)

select
	all_r.owner_company                                 as owner_company,
	COALESCE(an_r.count, 0)                             as count,
	COALESCE(an_r.count::float / all_r.count::float, 0) as percent
from
	all_requests as all_r
left join anomaly_requests as an_r on
  all_r.owner_company = an_r.owner_company
where
	all_r.owner_company != ''
order by count desc;
`

func (r *requestsAnomaliesPostgresRepository) CountAnomaliesByOwnerCompanies(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesByOwnerCompany, error) {
	as := make([]*entity.AnomaliesByOwnerCompany, 0, groupMapMinSize)

	rows, err := r.db.Pool.Query(ctx, countAnomaliesCountByOwnerCompaniesStmt, region, dateFrom, dateTo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to count anomalies by owner companies: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		a := &entity.AnomaliesByOwnerCompany{}

		if err := rows.Scan(&a.OwnerCompany, &a.Count, &a.Percent); err != nil {
			return nil, fmt.Errorf("failed to count anomalies by owner companies: %w", err)
		}

		as = append(as, a)
	}

	return as, nil
}

const countAnomaliesCountByServingCompaniesStmt = `
with all_requests as (
  select r.serving_company, count(*) from requests as r
  inner join requests_anomalies as ra on
    ra.root_id = r.root_id
  where
    r.hood = $1 AND
    r.closure_date BETWEEN $2 AND $3
  group by r.serving_company
),
anomaly_requests as (
  select r.serving_company, count(*) from requests as r
  inner join requests_anomalies as ra on
    ra.root_id = r.root_id and
    ra.is_anomaly
  where
    r.hood = $1 AND
    r.closure_date BETWEEN $2 AND $3
  group by r.serving_company
)

select
	all_r.serving_company                               as serving_company,
	COALESCE(an_r.count, 0)                             as count,
	COALESCE(an_r.count::float / all_r.count::float, 0) as percent
from
	all_requests as all_r
left join anomaly_requests as an_r on
  all_r.serving_company = an_r.serving_company
where
	all_r.serving_company != ''
order by count desc;
`

func (r *requestsAnomaliesPostgresRepository) CountAnomaliesByServingCompanies(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesByServingCompany, error) {
	as := make([]*entity.AnomaliesByServingCompany, 0, groupMapMinSize)

	rows, err := r.db.Pool.Query(ctx, countAnomaliesCountByServingCompaniesStmt, region, dateFrom, dateTo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to count anomalies by serving companies: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		a := &entity.AnomaliesByServingCompany{}

		if err := rows.Scan(&a.ServingCompany, &a.Count, &a.Percent); err != nil {
			return nil, fmt.Errorf("failed to count anomalies by serving companies: %w", err)
		}

		as = append(as, a)
	}

	return as, nil
}

const countAnomaliesCountByDeffectCategoriesStmt = `
with all_requests as (
  select r.deffect_category_name, count(*) from requests as r
  inner join requests_anomalies as ra on
    ra.root_id = r.root_id
  where
    r.hood = $1 AND
    r.closure_date BETWEEN $2 AND $3
  group by r.deffect_category_name
),
anomaly_requests as (
  select r.deffect_category_name, count(*) from requests as r
  inner join requests_anomalies as ra on
    ra.root_id = r.root_id and
    ra.is_anomaly
  where
    r.hood = $1 AND
    r.closure_date BETWEEN $2 AND $3
  group by r.deffect_category_name
)

select
	all_r.deffect_category_name                         as deffect_category_name,
	COALESCE(an_r.count, 0)                             as count,
	COALESCE(an_r.count::float / all_r.count::float, 0) as percent
from
	all_requests as all_r
left join anomaly_requests as an_r on
  all_r.deffect_category_name = an_r.deffect_category_name
where
	all_r.deffect_category_name != ''
order by count desc;
`

func (r *requestsAnomaliesPostgresRepository) CountAnomaliesByDeffectCategories(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	region string,
) ([]*entity.AnomaliesByDeffectCategory, error) {
	as := make([]*entity.AnomaliesByDeffectCategory, 0, groupMapMinSize)

	rows, err := r.db.Pool.Query(ctx, countAnomaliesCountByDeffectCategoriesStmt, region, dateFrom, dateTo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to count anomalies by deffect categories: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		a := &entity.AnomaliesByDeffectCategory{}

		if err := rows.Scan(&a.DeffectCategory, &a.Count, &a.Percent); err != nil {
			return nil, fmt.Errorf("failed to count anomalies by deffect categories: %w", err)
		}

		as = append(as, a)
	}

	return as, nil
}
