package repositories

import (
	"context"

	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepository struct {
	db *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (sr *ScheduleRepository) GetSchedule(rctx context.Context, id int) ([]models.Schedule, error) {
	sql := `SELECT
s.id,
s.date,
		m.title AS title,
		m.image AS image,
c.name AS cinema,
t.name AS time,
l.name AS location
	FROM schedule s
	JOIN movies m ON s.id_movie = m.id
	LEFT JOIN cinema c ON s.id_cinema = c.id
	LEFT JOIN time t ON s.id_time = t.id
	LEFT JOIN location l ON s.id_location = l.id
	WHERE s.id_movie = $1
ORDER BY s.date ASC`

	rows, err := sr.db.Query(rctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		if err := rows.Scan(&schedule.Id, &schedule.Date, &schedule.Title, &schedule.Image, &schedule.Cinema, &schedule.Time, &schedule.Location); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}
