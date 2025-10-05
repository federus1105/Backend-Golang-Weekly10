package configs

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	// "github.com/jackc/pgx/v5/pgxpool"
)

func InitRDB() (*redis.Client, string, error) {
	rdbUser := os.Getenv("REDISUSER")
	rdbPass := os.Getenv("REDISPASS")
	rdbHost := os.Getenv("REDISHOST")
	// rdbHost := os.Getenv("DBSHOST")
	rdbPort := os.Getenv("REDISPORT")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rdbHost, rdbPort),
		Username: rdbUser,
		Password: rdbPass,
		DB:       0,
	})
	return rdb, rdbUser, nil
}
