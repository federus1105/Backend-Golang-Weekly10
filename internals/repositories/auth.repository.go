package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) GetUserWithPasswordAndRole(rctx context.Context, email string) (models.User, error) {
	// validasi user
	// ambil data user berdasarkan input user
	sql := `SELECT id, email, password, role FROM users WHERE email = $1`

	var user models.User
	if err := a.db.QueryRow(rctx, sql, email).Scan(&user.Id, &user.Email, &user.Password, &user.Role); err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		return models.User{}, err
	}
	return user, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (r *AuthRepository) CreateUser(ctx context.Context, user models.User) error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	sql := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3)`
	_, err = r.db.Exec(ctx, sql, user.Email, hashedPassword, user.Role)
	return err
}

func (r *AuthRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	sql := `SELECT id, email, password, role FROM users`
	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return users, nil
}

func (r *AuthRepository) UpdateUserPassword(ctx context.Context, userID uint, hashedPassword string) error {
	sql := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, sql, hashedPassword, userID)
	return err
}