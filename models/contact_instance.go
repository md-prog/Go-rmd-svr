package Models

import "time"

//Contact move this later
type ContactInstance struct {
	ID            string    `json:"id" form:"id" query:"id"`
	UpdatedAt     time.Time `json:"updated_at" form:"updated_at" query:"updated_at"`
	ContactType   string    `json:"contact_type" form:"contact_type" query:"contact_type"`
	UserContactID uint      `json:"user_contact_id" form:"user_contact_id" query:"user_contact_id"`
}
