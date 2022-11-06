CREATE TABLE IF NOT EXISTS requests_anomalies (
    root_id       varchar   NOT NULL UNIQUE,
    is_anomaly    boolean   NOT NULL,
    anomaly_cases integer[] NOT NULL
);

