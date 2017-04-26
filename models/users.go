package Models

import (
	"time"

	null "gopkg.in/guregu/null.v3"
)

//User move this later
type User struct {
	ID        string      `json:"id" form:"id" query:"id"`
	UpdatedAt time.Time   `json:"updated_at" form:"updated_at" query:"updated_at"`
	Email     string      `json:"email" form:"email" query:"email"`
	Password  string      `json:"password" form:"password" query:"password"`
	PeopleID  null.String `json:"people_id" form:"people_id" query:"people_id"`
}
