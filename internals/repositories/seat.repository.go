package repositories

import (
	"context"

	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatRepository struct {
	db *pgxpool.Pool
}

func NewSeatRepository(db *pgxpool.Pool) *SeatRepository {
	return &SeatRepository{db: db}
}

func (sr *SeatRepository) GetSeats(rctx context.Context, Id int) ([]models.Seat, error) {
	sql := `WITH target_schedule AS (
	SELECT sc.id AS schedule_id, sc.id_cinema, c.price
	FROM schedule sc
	JOIN cinema c ON c.id = sc.id_cinema
	WHERE sc.id = $1
)
SELECT
	s.id,
	s.codeseat,
CASE
	WHEN o.id_schedule = t.schedule_id AND sc.id_cinema = t.id_cinema THEN false
	ELSE true
END AS isstatus,
	t.price AS seat_price
FROM seats s
LEFT JOIN order_seat os ON s.id = os.id_seats
LEFT JOIN orders o ON o.id = os.id_order
LEFT JOIN schedule sc ON sc.id = o.id_schedule
LEFT JOIN target_schedule t ON true
ORDER BY s.codeseat ASC;`

	rows, err := sr.db.Query(rctx, sql, Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var seats []models.Seat
	for rows.Next() {
		var Seat models.Seat
		if err := rows.Scan(&Seat.Id, &Seat.Code, &Seat.Status, &Seat.Price); err != nil {
			return nil, err
		}
		seats = append(seats, Seat)
	}
	return seats, err
}

// SELECT
// 	s.id,
// 	s.codeseat,
// 	CASE
// 		WHEN o.id_schedule = sc_target.id AND sc.id_cinema = sc_target.id_cinema THEN false
// 		ELSE true
// 	END AS isAvailable,
// 	c_target.price AS seat_price
// FROM seats s
// LEFT JOIN order_seat os ON s.id = os.id_seats
// LEFT JOIN orders o ON o.id = os.id_order
// LEFT JOIN schedule sc ON sc.id = o.id_schedule
// LEFT JOIN cinema c ON c.id = sc.id_cinema
// -- Ganti target schedule & cinema di sini
// LEFT JOIN schedule sc_target ON sc_target.id
// LEFT JOIN cinema c_target ON c_target.id = sc_target.id_cinema
// WHERE sc_target.id = 9 IS NOT NULL
// ORDER BY s.codeseat ASC;