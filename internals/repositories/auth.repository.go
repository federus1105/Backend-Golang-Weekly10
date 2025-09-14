package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/federus1105/weekly/internals/models"
	"github.com/federus1105/weekly/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (r *AuthRepository) Register(ctx context.Context, user models.UserRegister) (models.UserRegister, error) {
	// Start transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return models.UserRegister{}, err
	}
	defer tx.Rollback(ctx)

	hc := pkg.NewHashConfig()
	hc.UseRecommended()
	hashedPassword, err := hc.GenHash(user.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
	}

	sql := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING id, email, password, role`
	values := []any{user.Email, hashedPassword, user.Role}
	var newUser models.UserRegister
	if err := tx.QueryRow(ctx, sql, values...).Scan(&newUser.Id, &newUser.Email, &newUser.Password, &newUser.Role); err != nil {
		log.Println("Failed to insert into users: ", err.Error())
		return models.UserRegister{}, err
	}
	accountSQL := `
		INSERT INTO account (user_id)
		VALUES ($1)`

	_, err = tx.Exec(ctx, accountSQL, newUser.Id)
	if err != nil {
		log.Println("Failed to insert empty account:", err)
		return models.UserRegister{}, err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return models.UserRegister{}, err
	}

	return newUser, nil
}

func (r *AuthRepository) ResetPassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	// Step 1: Ambil hashed password dari database berdasarkan userID
	var hashedDB string
	err := r.db.QueryRow(ctx, "SELECT password FROM users WHERE id = $1", userID).Scan(&hashedDB)
	if err != nil {
		log.Println("Failed to get current password hash:", err)
		return fmt.Errorf("user tidak ditemukan")
	}

	// Step 2: Verify password lama cocok
	hc := pkg.NewHashConfig()
	hc.UseRecommended()
	ok, err := hc.CompareHashAndPassword(oldPassword, hashedDB)
	if err != nil {
		log.Println("Error comparing password hash:", err)
		return fmt.Errorf("gagal memverifikasi password lama")
	}
	if !ok {
		return fmt.Errorf("password lama tidak cocok")
	}

	// Step 3: Hash password baru
	newHashed, err := hc.GenHash(newPassword)
	if err != nil {
		log.Println("Failed to hash new password:", err)
		return fmt.Errorf("gagal hash password baru")
	}

	// Step 4: Update password di database
	_, err = r.db.Exec(ctx, "UPDATE users SET password = $1 WHERE id = $2", newHashed, userID)
	if err != nil {
		log.Println("Failed to update password:", err)
		return fmt.Errorf("gagal update password")
	}

	log.Println("Password updated for user id:", userID)
	return nil
}

// func (r *AuthRepository) Register(ctx context.Context, user models.UserRegister) (models.UserRegister, error) {
//     tx, err := r.db.Begin(ctx)
//     if err != nil {
//         log.Println("Failed to begin transaction:", err)
//         return models.UserRegister{}, err
//     }
//     // pastikan kita rollback jika gagal
//     defer tx.Rollback(ctx)

//     // 1. Insert ke tabel account dulu
//     var accountID int
//     accountSQL := `INSERT INTO account (image, firstname, lastname, phoneNumber, created_at, update_at, point)
//                    VALUES ($1, $2, $3, $4, NOW(), NOW(), $5) RETURNING id`
//     err = tx.QueryRow(ctx, accountSQL, user.Account.Image, user.Account.Firstname, user.Account.Lastname, user.Account.PhoneNumber, 0).Scan(&accountID)
//     if err != nil {
//         log.Println("Failed to insert account:", err)
//         return models.UserRegister{}, err
//     }

//     // 2. Hash password
//     hc := pkg.NewHashConfig()
//     hc.UseRecommended()
//     hashedPassword, err := hc.GenHash(user.Password)
//     if err != nil {
//         log.Println("Error hashing password:", err)
//         return models.UserRegister{}, err
//     }

//     // 3. Insert ke tabel user dengan accountID
//     userSQL := `INSERT INTO users (email, password, isactive, id_account, role)
//                 VALUES ($1, $2, $3, $4, $5)
//                 RETURNING id, email, password, isactive, id_account, role`
//     var newUser models.UserRegister
//     err = tx.QueryRow(ctx, userSQL, user.Email, hashedPassword, user.Isactive, accountID, user.Role).
//         Scan(&newUser.Id, &newUser.Email, &newUser.Password, &newUser.Isactive, &newUser.Id_account, &newUser.Role)
//     if err != nil {
//         log.Println("Failed to insert user:", err)
//         return models.UserRegister{}, err
//     }

//     // 4. Commit transaction
//     if err := tx.Commit(ctx); err != nil {
//         log.Println("Failed to commit transaction:", err)
//         return models.UserRegister{}, err
//     }

//     return newUser, nil
// }
