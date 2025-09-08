package models

import "encoding/json"

type Seat struct {
	Id     int    `db:"id" json:"id"`
	Code   string `db:"codeseat" json:"seat"`
	Status bool   `db:"isstatus" json:"status"`
}
type seatAlias Seat // alias to avoid recursion

func (s Seat) MarshalJSON() ([]byte, error) {
	alias := struct {
		seatAlias
		Status string `json:"status"`
	}{
		seatAlias: seatAlias(s),
		Status:    "tidak tersedia",
	}

	if s.Status {
		alias.Status = "tersedia"
	}

	return json.Marshal(alias)
}
