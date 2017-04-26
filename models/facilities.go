package Models

import (
	"github.com/jmoiron/sqlx/types"
	null "gopkg.in/guregu/null.v3"
	"time"
)

//Facility for reading in user data from api calls
type Facility struct {
	ID                 string         `json:"id" form:"id" query:"id"`
	UpdatedAt          time.Time      `json:"updated_at" form:"updated_at" query:"updated_at"`
	Name               string         `json:"name" form:"name" query:"name"`
	Hospitalist        Hospitalist    `json:"hospitalist" form:"hospitalist" query:"hospitalist"`
	UserFacilitiesID   null.String    `json:"user_facilities_id" form:"user_facilities_id" query:"user_facilities_id"`
	FacilityTypeID     null.String    `json:"facility_type_id" form:"facility_type_id" query:"facility_type_id"`
	EMRID              null.String    `json:"emr_id" form:"emr_id" query:"emr_id"`
	Address            null.String    `json:"address" form:"address" query:"address"`
	City               null.String    `json:"city" form:"city" query:"city"`
	Zip                null.String    `json:"zip" form:"zip" query:"zip"`
	Lat                null.Float     `json:"lat" form:"lat" query:"lat"`
	Lon                null.Float     `json:"lon" form:"lon" query:"lon"`
	TypeID             null.String    `json:"type_id" form:"type_id" query:"type_id"`
	ContractID         null.String    `json:"contract_id" form:"contract_id" query:"contract_id"`
	FacilityStatsID    null.String    `json:"facility_stats_id" form:"facility_stats_id" query:"facility_stats_id"`
	PrimaryContactID   null.String    `json:"primary_contact_id" form:"primary_contact_id" query:"primary_contact_id"`
	SecondaryContactID null.String    `json:"secondary_contact_id" form:"secondary_contact_id" query:"secondary_contact_id"`
	AdmissionOrders    null.Bool      `json:"admission_orders" form:"admission_orders" query:"admission_orders"`
	HospitalistID      null.String    `json:"hospitalist_id" form:"hospitalist_id" query:"hospitalist_id"`
	StateID            null.String    `json:"state_id" form:"state_id" query:"state_id"`
	LastContactID      null.String    `json:"last_contact_id" form:"last_contact_id" query:"last_contact_id"`
	Scribes            null.Bool      `json:"scribes" form:"scribes" query:"scribes"`
	Contacts           types.JSONText `json:"contacts" form:"contacts" query:"contacts"`
}

// last_contact_id	int	null

//FacilityType move this later
type FacilityType struct {
	ID        string    `json:"id" form:"id" query:"id"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" query:"updated_at"`
	Name      string    `json:"name" form:"name" query:"name"`
}

//Hospitalist move this later
type Hospitalist struct {
	ID               string    `json:"id" form:"id" query:"id"`
	UpdatedAt        time.Time `json:"updated_at" form:"updated_at" query:"updated_at"`
	Name             string    `json:"name" form:"name" query:"name"`
	PrimaryContactID uint      `json:"primary_contact_id" form:"primary_contact_id" query:"primary_contact_id"`
}

//EMR move this later
type EMR struct {
	ID        string    `json:"id" form:"id" query:"id"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" query:"updated_at"`
	Name      string    `json:"name" form:"name" query:"name"`
}

type FacilityStats struct {
	ID                     string    `json:"id" form:"id" query:"id"`
	UpdatedAt              time.Time `json:"updated_at" form:"updated_at" query:"updated_at"`
	EntryDate              time.Time `json:"entry_date" form:"entry_date" query:"entry_date"`
	BedCountFT             null.Int  `json:"bed_count_ft" form:"bed_count_ft" query:"bed_count_ft"`
	BedCountMain           null.Int  `json:"bed_count_main" form:"bed_count_main" query:"bed_count_main"`
	BedCountOverflow       null.Int  `json:"bed_count_overflow" form:"bed_count_overflow" query:"bed_count_overflow"`
	AnnualVolume           null.Int  `json:"annual_volume" form:"annual_volume" query:"annual_volume"`
	ProviderInTriage       null.Bool `json:"provider_in_triage" form:"provider_in_triage" query:"provider_in_triage"`
	HoursPhysician         null.Int  `json:"hours_physician" form:"hours_physician" query:"hours_physician"`
	HoursAPC               null.Int  `json:"hours_apc" form:"hours_apc" query:"hours_apc"`
	ShiftsPhysician        null.Int  `json:"shifts_physician" form:"shifts_physician" query:"shifts_physician"`
	ShiftsAPC              null.Int  `json:"shifts_apc" form:"shifts_apc" query:"shifts_apc"`
	AnnualHospitalizations null.Int  `json:"annual_hospitalizations" form:"annual_hospitalizations" query:"annual_hospitalizations"`
	AnnualObs              null.Int  `json:"annual_obs" form:"annual_obs" query:"annual_obs"`
	AnnualAdmit            null.Int  `json:"annual_admit" form:"annual_admit" query:"annual_admit"`
	LWBS                   null.Int  `json:"lwbs" form:"lwbs" query:"lwbs"`
	LWOT                   null.Int  `json:"lwot" form:"lwot" query:"lwot"`
	FastTrackHoursStart    null.Time `json:"fast_track_hours_start" form:"fast_track_hours_start" query:"fast_track_hours_start"`
	FastTrackHoursEnd      null.Time `json:"fast_track_hours_end" form:"fast_track_hours_end" query:"fast_track_hours_end"`
}
