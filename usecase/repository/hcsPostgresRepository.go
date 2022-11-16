package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"gitlab.com/techpotion/leadershack2022/api/entity"
	"gitlab.com/techpotion/leadershack2022/api/infrastructure/database"
	"gitlab.com/techpotion/leadershack2022/api/usecase/dto"
)

type HCSRepository interface {
	CountPointsWithFilters(ctx context.Context, filters *dto.CountPointsRequestDTO) (int, error)
	GetPointsWithFilters(ctx context.Context, filters *dto.GetPointsRequestQueryDTO) ([]entity.HCSPoint, error)
	GetRequestsByPointsIds(ctx context.Context, ids []string) ([]*entity.Request, error)
	GetPrimaryRequestDatetime(ctx context.Context, current *entity.Request) (*time.Time, error)
	CountRequestsByClosureTime(ctx context.Context, from, to time.Time) (int, error)
	GetRequestsByClosureTime(ctx context.Context, from, to time.Time, limit, offset int) ([]*entity.Request, error)
	CountRequestsFull(ctx context.Context, filters *dto.CountRequestsFullRequestDTO) (int, error)
	GetRequestsFull(ctx context.Context, filters *dto.GetRequestsFullRequestDTO) ([]*entity.RequestFull, error)
	GetUniqueRegions(ctx context.Context) ([]string, error)
	GetUniqueServingCompanies(ctx context.Context) ([]string, error)
	GetUniqueOwnerCompanies(ctx context.Context) ([]string, error)
	GetUniqueDeffectCategories(ctx context.Context) ([]string, error)
	GetUniqueWorkTypes(ctx context.Context) ([]string, error)
	GetRequestsByDispatcher(ctx context.Context, dateFrom, dateTo time.Time, dispNumber string) ([]*entity.Request, error)
	GetUniqueDispatchers(ctx context.Context) ([]string, error)
	GetRegionArea(ctx context.Context, region string) (string, error)
}

type hcsPostgresRepository struct {
	db *database.Postgres
}

func NewHCSPostgresRepository(db *database.Postgres) *hcsPostgresRepository {
	return &hcsPostgresRepository{db: db}
}

const countPointsWithFiltersStmt = `
    SELECT
        count(*)
    FROM addresses AS a
    INNER JOIN requests AS m on
        m.adress_unom = a.unom
    WHERE
        m.hood = $1 AND
        m.closure_date BETWEEN COALESCE($2, m.closure_date) AND COALESCE($3, m.closure_date) AND
        rectangle_contains($4, $5, $6, $7, a.center_x, a.center_y) AND
        (array_length($8::varchar[], 1)  IS NULL OR ARRAY[m.serving_company]       && ($8::varchar[])) AND
        (array_length($9::varchar[], 1)  IS NULL OR ARRAY[m.efficiency]            &&  ($9::varchar[])) AND
        (array_length($10::varchar[], 1) IS NULL OR ARRAY[m.grade_for_service]     && ($10::varchar[])) AND
        (array_length($11::varchar[], 1) IS NULL OR ARRAY[m.urgency_category]      && ($11::varchar[])) AND
        (array_length($12::varchar[], 1) IS NULL OR ARRAY[m.work_type_done]        && ($12::varchar[])) AND
        (array_length($13::varchar[], 1) IS NULL OR ARRAY[m.deffect_category_name] && ($13::varchar[])) AND
        (array_length($14::varchar[], 1) IS NULL OR ARRAY[m.owner_company]         && ($14::varchar[]))
;
`

func (r *hcsPostgresRepository) CountPointsWithFilters(ctx context.Context, filters *dto.CountPointsRequestDTO) (int, error) {
	var count int

	row := r.db.Pool.QueryRow(
		ctx,
		countPointsWithFiltersStmt,
		filters.Region,
		filters.DateTimeFrom, filters.DateTimeTo,
		filters.XMin, filters.Ymin,
		filters.XMax, filters.YMax,
		filters.ServingCompany,
		filters.Efficiency,
		filters.GradeForService,
		filters.UrgencyCategory,
		filters.WorkType,
		filters.DeffectCategory,
		filters.OwnerCompany,
	)

	if err := row.Scan(&count); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("failed to count points: %w", err)
	}

	return count, nil
}

