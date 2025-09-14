package repositories

import (
	"context"
	"log"

	"github.com/federus1105/weekly/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (pr *ProfileRepository) GetProfile(rctx context.Context, id int) (models.Profile, error) {
	sql := `
	SELECT 
		u.id,
		u.email,
  COALESCE(a.image, 'belum ada data'),
  COALESCE(a.firstname, 'belum ada data'),
  COALESCE(a.lastname, 'belum ada data'),
  COALESCE(a.phonenumber, 'belum ada data')
	FROM users u
	JOIN account a ON a.user_id = u.id
	WHERE u.id = $1;
	`
	var profile models.Profile
	err := pr.db.QueryRow(rctx, sql, id).Scan(
		&profile.UserID,
		&profile.Email,
		&profile.Image,
		&profile.FirstName,
		&profile.LastName,
		&profile.Phone,
	)
	if err != nil {
		log.Println("Error GetProfile:", err)
		return models.Profile{}, err
	}
	return profile, nil
}


func (s *ProfileRepository) EditProfile(
	rctx context.Context,
	Image string,
	firstname string,
	lastname string,
	phonenumber string,
	id int,
) (models.Profile, error) {
	sql := `
		UPDATE account 
		SET image = $1, firstname = $2, lastname = $3, phonenumber = $4 
		WHERE user_id = $5 
		RETURNING user_id, image, firstname, lastname, phonenumber;
	`

	values := []any{Image, firstname, lastname, phonenumber, id}

	var profile models.Profile
	err := s.db.QueryRow(rctx, sql, values...).Scan(
		&profile.UserID,
		&profile.Image,
		&profile.FirstName,
		&profile.LastName,
		&profile.Phone,
	)
	if err != nil {
		log.Println("Internal server error.\nCause: ", err.Error())
		return models.Profile{}, err
	}
	return profile, nil
}
