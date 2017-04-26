SELECT * FROM (
	SELECT people_id AS id
	FROM user_contacts uc
	WHERE uc.user_id=$1
) u
INNER JOIN people c
USING (id)
LEFT OUTER JOIN (
	SELECT id "provider_id"
	, id "provider.id"
	, npi "provider.npi"
	, w9 "provider.w9"
	, direct_deposit_form "provider.direct_deposit_form"
	, hourly_rate "provider.hourly_rate"
	, desired_shifts_month "provider.desired_shifts_month"
	, max_shifts_month "provider.max_shifts_month"
	, min_shifts_month "provider.min_shifts_month"
	, full_time "provider.full_time"
	, part_time "provider.part_time"
	, prn "provider.prn"
	, retired "provider.retired"
	, notes "provider.notes"
	, insurance_certificate "provider.insurance_certificate"
	, tb_expiration "provider.tb_expiration"
	, tb_file "provider.tb_file"
	, flu_expiration "provider.flu_expiration"
	, flu_file "provider.flu_file" FROM providers
) as provider
USING (provider_id)
LEFT OUTER JOIN (
	SELECT id "state.id", id "state_id", name "state.name", abbreviation as "state.abbr"
	FROM states
) st
USING (state_id)
LEFT OUTER JOIN (
	SELECT id "recruiter_id", id "recruiter.id"
	FROM recruiters
) rr
USING (recruiter_id)
