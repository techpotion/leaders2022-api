CREATE TABLE IF NOT EXISTS custom_requests_anomalies (
    root_id    varchar NOT NULL UNIQUE,
    is_anomaly boolean NOT NULL DEFAULT false
);


