-- ALTER TABLE IF EXISTS requests
--     ADD IF NOT EXISTS date_of_previous_request_close TIMESTAMP;
--
CREATE OR REPLACE FUNCTION update_date_of_previous_request_close() RETURNS trigger AS $$
BEGIN
    IF tg_op = 'INSERT' OR tg_op = 'UPDATE' THEN
		new.date_of_previous_request_close = (
			WITH requests AS MATERIALIZED (
			  SELECT closure_date FROM requests
			  WHERE
			  	request_number != COALESCE(new.number_of_maternal, '') AND
			    adress_unom = new.adress_unom AND
			    "floor" = new."floor" AND
			    flat_number = new.flat_number AND
			    deffect_id = new.deffect_id AND
			    closure_date < new.date_of_creation
			)

			select MAX(requests.closure_date) from requests
		);
RETURN new;
END IF;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_date_of_previous_request_close on requests;
CREATE TRIGGER update_date_of_previous_request_close BEFORE INSERT OR UPDATE ON
        requests
    FOR EACH ROW EXECUTE
        PROCEDURE update_date_of_previous_request_close();
