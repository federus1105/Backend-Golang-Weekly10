package models

import (
	"encoding/json"
)

type Seat struct {
	Id     int     `db:"id" json:"id"`
	Code   string  `db:"codeseat" json:"seat"`
	Status bool    `db:"isstatus" json:"status"`
}
func (s Seat) MarshalJSON() ([]byte, error) {
	type SeatAlias Seat

	var statusText string
	if s.Status {
		statusText = "tersedia"
	} else {
		statusText = "terjual"
	}

	return json.Marshal(struct {
		SeatAlias
		Status string `json:"status"`
	}{
		SeatAlias: SeatAlias(s),
		Status:    statusText,
	})
}
