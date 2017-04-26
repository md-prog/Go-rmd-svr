package Models

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"gopkg.in/guregu/null.v3"
)

//Contract move this later
type Contract struct {
	ID                    string         `json:"id" form:"id" query:"id"`
	UsersContractsID      string         `json:"user_contracts_id" form:"user_contracts_id" query:"user_contracts_id"`
	Name                  string         `json:"name" form:"name" query:"name"`
	UpdatedAt             time.Time      `json:"updated_at" form:"updated_at" query:"updated_at"`
	EvergreenClause       null.Bool      `json:"evergreen_clause" form:"evergreen_clause" query:"evergreen_clause"`
	RenewalDate           null.Time      `json:"renewal_date" form:"renewal_date" query:"renewal_date"`
	StartDate             null.Time      `json:"start_date" form:"start_date" query:"start_date"`
	LeadSourceID          null.Int       `json:"lead_source_id" form:"lead_source_id" query:"lead_source_id"`
	LastContactID         null.Int       `json:"last_contact_id" form:"last_contact_id" query:"last_contact_id"`
	File                  null.String    `json:"file" form:"file" query:"file"`
	Facilities            types.JSONText `json:"facilities" form:"facilities" query:"facilities"`
	Contacts              types.JSONText `json:"contacts" form:"contacts" query:"contacts"`
	StatusChanges         types.JSONText `json:"status_changes" form:"status_changes" query:"status_changes"`
	CurrentContractStatus string         `json:"current_contract_status" form:"current_contract_status" query:"current_contract_status"`
}

//ContractStatus move this later
type ContractStatus struct {
	ID               string    `json:"id" form:"id" query:"id"`
	UpdatedAt        time.Time `json:"updated_at" form:"updated_at" query:"updated_at"`
	LeadNoContact    null.Time `json:"lead_no_contact" form:"lead_no_contact" query:"lead_no_contact"`
	InitialContact   null.Time `json:"initial_contact" form:"initial_contact" query:"initial_contact"`
	Interested       null.Time `json:"interested" form:"interested" query:"interested"`
	NotInterestedATM null.Time `json:"not_interested_atm" form:"not_interested_atm" query:"not_interested_atm"`
	ProposalSent     null.Time `json:"proposal_sent" form:"proposal_sent" query:"proposal_sent"`
	ContractSent     null.Time `json:"contract_sent" form:"contract_sent" query:"contract_sent"`
	ContractSigned   null.Time `json:"contract_signed" form:"contract_signed" query:"contract_signed"`
	Working          null.Time `json:"working" form:"working" query:"working"`
	Dead             null.Time `json:"dead" form:"dead" query:"dead"`
	Undesirable      null.Time `json:"undesirable" form:"undesirable" query:"undesirable"`
	ContractID       uint      `json:"contract_id" form:"contract_id" query:"contract_id"`
}