const getPointsWithFiltersStmt = `
    SELECT
        m.root_id  AS root_id,
        a.center_x AS x,
        a.center_y AS y
    FROM addresses AS a
    INNER JOIN requests AS m on
        m.adress_unom = a.unom
    WHERE
        m.hood = $1 AND
        m.closure_date BETWEEN COALESCE($2, m.closure_date) AND COALESCE($3, m.closure_date) AND
        rectangle_contains($4, $5, $6, $7, a.center_x, a.center_y) AND
        (array_length($8::varchar[], 1)  IS NULL OR ARRAY[m.serving_company]       && ($8::varchar[])) AND
        (array_length($9::varchar[], 1)  IS NULL OR ARRAY[m.efficiency]            &&  ($9::varchar[])) AND
        (array_length($10::varchar[], 1) IS NULL OR ARRAY[m.grade_for_service]     && ($10::varchar[])) AND
        (array_length($11::varchar[], 1) IS NULL OR ARRAY[m.urgency_category]      && ($11::varchar[])) AND
        (array_length($12::varchar[], 1) IS NULL OR ARRAY[m.work_type_done]        && ($12::varchar[])) AND
        (array_length($13::varchar[], 1) IS NULL OR ARRAY[m.deffect_category_name] && ($13::varchar[])) AND
        (array_length($14::varchar[], 1) IS NULL OR ARRAY[m.owner_company]         && ($14::varchar[]))
    LIMIT $15
    OFFSET COALESCE($16, 0);
`

func (r *hcsPostgresRepository) GetPointsWithFilters(
	ctx context.Context,
	filters *dto.GetPointsRequestQueryDTO,
) ([]entity.HCSPoint, error) {
	points := make([]entity.HCSPoint, 0, filters.Limit)

	rows, err := r.db.Pool.Query(
		ctx,
		getPointsWithFiltersStmt,
		filters.Region,
		filters.DateTimeFrom, filters.DateTimeTo,
		filters.XMin, filters.Ymin,
		filters.XMax, filters.YMax,
		filters.ServingCompany,
		filters.Efficiency,
		filters.GradeForService,
		filters.UrgencyCategory,
		filters.WorkType,
		filters.DeffectCategory,
		filters.OwnerCompany,
		filters.Limit, filters.Offest,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return points, nil
		}

		return nil, fmt.Errorf("failed to get HCS points: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		point := new(entity.HCSPoint)

		if err := rows.Scan(
			&point.RootID,
			&point.X,
			&point.Y,
		); err != nil {
			return nil, fmt.Errorf("failed to scan HCS point: %w", err)
		}

		points = append(points, *point)
	}

	return points, nil
}

const getRequestsByPointsIdsStmt = `
    SELECT
        root_id,
        version_id,
        request_number,
        mos_ru_request_number,
        date_of_creation,
        date_of_start,
        name_of_source,
        name_of_source_eng,
        name_of_creator,
        incident_feature,
        root_identificator_of_maternal,
        number_of_maternal,
        last_name_redacted,
        role_of_user,
        commentaries,
        deffect_category_name,
        deffect_category_id,
        deffect_category_name_eng,
        deffect_name,
        short_deffect_name,
        deffect_id,
        code_of_deffect,
        need_for_revision,
        description,
        presence_of_question,
        urgency_category,
        urgency_category_eng,
        district,
        district_code,
        hood,
        hood_code,
        adress_of_problem,
        adress_unom,
        porch,
        floor,
        flat_number,
        dispetchers_number,
        owner_company,
        serving_company,
        performing_company,
        inn_of_performing_company,
        request_status,
        request_status_eng,
        reason_for_decline,
        id_of_reason_for_decline,
        reason_for_decline_of_org,
        id_of_reason_for_decline_of_org,
        work_type_done,
        id_work_type_done,
        used_material,
        guarding_events,
        id_guarding_events,
        efficiency,
        efficiency_eng,
        times_returned,
        date_of_last_return_for_revision,
        being_on_revision,
        alerted_feature,
        closure_date,
        wanted_time_from,
        wanted_time_to,
        date_of_review,
        review,
        grade_for_service,
        grade_for_service_eng,
        payment_category,
        payment_category_eng,
        payed_by_card,
        date_of_previous_request_close
    FROM
        requests
    WHERE
        root_id = ANY($1);
`

