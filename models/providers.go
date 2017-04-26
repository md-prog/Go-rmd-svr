package Models

import (
	"time"

	null "gopkg.in/guregu/null.v3"
	"gopkg.in/guregu/null.v3/zero"
)

//Provider move this later
type Provider struct {
	ID                    null.String `json:"id" form:"id" query:"id" db:"provider_id"`
	UpdatedAt             time.Time   `json:"updated_at" form:"updated_at" query:"updated_at"`
	Npi                   null.Int    `json:"npi" form:"npi" query:"npi" db:"provider_npi"`
	W9                    null.String `json:"w9" form:"w9" query:"w9"   db:"provider_w9"`
	DirectDepositForm     null.String `json:"direct_deposit_form" form:"direct_deposit_form" query:"direct_deposit_form"  db:"provider_direct_deposit_form"`
	HourlyRate            null.Int    `json:"hourly_rate" form:"hourly_rate" query:"hourly_rate"  db:"provider_hourly_rate"`
	DesiredShiftsPerMonth null.Int    `json:"desired_shifts_month" form:"desired_shifts_month" query:"desired_shifts_month"  db:"provider_desired_shifts_month"`
	MaxShiftsPerMonth     null.Int    `json:"max_shifts_month" form:"max_shifts_month" query:"max_shifts_month"  db:"provider_max_shifts_month"`
	MinShiftsPerMonth     null.Int    `json:"min_shifts_month" form:"min_shifts_month" query:"min_shifts_month"  db:"provider_min_shifts_month"`
	FullTime              null.Bool   `json:"full_time" form:"full_time" query:"full_time"  db:"provider_full_time"`
	PartTime              null.Bool   `json:"part_time" form:"part_time" query:"part_time"  db:"provider_part_time"`
	Prn                   null.Bool   `json:"prn" form:"prn" query:"prn"  db:"provider_prn"`
	Retired               null.Bool   `json:"retired" form:"retired" query:"retired"  db:"provider_retired"`
	Notes                 null.String `json:"notes" form:"notes" query:"notes"  db:"provider_notes"`
	InsuranceCertificate  null.String `json:"insurance_certificate" form:"insurance_certificate" query:"insurance_certificate"  db:"provider_insurance_certificate"`
	TBExpiration          null.Time   `json:"tb_expiration" form:"tb_expiration" query:"tb_expiration"  db:"provider_tb_expiration"`
	TBFile                null.String `json:"tb_file" form:"tb_file" query:"tb_file"  db:"provider_tb_file"`
	FluExpiration         null.Time   `json:"flu_expiration" form:"flu_expiration" query:"flu_expiration"  db:"provider_flu_expiration"`
	FluFile               null.String `json:"flu_file" form:"flu_file" query:"flu_file"  db:"provider_flu_file"`
	PA                    zero.Bool   `json:"pa" form:"pa" query:"pa"  db:"provider_pa"`
	MD                    zero.Bool   `json:"md" form:"md" query:"md"  db:"provider_md"`
	DO                    zero.Bool   `json:"do" form:"do" query:"do"  db:"provider_do"`
	BoardCertified        zero.Bool   `json:"board_certified" form:"board_certified" query:"board_certified"  db:"provider_board_certified"`
	NP                    zero.Bool   `json:"np" form:"np" query:"np" db:"provider_np"`
}

//ProviderStatus move this later
type ProviderStatus struct {
	ID               string    `json:"id" form:"id" query:"id"`
	UpdatedAt        time.Time `json:"updated_at" form:"updated_at" query:"updated_at"`
	LeadNoContact    null.Time `json:"lead_no_contact" form:"lead_no_contact" query:"lead_no_contact"`
	InitialContact   null.Time `json:"initial_contact" form:"initial_contact" query:"initial_contact"`
	Interested       null.Time `json:"interested" form:"interested" query:"interested"`
	NotInterestedATM null.Time `json:"not_interested_atm" form:"not_interested_atm" query:"not_interested_atm"`
	VerbalCommitment null.Time `json:"verbal_commitment" form:"verbal_commitment" query:"verbal_commitment"`
	ContractSent     null.Time `json:"contract_sent" form:"contract_sent" query:"contract_sent"`
	ContractSigned   null.Time `json:"contract_signed" form:"contract_signed" query:"contract_signed"`
	Dead             null.Time `json:"dead" form:"dead" query:"dead"`
	DoNotUse         null.Time `json:"do_not_use" form:"do_not_use" query:"do_not_use"`
	ProviderID       uint      `json:"provider_id" form:"provider_id" query:"provider_id"`
}

//ProviderLicense move this later
type ProviderLicense struct {
	ID            string      `json:"id" form:"id" query:"id"`
	UpdatedAt     time.Time   `json:"updated_at" form:"updated_at" query:"updated_at"`
	StateID       uint        `json:"state_id" form:"state_id" query:"state_id"`
	LicenseNumber uint        `json:"license_number" form:"license_number" query:"license_number"`
	File          null.String `json:"file" form:"file" query:"file"`
	ProviderID    uint        `json:"provider_id" form:"provider_id" query:"provider_id"`
}

//ProviderCertification move this later
type ProviderCertification struct {
	ID             string      `json:"id" form:"id" query:"id"`
	UpdatedAt      time.Time   `json:"updated_at" form:"updated_at" query:"updated_at"`
	Name           string      `json:"name" form:"name" query:"name"`
	ExpirationDate time.Time   `json:"expiration_date" form:"expiration_date" query:"expiration_date"`
	LicenseNumber  uint        `json:"license_number" form:"license_number" query:"license_number"`
	File           null.String `json:"file" form:"file" query:"file"`
	ProviderID     uint        `json:"provider_id" form:"provider_id" query:"provider_id"`
}

//ProviderFacility move this later
type ProvideFacility struct {
	ID               string    `json:"id" form:"id" query:"id"`
	UpdatedAt        time.Time `json:"updated_at" form:"updated_at" query:"updated_at"`
	ProviderID       uint      `json:"provider_id" form:"provider_id" query:"provider_id"`
	FacilityID       uint      `json:"facility_id" form:"facility_id" query:"facility_id"`
	CredentialedDate null.Time `json:"credentialed_date" form:"credentialed_date" query:"credentialed_date"`
	RenewalDate      null.Time `json:"renewal_date" form:"renewal_date" query:"renewal_date"`
}

//ProviderFacility move this later
type ProvideContract struct {
	ID               string      `json:"id" form:"id" query:"id"`
	UpdatedAt        time.Time   `json:"updated_at" form:"updated_at" query:"updated_at"`
	StartDate        null.Time   `json:"start_date" form:"start_date" query:"start_date"`
	RenewalDate      null.Time   `json:"renewal_date" form:"renewal_date" query:"renewal_date"`
	ProviderID       uint        `json:"provider_id" form:"provider_id" query:"provider_id"`
	UsersFacilityID  uint        `json:"users_facility_id" form:"users_facility_id" query:"users_facility_id"`
	UsersContractsID uint        `json:"users_contracts_id" form:"users_contracts_id" query:"users_contracts_id"`
	File             null.String `json:"file" form:"file" query:"file"`
}

type Hobby struct {
	ID        string    `json:"id" form:"id" query:"id"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" query:"updated_at"`
	Name      string    `json:"name" form:"name" query:"name"`
}
