package repositories

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/federus1105/weekly/internals/middlewares"
	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (pr *ProfileRepository) GetProfile(rctx context.Context, UserID int) ([]models.Profile, error) {
	userIDRaw := rctx.Value(middlewares.UserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok {
		return nil, fmt.Errorf("invalid or missing user ID in context")
	}
	sql := `
	SELECT 
		u.id,
		u.email,
  COALESCE(a.image, ''),
  COALESCE(a.firstname, ''),
  COALESCE(a.lastname, ''),
  COALESCE(a.phonenumber, ''),
  a.point
	FROM users u
	JOIN account a ON a.user_id = u.id
	WHERE u.id = $1;
	`
	rows, err := pr.db.Query(rctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var profile []models.Profile
	for rows.Next() {
		var profiles models.Profile
		if err := rows.Scan(&profiles.UserID,
			&profiles.Email,
			&profiles.Image,
			&profiles.FirstName,
			&profiles.LastName,
			&profiles.Phone,
			&profiles.Point); err != nil {
			return nil, err
		}
		profile = append(profile, profiles)
	}
	return profile, nil
}

func (s *ProfileRepository) EditProfile(
	rctx context.Context,
	image *string,
	firstname *string,
	lastname *string,
	phonenumber *string,
) (models.Profile, error) {
	// Ambil user_id dari context
	userIDRaw := rctx.Value(middlewares.UserIDKey)
	userID, ok := userIDRaw.(int)
	if !ok {
		return models.Profile{}, fmt.Errorf("invalid or missing user ID in context")
	}

	// Build SQL SET clause secara dinamis
	setClauses := []string{}
	args := []any{}
	argID := 1

	if image != nil {
		setClauses = append(setClauses, fmt.Sprintf("image = $%d", argID))
		args = append(args, *image)
		argID++
	}
	if firstname != nil {
		setClauses = append(setClauses, fmt.Sprintf("firstname = $%d", argID))
		args = append(args, *firstname)
		argID++
	}
	if lastname != nil {
		setClauses = append(setClauses, fmt.Sprintf("lastname = $%d", argID))
		args = append(args, *lastname)
		argID++
	}
	if phonenumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("phonenumber = $%d", argID))
		args = append(args, *phonenumber)
		argID++
	}

	// Kalau tidak ada field yang ingin diupdate
	if len(setClauses) == 0 {
		return models.Profile{}, fmt.Errorf("no fields to update")
	}

	// Tambahkan kondisi WHERE
	query := fmt.Sprintf(`
		UPDATE account 
		SET %s 
		WHERE user_id = $%d 
		RETURNING user_id, image, firstname, lastname, phonenumber;
	`, strings.Join(setClauses, ", "), argID)

	args = append(args, userID)

	var profile models.Profile
	err := s.db.QueryRow(rctx, query, args...).Scan(
		&profile.UserID,
		&profile.Image,
		&profile.FirstName,
		&profile.LastName,
		&profile.Phone,
	)
	if err != nil {
		log.Println("Internal server error.\nCause:", err.Error())
		return models.Profile{}, err
	}

	return profile, nil
}