// nolint // duplicates ignoring
func (r *hcsPostgresRepository) GetRequestsByPointsIds(ctx context.Context, ids []string) ([]*entity.Request, error) {
	requests := make([]*entity.Request, 0, len(ids))

	rows, err := r.db.Pool.Query(
		ctx,
		getRequestsByPointsIdsStmt,
		ids,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return requests, nil
		}

		return nil, fmt.Errorf("failed to get HCS points: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		req := new(entity.Request)
		if err := rows.Scan(
			&req.RootID,
			&req.VersionID,
			&req.RequestNumber,
			&req.MosRuRequestNumber,
			&req.DateOfCreation,
			&req.DateOfStart,
			&req.NameOfSource,
			&req.NameOsSourceENG,
			&req.NameOfCreator,
			&req.IncidentFeature,
			&req.RootIdentifcatorOfMaternal,
			&req.NumberOfMaternal,
			&req.LastNameRedacted,
			&req.RoleOfUser,
			&req.Commentaries,
			&req.DeffectCategoryName,
			&req.DeffectCategoryID,
			&req.DeffectCategoryNameENG,
			&req.DeffectName,
			&req.ShortDeffectName,
			&req.DeffectID,
			&req.CodeOfDeffect,
			&req.NeedForRevision,
			&req.Description,
			&req.PresenceOfQuestion,
			&req.UrgencyCategory,
			&req.UrgencyCategoryENG,
			&req.District,
			&req.DistrictCode,
			&req.Hood,
			&req.HoodCode,
			&req.AddressOfProblem,
			&req.AddressUNOM,
			&req.Porch,
			&req.Floor,
			&req.FlatNumber,
			&req.DispetchersNumber,
			&req.OwnerCompany,
			&req.ServingCompany,
			&req.PerformingCompany,
			&req.InnOfPerformingCompany,
			&req.RequestStatus,
			&req.RequestStatusENG,
			&req.ReasonForDecline,
			&req.ReasonForDeclineID,
			&req.ReasonForDeclineOfOrg,
			&req.ReasonForDeclineOfOrgID,
			&req.WorkTypeDone,
			&req.IDWorkTypeDone,
			&req.UsedMaterial,
			&req.GuardingEvents,
			&req.GuardingEventsID,
			&req.Effeciency,
			&req.EffeciencyENG,
			&req.TimesReturnes,
			&req.DateOfLastReturnForRevision,
			&req.BeingOnRevision,
			&req.AlertedFeature,
			&req.ClosureDate,
			&req.WantedTimeFrom,
			&req.WatnedTimeTo,
			&req.DateOfReview,
			&req.Review,
			&req.GradeForService,
			&req.GradeForServiceENG,
			&req.PaymentCategory,
			&req.PaymentCategoryENG,
			&req.PayedByCard,
			&req.DateOfPreviousRequestClose,
		); err != nil {
			return nil, fmt.Errorf("failed to scan request: %w", err)
		}

		requests = append(requests, req)
	}

	return requests, nil
}

const getPreviousRequestDatetimeStmt = `
WITH requests AS MATERIALIZED (
	SELECT * FROM requests
	WHERE
		adress_unom  = $1 AND
		"floor"      = $2 AND
		flat_number  = $3 AND
		deffect_id   = $4 AND
		closure_date < $5
    ORDER BY closure_date
)

SELECT MAX(closure_date) FROM requests;
`

