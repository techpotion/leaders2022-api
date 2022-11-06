ALTER TABLE IF EXISTS requests_anomalies
    ADD IF NOT EXISTS net_probability double precision;
