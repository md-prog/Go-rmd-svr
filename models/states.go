package Models

import (
	"time"

	null "gopkg.in/guregu/null.v3"
)

type State struct {
	ID           null.String `json:"id" form:"id" query:"id"`
	UpdatedAt    time.Time   `json:"updated_at" form:"updated_at" query:"updated_at"`
	Name         null.String `json:"name" form:"name" query:"name" db:"state_name"`
	Abbreviation null.String `json:"abbr" form:"abbr" query:"abbr" db:"state_abbr"`
}
