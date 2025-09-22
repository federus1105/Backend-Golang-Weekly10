package repositories

import (
	"context"
	"fmt"

	"github.com/federus1105/weekly/internals/middlewares"
	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HistoryRepository struct {
	db *pgxpool.Pool
}

func NewHistoryRepository(db *pgxpool.Pool) *HistoryRepository {
	return &HistoryRepository{db: db}
}

func (hr *HistoryRepository) GetHistory(rctx context.Context, UserID int) ([]models.History, error) {
	userIDRaw := rctx.Value(middlewares.UserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok {
		return nil, fmt.Errorf("invalid or missing user ID in context")
	}

	sql := `SELECT
      o.id AS id_order,
      m.title AS movie_title,
      STRING_AGG(s2.codeseat, ', ') AS seat_codes,
      COUNT(os.id_seats) AS total_seats,
      t.name AS time_name,
      o.total,
      c.name AS cinema_name,
      o.paid
    FROM orders o
    JOIN schedule s ON o.id_schedule = s.id
    JOIN movies m ON s.id_movie = m.id
    JOIN cinema c ON s.id_cinema = c.id
    JOIN time t ON s.id_time = t.id
    LEFT JOIN order_seat os ON o.id = os.id_order
    LEFT JOIN seats s2 ON os.id_seats = s2.id 
    WHERE o.id_user = $1
    GROUP BY o.id, m.title, t.name, o.total, c.name, o.paid
    ORDER BY o.created_at ASC;`

	rows, err := hr.db.Query(rctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []models.History
	for rows.Next() {
		var history models.History
		if err := rows.Scan(&history.IDOrder, &history.Movie, &history.Seat, &history.TotalSeat, &history.Time, &history.Total, &history.Cinema, &history.Paid); err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}
	return histories, nil
}