func (r *hcsPostgresRepository) GetPrimaryRequestDatetime(ctx context.Context, current *entity.Request) (*time.Time, error) {
	prDateTime := new(time.Time)

	row := r.db.Pool.QueryRow(
		ctx,
		getPreviousRequestDatetimeStmt,
		current.AddressUNOM,
		current.Floor,
		current.FlatNumber,
		current.DeffectID,
		current.DateOfCreation,
	)

	if err := row.Scan(&prDateTime); err != nil {
		return nil, fmt.Errorf("failed to get primary request datetime: %w", err)
	}

	return prDateTime, nil
}

const countRequestsByClosureTimeStmt = `
    SELECT
        count(*)
    FROM
        requests
    WHERE
        closure_date IS NOT NULL AND
        closure_date BETWEEN $1 AND $2;
`

func (r *hcsPostgresRepository) CountRequestsByClosureTime(ctx context.Context, from, to time.Time) (int, error) {
	var count int

	row := r.db.Pool.QueryRow(ctx, countRequestsByClosureTimeStmt, from, to)

	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count requests by closure time: %w", err)
	}

	return count, nil
}

const getRequestsByClosureTimeStmt = `
    SELECT
        root_id,
        version_id,
        request_number,
        mos_ru_request_number,
        date_of_creation,
        date_of_start,
        name_of_source,
        name_of_source_eng,
        name_of_creator,
        incident_feature,
        root_identificator_of_maternal,
        number_of_maternal,
        last_name_redacted,
        role_of_user,
        commentaries,
        deffect_category_name,
        deffect_category_id,
        deffect_category_name_eng,
        deffect_name,
        short_deffect_name,
        deffect_id,
        code_of_deffect,
        need_for_revision,
        description,
        presence_of_question,
        urgency_category,
        urgency_category_eng,
        district,
        district_code,
        hood,
        hood_code,
        adress_of_problem,
        adress_unom,
        porch,
        floor,
        flat_number,
        dispetchers_number,
        owner_company,
        serving_company,
        performing_company,
        inn_of_performing_company,
        request_status,
        request_status_eng,
        reason_for_decline,
        id_of_reason_for_decline,
        reason_for_decline_of_org,
        id_of_reason_for_decline_of_org,
        work_type_done,
        id_work_type_done,
        used_material,
        guarding_events,
        id_guarding_events,
        efficiency,
        efficiency_eng,
        times_returned,
        date_of_last_return_for_revision,
        being_on_revision,
        alerted_feature,
        closure_date,
        wanted_time_from,
        wanted_time_to,
        date_of_review,
        review,
        grade_for_service,
        grade_for_service_eng,
        payment_category,
        payment_category_eng,
        payed_by_card,
        date_of_previous_request_close
    FROM
        requests
    WHERE
        closure_date IS NOT NULL AND
        closure_date BETWEEN $1 AND $2
    ORDER BY root_id
    LIMIT $3
    OFFSET COALESCE($4, 0);
`

