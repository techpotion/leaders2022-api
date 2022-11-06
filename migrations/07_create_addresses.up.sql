CREATE TABLE IF NOT EXISTS addresses
(
    unom               bigint,
    address            text,
    pos_geojson        text,
    pos_geojson_center text,
    center_gis         geometry,
    center_x           double precision,
    center_y           double precision
);

CREATE INDEX IF NOT EXISTS addr_idx
    ON addresses (address);

CREATE UNIQUE INDEX IF NOT EXISTS addr_unom_idx
    ON addresses (unom);

CREATE INDEX IF NOT EXISTS addr_x
    ON addresses (center_x);

CREATE INDEX IF NOT EXISTS addr_y
    ON addresses (center_y);
