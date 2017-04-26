package Models

import (
	"time"

	null "gopkg.in/guregu/null.v3"
)

//User move this later
type Recruiter struct {
	ID        null.String `json:"id" form:"id" query:"id" db:"recruiter_id"`
	UpdatedAt time.Time   `json:"updated_at" form:"updated_at" query:"updated_at"`
}