// nolint // duplicates ignoring
func (r *hcsPostgresRepository) GetRequestsByClosureTime(
	ctx context.Context,
	from, to time.Time,
	limit, offset int,
) ([]*entity.Request, error) {
	requests := make([]*entity.Request, 0, limit)

	rows, err := r.db.Pool.Query(ctx, getRequestsByClosureTimeStmt, from, to, limit, offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return requests, nil
		}

		return nil, fmt.Errorf("failed to get requests by closure time: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		req := new(entity.Request)

		if err := rows.Scan(
			&req.RootID,
			&req.VersionID,
			&req.RequestNumber,
			&req.MosRuRequestNumber,
			&req.DateOfCreation,
			&req.DateOfStart,
			&req.NameOfSource,
			&req.NameOsSourceENG,
			&req.NameOfCreator,
			&req.IncidentFeature,
			&req.RootIdentifcatorOfMaternal,
			&req.NumberOfMaternal,
			&req.LastNameRedacted,
			&req.RoleOfUser,
			&req.Commentaries,
			&req.DeffectCategoryName,
			&req.DeffectCategoryID,
			&req.DeffectCategoryNameENG,
			&req.DeffectName,
			&req.ShortDeffectName,
			&req.DeffectID,
			&req.CodeOfDeffect,
			&req.NeedForRevision,
			&req.Description,
			&req.PresenceOfQuestion,
			&req.UrgencyCategory,
			&req.UrgencyCategoryENG,
			&req.District,
			&req.DistrictCode,
			&req.Hood,
			&req.HoodCode,
			&req.AddressOfProblem,
			&req.AddressUNOM,
			&req.Porch,
			&req.Floor,
			&req.FlatNumber,
			&req.DispetchersNumber,
			&req.OwnerCompany,
			&req.ServingCompany,
			&req.PerformingCompany,
			&req.InnOfPerformingCompany,
			&req.RequestStatus,
			&req.RequestStatusENG,
			&req.ReasonForDecline,
			&req.ReasonForDeclineID,
			&req.ReasonForDeclineOfOrg,
			&req.ReasonForDeclineOfOrgID,
			&req.WorkTypeDone,
			&req.IDWorkTypeDone,
			&req.UsedMaterial,
			&req.GuardingEvents,
			&req.GuardingEventsID,
			&req.Effeciency,
			&req.EffeciencyENG,
			&req.TimesReturnes,
			&req.DateOfLastReturnForRevision,
			&req.BeingOnRevision,
			&req.AlertedFeature,
			&req.ClosureDate,
			&req.WantedTimeFrom,
			&req.WatnedTimeTo,
			&req.DateOfReview,
			&req.Review,
			&req.GradeForService,
			&req.GradeForServiceENG,
			&req.PaymentCategory,
			&req.PaymentCategoryENG,
			&req.PayedByCard,
			&req.DateOfPreviousRequestClose,
		); err != nil {
			return nil, fmt.Errorf("failed to scan request: %w", err)
		}

		requests = append(requests, req)
	}

	return requests, nil
}

const countRequestsFullStmt = `
    SELECT
        count(*)
    FROM
        requests as mt
    INNER JOIN addresses AS a ON
        a.unom = mt.adress_unom
    INNER JOIN requests_anomalies AS ra ON
    	ra.root_id = mt.root_id
    WHERE
        mt.hood = $1 AND
        mt.closure_date BETWEEN COALESCE($2, mt.closure_date) AND COALESCE($3, mt.closure_date) AND
        rectangle_contains($4, $5, $6, $7, a.center_x, a.center_y) AND
        mt.urgency_category = COALESCE($8, mt.urgency_category) AND
        (array_length($9::int[], 1) IS NULL OR ra.anomaly_cases && ($9));
`

func (r *hcsPostgresRepository) CountRequestsFull(ctx context.Context, filters *dto.CountRequestsFullRequestDTO) (int, error) {
	var count int

	row := r.db.Pool.QueryRow(
		ctx,
		countRequestsFullStmt,
		filters.Region,
		filters.DateTimeFrom,
		filters.DateTimeTo,
		filters.XMin,
		filters.Ymin,
		filters.XMax,
		filters.YMax,
		filters.UrgencyCategory,
		filters.AnomalyCases,
	)

	if err := row.Scan(&count); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, fmt.Errorf("failed to count points: %w", err)
	}

	return count, nil
}

