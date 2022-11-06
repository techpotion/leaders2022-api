UPDATE requests
	SET date_of_previous_request_close = prev_mue.closure_date
FROM requests AS prev_mue
	WHERE requests.number_of_maternal = prev_mue.request_number

-----------------------------------------

UPDATE requests mt
SET date_of_previous_request_close = (
	WITH requests AS (
	  SELECT prev_mue.prev_mue FROM requests AS prev_mue
	  WHERE
        prev_mue.request_number     != COALESCE(mt.number_of_maternal, '') AND
	    prev_mue.adress_unom         = mt.adress_unom        AND
	    prev_mue."floor"             = mt."floor"            AND
	    prev_mue.flat_number         = mt.flat_number        AND
	    prev_mue.deffect_id          = mt.deffect_id         AND
	    prev_mue.closure_date        < mt.date_of_creation
	)

	select MAX(req.closure_date) FROM requests AS req
);

---------------------------------------

WITH tmp AS (
	SELECT
		root_q.root_id,
		MAX(prev_q.closure_date) as max_date_of_previous_request_close
	FROM requests AS root_q
	LEFT JOIN requests AS prev_q ON
		prev_q.request_number != COALESCE(root_q.number_of_maternal, '') AND
		prev_q.adress_unom     = root_q.adress_unom AND
		prev_q."floor"         = root_q."floor" AND
		prev_q.flat_number     = root_q.flat_number AND
		prev_q.deffect_id      = root_q.deffect_id AND
		prev_q.closure_date    < root_q.date_of_creation
	GROUP BY
		root_q.root_id
)

UPDATE
	requests
SET
	date_of_previous_request_close = tmp.max_date_of_previous_request_close
FROM
	tmp
WHERE
	requests.root_id = tmp.root_id
