-- GET ALL USER FACILITIES
SELECT * FROM (
	SELECT facility_id as id
	FROM user_facilities
	WHERE user_id=1
) uf
LEFT OUTER JOIN facilities f
USING (id)