const getRequestsFullStmt = `
    SELECT
        mt.root_id,
        mt.version_id,
        mt.request_number,
        mt.mos_ru_request_number,
        mt.date_of_creation,
        mt.date_of_start,
        mt.name_of_source,
        mt.name_of_source_eng,
        mt.name_of_creator,
        mt.incident_feature,
        mt.root_identificator_of_maternal,
        mt.number_of_maternal,
        mt.last_name_redacted,
        mt.role_of_user,
        mt.commentaries,
        mt.deffect_category_name,
        mt.deffect_category_id,
        mt.deffect_category_name_eng,
        mt.deffect_name,
        mt.short_deffect_name,
        mt.deffect_id,
        mt.code_of_deffect,
        mt.need_for_revision,
        mt.description,
        mt.presence_of_question,
        mt.urgency_category,
        mt.urgency_category_eng,
        mt.district,
        mt.district_code,
        mt.hood,
        mt.hood_code,
        mt.adress_of_problem,
        mt.adress_unom,
        mt.porch,
        mt.floor,
        mt.flat_number,
        mt.dispetchers_number,
        mt.owner_company,
        mt.serving_company,
        mt.performing_company,
        mt.inn_of_performing_company,
        mt.request_status,
        mt.request_status_eng,
        mt.reason_for_decline,
        mt.id_of_reason_for_decline,
        mt.reason_for_decline_of_org,
        mt.id_of_reason_for_decline_of_org,
        mt.work_type_done,
        mt.id_work_type_done,
        mt.used_material,
        mt.guarding_events,
        mt.id_guarding_events,
        mt.efficiency,
        mt.efficiency_eng,
        mt.times_returned,
        mt.date_of_last_return_for_revision,
        mt.being_on_revision,
        mt.alerted_feature,
        mt.closure_date,
        mt.wanted_time_from,
        mt.wanted_time_to,
        mt.date_of_review,
        mt.review,
        mt.grade_for_service,
        mt.grade_for_service_eng,
        mt.payment_category,
        mt.payment_category_eng,
        mt.payed_by_card,
        mt.date_of_previous_request_close,
        a.center_x,
        a.center_y,
        ra.is_anomaly,
        ra.anomaly_cases
    FROM
        requests as mt
    INNER JOIN addresses AS a ON
        a.unom = mt.adress_unom
    INNER JOIN requests_anomalies AS ra ON
    	ra.root_id = mt.root_id
    WHERE
        mt.hood = $1 AND
        mt.closure_date BETWEEN COALESCE($2, mt.closure_date) AND COALESCE($3, mt.closure_date) AND
        rectangle_contains($4, $5, $6, $7, a.center_x, a.center_y) AND
        mt.urgency_category = COALESCE($8, mt.urgency_category) AND
        (array_length($9::int[], 1) IS NULL OR ra.anomaly_cases && ($9))
    LIMIT $10
    OFFSET COALESCE($11, 0);
`

