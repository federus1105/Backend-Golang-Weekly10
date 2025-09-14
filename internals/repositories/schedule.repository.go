package repositories

import (
	"context"

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
	WHERE s.id = $1
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

// func (sr *ScheduleRepository) GetSchedule(ctx context.Context, id int) ([]models.Schedule, error) {
// 	start := time.Now()
// 	redisKey := "firdaus:schedule-all"
// 	cmd := sr.rdb.Get(ctx, redisKey)
// 	if cmd.Err() != nil {
// 		if cmd.Err() == redis.Nil {
// 			log.Printf("Key %s does not exist\n", redisKey)
// 		} else {
// 			log.Println("Redis Error. \nCause: ", cmd.Err().Error())
// 		}
// 	} else {
// 		// cache hit
// 		var cachedSchedules []models.Schedule
// 		cmdByte, err := cmd.Bytes()
// 		if err != nil {
// 			log.Println("Internal server error.\nCause: ", err.Error())
// 		} else {
// 			if err := json.Unmarshal(cmdByte, &cachedSchedules); err != nil {
// 				log.Println("Internal Server Error. \nCause: ", err.Error())
// 			}
// 		}
// 		if len(cachedSchedules) > 0 {
// 			log.Printf("Key %s found in cache ✅", redisKey)
// 			log.Printf("Served in %s using Redis", time.Since(start))
// 			return cachedSchedules, nil
// 		}
// 	}
// 	sql := `SELECT
// s.id,
// s.date,
// 		m.title AS title,
// 		m.image AS image,
// c.name AS cinema,
// t.name AS time,
// l.name AS location
// 	FROM schedule s
// 	JOIN movies m ON s.id_movie = m.id
// 	LEFT JOIN cinema c ON s.id_cinema = c.id
// 	LEFT JOIN time t ON s.id_time = t.id
// 	LEFT JOIN location l ON s.id_location = l.id
// 	WHERE s.id = $1
// ORDER BY s.date ASC`
// 	rows, err := sr.db.Query(ctx, sql, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var schedules []models.Schedule
// 	for rows.Next() {
// 		var schedule models.Schedule
// 		if err := rows.Scan(&schedule.Id, &schedule.Date, &schedule.Title, &schedule.Image, &schedule.Cinema, &schedule.Time, &schedule.Location); err != nil {
// 			return nil, err
// 		}
// 		schedules = append(schedules, schedule)
// 	}
// 	// renew cache
// 	bt, err := json.Marshal(schedules)
// 	if err != nil {
// 		log.Println("Internal Server Error.\n Cause: ", err.Error())
// 	}
// 	if err := sr.rdb.Set(ctx, redisKey, string(bt), 5*time.Minute).Err(); err != nil {
// 		log.Println("Redis Error. \nCause: ", err.Error())
// 	}
// 	log.Printf("[REDIS TIMING] Served in %s using DB (cache miss)", time.Since(start))
// 	return schedules, nil
// }
