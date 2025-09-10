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

func (pr *ProfileRepository) GetProfile(rctx context.Context, Id int) ([]models.Profile, error) {
	sql := `SELECT 
  u.id,
  u.email,
  a.image, 
  a.firstname, 
  a.lastname, 
  a.phoneNumber
FROM users u
JOIN account a ON u.id_account = a.id
WHERE u.id = $1;`

	rows, err := pr.db.Query(rctx, sql, Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var profiles []models.Profile
	for rows.Next() {
		var Profile models.Profile
		if err := rows.Scan(&Profile.ID, &Profile.Email, &Profile.Image, &Profile.FirstName, &Profile.LastName, &Profile.Phone); err != nil {
			return nil, err
		}
		profiles = append(profiles, Profile)
	}
	return profiles, nil
}

func (s *ProfileRepository) EditProfile(rctx context.Context, firstname string, lastname string, phonenumber string, id int) (models.Profile, error) {
	sql := "UPDATE account SET firstname=$1, lastname=$2, phonenumber=$3 WHERE id =$4 RETURNING id, firstname, lastname, phonenumber"
	log.Println("Updating profile: ", firstname, lastname, phonenumber, id)
	values := []any{ firstname, lastname, phonenumber, id}
	var profile models.Profile
	err := s.db.QueryRow(rctx, sql, values...).Scan(&profile.ID, &profile.FirstName, &profile.LastName, &profile.Phone)
	if err != nil {
		log.Println("Internal server error.\nCause: ", err.Error())
		return models.Profile{}, err
	}
	return profile, nil
}