// nolint // length ignoring
func (r *hcsPostgresRepository) GetRequestsFull(ctx context.Context, filters *dto.GetRequestsFullRequestDTO) ([]*entity.RequestFull, error) {
	requests := make([]*entity.RequestFull, 0, filters.Limit)

	rows, err := r.db.Pool.Query(
		ctx,
		getRequestsFullStmt,
		filters.Region,
		filters.DateTimeFrom,
		filters.DateTimeTo,
		filters.XMin,
		filters.Ymin,
		filters.XMax,
		filters.YMax,
		filters.UrgencyCategory,
		filters.AnomalyCases,
		filters.Limit,
		filters.Offest,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return requests, nil
		}

		return nil, fmt.Errorf("failed to get requests by closure time: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		full := new(entity.RequestFull)

		req := new(entity.Request)
		p := new(entity.HCSPoint)
		ra := new(entity.RequestAnomaly)

		if err := rows.Scan(
			&full.RootID,
			&req.VersionID,
			&req.RequestNumber,
			&req.MosRuRequestNumber,
			&req.DateOfCreation,
			&req.DateOfStart,
			&req.NameOfSource,
			&req.NameOsSourceENG,
			&req.NameOfCreator,
			&req.IncidentFeature,
			&req.RootIdentifcatorOfMaternal,
			&req.NumberOfMaternal,
			&req.LastNameRedacted,
			&req.RoleOfUser,
			&req.Commentaries,
			&req.DeffectCategoryName,
			&req.DeffectCategoryID,
			&req.DeffectCategoryNameENG,
			&req.DeffectName,
			&req.ShortDeffectName,
			&req.DeffectID,
			&req.CodeOfDeffect,
			&req.NeedForRevision,
			&req.Description,
			&req.PresenceOfQuestion,
			&req.UrgencyCategory,
			&req.UrgencyCategoryENG,
			&req.District,
			&req.DistrictCode,
			&req.Hood,
			&req.HoodCode,
			&req.AddressOfProblem,
			&req.AddressUNOM,
			&req.Porch,
			&req.Floor,
			&req.FlatNumber,
			&req.DispetchersNumber,
			&req.OwnerCompany,
			&req.ServingCompany,
			&req.PerformingCompany,
			&req.InnOfPerformingCompany,
			&req.RequestStatus,
			&req.RequestStatusENG,
			&req.ReasonForDecline,
			&req.ReasonForDeclineID,
			&req.ReasonForDeclineOfOrg,
			&req.ReasonForDeclineOfOrgID,
			&req.WorkTypeDone,
			&req.IDWorkTypeDone,
			&req.UsedMaterial,
			&req.GuardingEvents,
			&req.GuardingEventsID,
			&req.Effeciency,
			&req.EffeciencyENG,
			&req.TimesReturnes,
			&req.DateOfLastReturnForRevision,
			&req.BeingOnRevision,
			&req.AlertedFeature,
			&req.ClosureDate,
			&req.WantedTimeFrom,
			&req.WatnedTimeTo,
			&req.DateOfReview,
			&req.Review,
			&req.GradeForService,
			&req.GradeForServiceENG,
			&req.PaymentCategory,
			&req.PaymentCategoryENG,
			&req.PayedByCard,
			&req.DateOfPreviousRequestClose,
			&p.X,
			&p.Y,
			&ra.IsAnomaly,
			&ra.AnomalyCases,
		); err != nil {
			return nil, fmt.Errorf("failed to scan request: %w", err)
		}

		full.Request = req
		full.HCSPoint = p
		full.RequestAnomaly = ra

		requests = append(requests, full)
	}

	return requests, nil
}

const maxHoods = 125

func (r *hcsPostgresRepository) GetUniqueRegions(ctx context.Context) ([]string, error) {
	hoods := make([]string, 0, maxHoods)

	rows, err := r.db.Pool.Query(
		ctx,
		`SELECT DISTINCT(hood) FROM requests WHERE hood IS NOT NULL AND hood != '';`,
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

func (r *hcsPostgresRepository) GetUniqueServingCompanies(ctx context.Context) ([]string, error) {
	scs := make([]string, 0, maxHoods)

	rows, err := r.db.Pool.Query(
		ctx,
		`SELECT DISTINCT(serving_company) FROM requests WHERE serving_company IS NOT NULL AND serving_company != '';`,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return scs, nil
		}

		return nil, fmt.Errorf("failed to get unqiue serving_company: %w", err)
	}

	defer rows.Close()

	var sc string

	for rows.Next() {
		if err = rows.Scan(&sc); err != nil {
			return nil, fmt.Errorf("failed to scan serving_company: %w", err)
		}

		scs = append(scs, sc)
	}

	return scs, nil
}

func (r *hcsPostgresRepository) GetUniqueOwnerCompanies(ctx context.Context) ([]string, error) {
	scs := make([]string, 0, maxHoods)

	rows, err := r.db.Pool.Query(
		ctx,
		`SELECT DISTINCT(owner_company) FROM requests WHERE owner_company IS NOT NULL AND owner_company != '';`,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return scs, nil
		}

		return nil, fmt.Errorf("failed to get unqiue owner_company: %w", err)
	}

	defer rows.Close()

	var sc string

	for rows.Next() {
		if err = rows.Scan(&sc); err != nil {
			return nil, fmt.Errorf("failed to scan owner_company: %w", err)
		}

		scs = append(scs, sc)
	}

	return scs, nil
}

