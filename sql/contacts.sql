
-- GET USER CONTACTS
SELECT * FROM (
			SELECT people_id AS id
			FROM user_contacts uc
			WHERE uc.user_id=$1
		) u
INNER JOIN people c
USING (id)

-- GET USER CONTACT BY ID
SELECT * FROM (
			SELECT people_id AS id
			FROM user_contacts uc
			WHERE uc.user_id=$1
		) u
INNER JOIN people c
USING (id)
WHERE u.id=$1
