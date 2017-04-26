-- GET USER CONTRACTS, CONTACTS, AND FACILITIES AS ID ARRAY

SELECT * FROM (
		SELECT uf.facility_id as id, uf.id as user_facilities_id, json_agg(ci) as contacts
		FROM user_facilities uf
		LEFT OUTER JOIN contact_instance ci
		ON ci.user_facilities_id=uf.id
		WHERE uf.user_id=$1
		GROUP BY uf.facility_id, uf.id
) uf
LEFT JOIN facilities c
USING (id)
