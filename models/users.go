package Models

import (
	"time"

	null "gopkg.in/guregu/null.v3"
)

//User move this later
type User struct {
	ID                    string      `json:"id" `
	UpdatedAt             time.Time   `json:"updated_at" `
	Email                 string      `json:"email" `
	Password              string      `json:"password"`
	PasswordConfirm       string      `json:"password_confirm" `
	PeopleID              null.String `json:"people_id" `
	Token                 null.String `json:"token"`
	IntercomHash          null.String `json:"intercom_hash"`
	EmailForwardingHandle null.String `json:"email_forwarding_handle"`
	Plan                  null.String `json:"plan"`
	PlanStart             null.Time   `json:"plan_start"`
	PlanRenew             null.Time   `json:"plan_renew"`
	Trial                 bool        `json:"trial"`
	TrialExpiration       null.Time   `json:"trial_expiration"`
	OnboardingID          null.String `json:"onboarding_id"`
	AccountType           null.String `json:"account_type"`
	AccountTypeID         null.String `json:"account_type_id"`
}

const (
	//CreateUser sql statement to create a user
	CreateUser string = `
		WITH new_onboarding as (
			INSERT INTO onboarding(completed)
			VALUES (false)
			RETURNING id
		)
		INSERT INTO users
			( email
			, password
			, onboarding_id
		 )
		VALUES (
			$1,
			$2,
			(SELECT id FROM new_onboarding)
		)
		RETURNING *
	`

	//UpdateUser sql statement to update a user
	UpdateUser string = `
		UPDATE users
		SET password=$2
		WHERE id=$1
		RETURNING *
	`

	//GetUserByEmail SQL Statement to retrieve user by email
	GetUserByEmail string = `
		SELECT u.id,
		at.name as account_type,
		u.email,
		u.email_forwarding_handle,
		u.plan,
		u.plan_start,
		u.plan_renew,
		u.trial,
		u.trial_expiration,
		u.onboarding_id,
		u.password
		FROM users u, account_types at
		WHERE email=$1
		AND u.account_type_id=at.id
	`
)
