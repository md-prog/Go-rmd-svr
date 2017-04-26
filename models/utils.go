package Models

import (
	null "gopkg.in/guregu/null.v3"
)

func HandleNullString(s null.String) null.String {
	x, _ := s.Value()
	if x == "" {
		return null.NewString("", false)
	} else {
		return s
	}
}
