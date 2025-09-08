package repositories

import (
	"context"
	"fmt"

	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatRepository struct {
	db *pgxpool.Pool
}

func NewSeatRepository(db *pgxpool.Pool) *SeatRepository {
	return &SeatRepository{db: db}
}

// func (sr *SeatRepository) GetSeats(rctx context.Context) ([]models.Seat, error) {
// 	sql := `SELECT id, codeseat, isstatus FROM seats WHERE isstatus = TRUE`

// 	rows, err := sr.db.Query(rctx, sql)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var seats []models.Seat
// 	for rows.Next() {
// 		var seat models.Seat
// 		if err := rows.Scan(&seat.Id, &seat.Code, &seat.Status); err != nil {
// 			return nil, err
// 		}
// 		seats = append(seats, seat)
// 	}
// 	return seats, nil
// }

func (sr *SeatRepository) GetSeats(ctx context.Context, cinemaID, locationID *int) ([]models.Seat, error) {
	sql := `SELECT s.id, s.codeseat, s.isstatus 
			FROM seats s
			JOIN order_seat os ON s.id = os.id_seats
			JOIN orders o ON o.id = os.id_order
			JOIN schedule sc ON sc.id = o.id_schedule
			WHERE s.isstatus = TRUE`

	var args []interface{}
	argNum := 1

	// Filter berdasarkan id_cinema jika diberikan
	if cinemaID != nil {
		sql += fmt.Sprintf(" AND sc.id_cinema = $%d", argNum)
		args = append(args, *cinemaID)
		argNum++
	}

	// Filter berdasarkan id_location jika diberikan
	if locationID != nil {	
		sql += fmt.Sprintf(" AND sc.id_location = $%d", argNum)
		args = append(args, *locationID)
		argNum++
	}

	rows, err := sr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var seat models.Seat
		if err := rows.Scan(&seat.Id, &seat.Code, &seat.Status); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	return seats, nil
}