func (r *hcsPostgresRepository) GetUniqueDeffectCategories(ctx context.Context) ([]string, error) {
	scs := make([]string, 0, maxHoods)

	rows, err := r.db.Pool.Query(
		ctx,
		`SELECT DISTINCT(deffect_category_name) FROM requests WHERE deffect_category_name IS NOT NULL AND deffect_category_name != '';`,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return scs, nil
		}

		return nil, fmt.Errorf("failed to get unqiue deffect_category_name: %w", err)
	}

	defer rows.Close()

	var sc string

	for rows.Next() {
		if err = rows.Scan(&sc); err != nil {
			return nil, fmt.Errorf("failed to scan deffect_category_name: %w", err)
		}

		scs = append(scs, sc)
	}

	return scs, nil
}

func (r *hcsPostgresRepository) GetUniqueWorkTypes(ctx context.Context) ([]string, error) {
	scs := make([]string, 0, maxHoods)

	rows, err := r.db.Pool.Query(
		ctx,
		`SELECT DISTINCT(work_type_done) FROM requests WHERE work_type_done IS NOT NULL AND work_type_done != '';`,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return scs, nil
		}

		return nil, fmt.Errorf("failed to get unqiue work_type_done: %w", err)
	}

	defer rows.Close()

	var sc string

	for rows.Next() {
		if err = rows.Scan(&sc); err != nil {
			return nil, fmt.Errorf("failed to scan work_type_done: %w", err)
		}

		scs = append(scs, sc)
	}

	return scs, nil
}

const getRequestsByDispatcherStmt = `
SELECT
	root_id 	  	   			  AS root_id,
	dispetchers_number 			  AS dispetchers_number,
	efficiency 		   			  AS efficiency,
	closure_date 	   			  AS closure_date,
	NULLIF(grade_for_service, '') AS grade_for_service
FROM requests
WHERE
	closure_date BETWEEN $1 AND $2 AND
	dispetchers_number = $3
;
`

func (r *hcsPostgresRepository) GetRequestsByDispatcher(
	ctx context.Context,
	dateFrom, dateTo time.Time,
	dispNumber string,
) ([]*entity.Request, error) {
	requests := make([]*entity.Request, 0)

	rows, err := r.db.Pool.Query(ctx, getRequestsByDispatcherStmt, dateFrom, dateTo, dispNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return requests, nil
		}

		return nil, fmt.Errorf("failed to get requests by dispatcher: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		req := new(entity.Request)

		if err := rows.Scan(
			&req.RootID,
			&req.DispetchersNumber,
			&req.Effeciency,
			&req.ClosureDate,
			&req.GradeForService,
		); err != nil {
			return nil, fmt.Errorf("failed to scan request: %w", err)
		}

		requests = append(requests, req)
	}

	return requests, nil
}

func (r *hcsPostgresRepository) GetUniqueDispatchers(ctx context.Context) ([]string, error) {
	scs := make([]string, 0, maxHoods)

	rows, err := r.db.Pool.Query(
		ctx,
		`SELECT DISTINCT(dispetchers_number) FROM requests
        WHERE
            dispetchers_number IS NOT NULL AND
            dispetchers_number != '' AND
            dispetchers_number != 'тест мос. ру'
        ORDER BY dispetchers_number desc;`,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return scs, nil
		}

		return nil, fmt.Errorf("failed to get unqiue dispetchers_number: %w", err)
	}

	defer rows.Close()

	var sc string

	for rows.Next() {
		if err = rows.Scan(&sc); err != nil {
			return nil, fmt.Errorf("failed to scan dispetchers_number: %w", err)
		}

		scs = append(scs, sc)
	}

	return scs, nil
}

func (r *hcsPostgresRepository) GetRegionArea(ctx context.Context, region string) (string, error) {
	var areaGeojson string

	row := r.db.Pool.QueryRow(
		ctx,
		`SELECT ST_AsGeoJSON(wkt) FROM hoods_areas WHERE hood = $1;`,
		region,
	)

	if err := row.Scan(&areaGeojson); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("failed to count points: %w", err)
	}

	return areaGeojson, nil
}
