CREATE TABLE IF NOT EXISTS requests
(
    root_id                          varchar,
    version_id                       varchar,
    request_number                   varchar,
    mos_ru_request_number            varchar,
    date_of_creatiON                 timestamp,
    date_of_start                    timestamp,
    name_of_source                   varchar,
    name_of_source_eng               varchar,
    name_of_creator                  varchar,
    incident_feature                 varchar,
    root_identificator_of_maternal   varchar,
    number_of_maternal               varchar,
    last_name_redacted               varchar,
    role_of_user                     varchar,
    commentaries                     varchar,
    deffect_category_name            varchar,
    deffect_category_id              varchar,
    deffect_category_name_eng        varchar,
    deffect_name                     varchar,
    short_deffect_name               varchar,
    deffect_id                       varchar,
    code_of_deffect                  varchar,
    need_for_revisiON                varchar,
    descriptiON                      varchar,
    presence_of_questiON             varchar,
    urgency_category                 varchar,
    urgency_category_eng             varchar,
    district                         varchar,
    district_code                    varchar,
    hood                             varchar,
    hood_code                        varchar,
    adress_of_problem                varchar,
    adress_unom                      integer,
    porch                            varchar,
    floor                            varchar,
    flat_number                      varchar,
    dispetchers_number               varchar,
    owner_company                    varchar,
    serving_company                  varchar,
    performing_company               varchar,
    inn_of_performing_company        varchar,
    request_status                   varchar,
    request_status_eng               varchar,
    reason_for_decline               varchar,
    id_of_reason_for_decline         varchar,
    reason_for_decline_of_org        varchar,
    id_of_reason_for_decline_of_org  varchar,
    work_type_done                   varchar,
    id_work_type_done                varchar,
    used_material                    varchar,
    guarding_events                  varchar,
    id_guarding_events               varchar,
    efficiency                       varchar,
    efficiency_eng                   varchar,
    times_returned                   integer,
    date_of_last_return_for_revisiON timestamp,
    being_on_revisiON                varchar,
    alerted_feature                  varchar,
    closure_date                     timestamp,
    wanted_time_from                 timestamp,
    wanted_time_to                   timestamp,
    date_of_review                   timestamp,
    review                           varchar,
    grade_for_service                varchar,
    grade_for_service_eng            varchar,
    payment_category                 varchar,
    payment_category_eng             varchar,
    payed_by_card                    varchar,
    date_of_previous_request_close   timestamp
);

CREATE INDEX IF NOT EXISTS requests_addr_idx
    ON requests (adress_of_problem);

CREATE INDEX IF NOT EXISTS requests_date_of_creation_idx
    ON requests (date_of_creation);

CREATE INDEX IF NOT EXISTS requests_date_of_start_idx
    ON requests (date_of_start);

create UNIQUE INDEX IF NOT EXISTS requests_root_id_uidx
    ON requests (root_id);

CREATE INDEX IF NOT EXISTS requests_unom_idx
    ON requests (adress_unom);

CREATE INDEX IF NOT EXISTS requests_hood_idx
    ON requests (hood);

CREATE INDEX IF NOT EXISTS requests_closure_date_idx
    ON requests (closure_date);

CREATE INDEX IF NOT EXISTS requests_floor_idx
    ON requests (floor);

CREATE INDEX IF NOT EXISTS requests_flat_number_idx
    ON requests (flat_number);

CREATE INDEX IF NOT EXISTS requests_deffect_id_idx
    ON requests (deffect_id);

CREATE UNIQUE INDEX IF NOT EXISTS requests_request_number_uidx
    ON requests (request_number);

CREATE INDEX IF NOT EXISTS requests_number_of_maternal_idx
    ON requests (number_of_maternal);

CREATE INDEX IF NOT EXISTS requests_date_of_previous_request_close_idx
    ON requests (date_of_previous_request_close);
