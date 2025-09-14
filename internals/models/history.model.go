package models

type History struct {
	IDOrder   int     `db:"id" json:"id_order"`
	Movie     string  `db:"movie" json:"movie_title"`
	Seat      string  `db:"seat_codes" json:"seat"`
	TotalSeat int     `db:"total_seats" json:"total_seats"`
	Time      string  `db:"time_name" json:"time"`
	Total     float32 `db:"total" json:"total"`
	Cinema    string  `db:"name" json:"cinema"`
	Paid      bool    `db:"paid" json:"paid"`
}
