package repositories

import (
	"context"

	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentRepository struct {
	db *pgxpool.Pool
}

func NewPaymentRepository(db *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (pr *PaymentRepository) GetPayment(rctx context.Context) ([]models.Payment, error) {
	rows, err := pr.db.Query(rctx, "SELECT id, name, image FROM payment_method")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payment []models.Payment
	for rows.Next() {
		var Payment models.Payment
		if err := rows.Scan(&Payment.Id, &Payment.Name, &Payment.Image); err != nil {
			return nil, err
		}
		payment = append(payment, Payment)
	}

	return payment, nil
}
