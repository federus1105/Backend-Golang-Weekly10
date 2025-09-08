package repositories

import (
	"context"
	"log"

	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (or *OrderRepository) CreateOrder(rctx context.Context, body models.Order) (models.Order, error) {
	sql := `INSERT INTO orders (id_schedule, id_user, id_payment_method, total, fullname, email, phone_number, paid) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)RETURNING id, id_schedule, id_user, id_payment_method, total, fullname, email, phone_number, paid;`
	values := []any{body.Schedule, body.User, body.Payment, body.Total, body.Fullname, body.Email, body.Phone, body.Paid}
	var newOrder models.Order
	if err := or.db.QueryRow(rctx, sql, values...).Scan(&newOrder.Id, &newOrder.Schedule, &newOrder.User, &newOrder.Payment, &newOrder.Total, &newOrder.Fullname, &newOrder.Email, &newOrder.Phone, &newOrder.Paid); err != nil {
		log.Println("Internal Server Error: ", err.Error())
		return models.Order{}, err
	}
	return newOrder, nil
}
