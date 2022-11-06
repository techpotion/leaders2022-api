CREATE TABLE IF NOT EXISTS last_anomaly_check_job (
    name varchar(30)  NOT NULL UNIQUE DEFAULT 'anomaly_job',
    ts timestamp      NOT NULL        DEFAULT to_timestamp(0),
    is_active boolean NOT NULL        DEFAULT false
);

INSERT INTO last_anomaly_check_job VALUES('anomaly_job', to_timestamp(0), false) ON CONFLICT DO NOTHING;
