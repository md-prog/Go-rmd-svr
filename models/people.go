package Models

import (
	"github.com/jmoiron/sqlx/types"
	"gopkg.in/guregu/null.v3"
	"time"
)

//Person move this later
type Person struct {
	ID          string      `json:"id" form:"id" query:"id"`
	UpdatedAt   time.Time   `json:"updated_at" form:"updated_at" query:"updated_at"`
	FirstName   string      `json:"first_name" form:"first_name" query:"first_name"`
	LastName    string      `json:"last_name" form:"last_name" query:"last_name"`
	MiddleName  null.String `json:"middle_name" form:"middle_name" query:"middle_name"`
	Title       null.String `json:"title" form:"title" query:"title"`
	Email       string      `json:"email" form:"email" query:"email"`
	CellPhone   null.String `json:"cell_phone" form:"cell_phone" query:"cell_phone"`
	HomePhone   null.String `json:"home_phone" form:"home_phone" query:"home_phone"`
	OfficePhone null.String `json:"office_phone" form:"office_phone" query:"office_phone"`
	Fax         null.String `json:"fax" form:"fax" query:"fax"`
	Address     null.String `json:"address" form:"address" query:"address"`
	City        null.String `json:"city" form:"city" query:"city"`
	Zip         null.String `json:"zip" form:"zip" query:"zip"`
	Lat         null.Float  `json:"lat" form:"lat" query:"lat"`
	Lon         null.Float  `json:"lon" form:"lon" query:"lon"`
	StateID     null.String `json:"state_id" form:"state_id" query:"state_id"`
	State       `json:"state" form:"state" query:"state" db:"state"`
	RecruiterID null.String `json:"recruiter_id" form:"recruiter_id" query:"recruiter_id"`
	Recruiter   `json:"recruiter" form:"recruiter" query:"recruiter" db:"recruiter"`
	ProviderID  null.String `json:"provider_id" form:"provider_id" query:"provider_id"`
	Provider    `json:"provider" form:"provider" query:"provider" db:"provider"`
	Notes       null.String    `json:"notes" form:"notes" query:"notes"`
	Contacts    types.JSONText `json:"contacts" form:"contacts" query:"contacts"`
}
