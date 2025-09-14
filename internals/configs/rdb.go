package configs

import (
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	// "github.com/jackc/pgx/v5/pgxpool"
)

func InitRDB() (*redis.Client, error) {
	rdbUser := os.Getenv("REDISUSER")
	rdbPass := os.Getenv("REDISPASS")
	rdbHost := os.Getenv("DBHOST")
	rdbPort := os.Getenv("REDISPORT")

	// connstring := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	// return pgxpool.New(context.Background(), connstring)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rdbHost, rdbPort),
		Username: rdbUser,
		Password: rdbPass,
		DB:       0,
	})
	return rdb, nil
}

// func TestRDB(rdb *redis.Client) error {
// 	return rdb.Ping(ctx context.Context.Background())
// }
