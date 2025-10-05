package repositories

import (
	"context"
	"log"

	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ScheduleRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewScheduleRepository(db *pgxpool.Pool, rdb *redis.Client) *ScheduleRepository {
	return &ScheduleRepository{db: db, rdb: rdb}
}

func (sr *ScheduleRepository) GetSchedule(rctx context.Context, id_movie int) ([]models.Schedule, error) {
	sql := `SELECT
s.id,
s.id_movie,
s.date,
		m.title AS title,
c.id AS idcinema,
c.name AS cinema,
t.name AS time,
l.name AS location,
c.image as icon
	FROM schedule s
	JOIN movies m ON s.id_movie = m.id
	LEFT JOIN cinema c ON s.id_cinema = c.id
	LEFT JOIN time t ON s.id_time = t.id
	LEFT JOIN location l ON s.id_location = l.id
	WHERE m.id = $1
ORDER BY s.date ASC`

	rows, err := sr.db.Query(rctx, sql, id_movie)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		if err := rows.Scan(&schedule.Id, &schedule.Idmovie, &schedule.Date, &schedule.Title, &schedule.Id_Cinema, &schedule.Cinema, &schedule.Time, &schedule.Location, &schedule.Image); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

func (sr *ScheduleRepository) CreateSchedule(
	rctx context.Context,
	input models.BodyScheduleInput,
) ([]models.BodySchedule, error) {
	sql := `INSERT INTO schedule (id_movie, date, id_cinema, id_time, id_location)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id, id_movie, date, id_cinema, id_time, id_location`

	var createdSchedules []models.BodySchedule

	for _, cinemaID := range input.Id_Cinema {
		for _, timeID := range input.Time {
			for _, locationID := range input.Location {
				values := []any{input.Id_movie, input.Date, cinemaID, timeID, locationID}
				var newSchedule models.BodySchedule

				err := sr.db.QueryRow(rctx, sql, values...).Scan(
					&newSchedule.Id,
					&newSchedule.Id_movie,
					&newSchedule.Date,
					&newSchedule.Id_Cinema,
					&newSchedule.Id_Time,
					&newSchedule.Id_Location,
				)
				if err != nil {
					log.Println("Failed to insert schedule:", err)
					return nil, err
				}

				createdSchedules = append(createdSchedules, newSchedule)
			}
		}
	}

	return createdSchedules, nil
}
