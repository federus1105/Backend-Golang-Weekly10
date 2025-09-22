package repositories

import (
	"context"
	"fmt"
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

//	func (or *OrderRepository) CreateOrder(rctx context.Context, body models.Order) (models.Order, error) {
//		sql := `INSERT INTO orders (id_schedule, id_user, id_payment_method, total, fullname, email, phone_number, paid)
//		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)RETURNING id, id_schedule, id_user, id_payment_method, total, fullname, email, phone_number, paid;`
//		values := []any{body.Schedule, body.User, body.Payment, body.Total, body.Fullname, body.Email, body.Phone, body.Paid}
//		var newOrder models.Order
//		if err := or.db.QueryRow(rctx, sql, values...).Scan(&newOrder.Id, &newOrder.Schedule, &newOrder.User, &newOrder.Payment, &newOrder.Total, &newOrder.Fullname, &newOrder.Email, &newOrder.Phone, &newOrder.Paid); err != nil {
//			log.Println("Internal Server Error: ", err.Error())
//			return models.Order{}, err
//		}
//		return newOrder, nil
//	}
func (or *OrderRepository) CreateOrder(
	rctx context.Context,
	body models.Order,
	seatIDs []int,
) (newOrder models.Order, err error) {
	// Begin transaction
	tx, err := or.db.Begin(rctx)
	if err != nil {
		log.Println("Failed to start transaction:", err)
		return
	}
	defer tx.Rollback(rctx)

	// defer func() {
	// 	if err != nil {
	// 		_ = tx.Rollback(rctx)
	// 	}
	// }()

	// ✅ Ambil harga cinema berdasarkan schedule
	var price int
	sqlPrice := `
		SELECT c.price
		FROM schedule s
		JOIN cinema c ON s.id_cinema = c.id
		WHERE s.id = $1;
	`
	err = tx.QueryRow(rctx, sqlPrice, body.Schedule).Scan(&price)
	if err != nil {
		log.Println("Failed to get cinema price:", err)
		return
	}
	

	// ✅ Hitung total otomatis
	body.Total = float32(price * len(seatIDs))

	// Step 1: Insert ke orders
	sqlOrder := `INSERT INTO orders 
	(id_schedule, id_user, id_payment_method, total, fullname, email, phone_number, paid) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, id_schedule, id_user, id_payment_method, total, fullname, email, phone_number, paid;`

	err = tx.QueryRow(rctx, sqlOrder,
		body.Schedule, body.User, body.Payment, body.Total,
		body.Fullname, body.Email, body.Phone, body.Paid,
	).Scan(
		&newOrder.Id, &newOrder.Schedule, &newOrder.User,
		&newOrder.Payment, &newOrder.Total, &newOrder.Fullname,
		&newOrder.Email, &newOrder.Phone, &newOrder.Paid,
	)

	if err != nil {
		log.Println("Failed to insert order:", err)
		return
	}

	for _, seatID := range seatIDs {
		// Validasi kursi belum diambil
		var isAvailable bool
		sqlCheck := `SELECT isstatus FROM seats WHERE id = $1`
		err = tx.QueryRow(rctx, sqlCheck, seatID).Scan(&isAvailable)
		if err != nil {
			log.Println("Failed to check seat status:", err)
			return
		}

		if !isAvailable {
			err = fmt.Errorf("seat with ID %d is already booked", seatID)
			log.Println(err)
			return
		}

		// Insert ke order_seat
		sqlSeat := `INSERT INTO order_seat (id_order, id_seats) VALUES ($1, $2);`
		if _, err = tx.Exec(rctx, sqlSeat, newOrder.Id, seatID); err != nil {
			log.Println("Failed to insert order_seat:", err)
			return
		}

		// Update status kursi
		sqlUpdate := `UPDATE seats SET isstatus = false WHERE id = $1;`
		if _, err = tx.Exec(rctx, sqlUpdate, seatID); err != nil {
			log.Println("Failed to update seat status:", err)
			return
		}
	}

	// Commit
	if err = tx.Commit(rctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return
	}

	return newOrder, nil
}
