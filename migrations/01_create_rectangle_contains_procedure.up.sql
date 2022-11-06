CREATE OR REPLACE FUNCTION rectangle_contains(
	x_min float,
	y_min float,
	x_max float,
	y_max float,
	x float,
	y float
)
RETURNS boolean
LANGUAGE plpgsql
AS $function$
BEGIN
	RETURN (x BETWEEN x_min AND x_max) AND (y BETWEEN y_min AND y_max);
END;
$function$
