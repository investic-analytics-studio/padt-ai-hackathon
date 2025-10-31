package repo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/quantsmithapp/datastation-backend/internal/core/port"
	"github.com/quantsmithapp/datastation-backend/internal/model"
)

type AuthRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) port.AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) Create(body *model.SignUpBody, methodSignUp string) error {
	fmt.Println("methodSignUp", methodSignUp)

	query := `
        INSERT INTO crypto_user (uuid, email,last_update,twitter_uid,twitter_name,method_login)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	if methodSignUp == model.SignUpMethodEmail || methodSignUp == model.SignUpMethodGoogle {
		_, err := r.db.Exec(query, body.UID, body.Email, body.Email, nil, nil, "email")
		return err
	} else if methodSignUp == model.SignUpXMethod {
		_, err := r.db.Exec(query, body.UID, body.Email, body.TwitterUID, body.TwitterUID, body.TwitterName, "twitter")
		return err
	}
	return nil
}

func (r *AuthRepo) GetUserByEmail(email string) (*model.User, error) {
	query := `
        SELECT email 
		FROM crypto_user
		WHERE email = $1
    `

	var emailResult []string
	err := r.db.Select(&emailResult, query, email)
	if err != nil {
		return nil, err
	}

	return &model.User{Email: emailResult[0]}, nil
}

func (r *AuthRepo) ExistEmail(email string) (bool, error) {
	query := `
        SELECT COUNT(email)
        FROM crypto_user
        WHERE email = $1
    `
	var count []int
	err := r.db.Select(&count, query, email)
	return count[0] > 0, err
}

func (r *AuthRepo) CheckUserByUid(uid string) (bool, error) {
	query := `
    SELECT COUNT(uuid)
    FROM crypto_user
    WHERE uuid = $1
`
	var count []int
	err := r.db.Select(&count, query, uid)
	return count[0] > 0, err
}
func (r *AuthRepo) CheckCRMUserByUid(uid string) (bool, error) {
	query := `
    SELECT COUNT(id)
    FROM crypto_crm_user
    WHERE id = $1
`
	var count []int
	err := r.db.Select(&count, query, uid)
	return count[0] > 0, err
}
func (r *AuthRepo) GetCRMUserInfo(uid string) (model.CRMUser, error) {
	query := `
	SELECT id, username
	FROM crypto_crm_user
	WHERE id = $1
	`

	var result struct {
		ID       string `db:"id"`
		Username string `db:"username"`
	}
	// Execute the query
	err := r.db.Get(&result, query, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found for UID: %s\n", uid)
			return model.CRMUser{}, nil // Returning empty user info without error
		}
		log.Printf("Error fetching user info for UID: %s - %v\n", uid, err)
		return model.CRMUser{}, err
	}
	// Return the result

	// Return the result with null values converted to empty strings
	return model.CRMUser{
		ID:       result.ID,
		Username: result.Username,
	}, nil
}
func (r *AuthRepo) GetUserInfo(uid string) (port.UserInfo, error) {
	query := `
	SELECT email, twitter_uid, twitter_name
	FROM crypto_user
	WHERE uuid = $1
	`
	var result struct {
		Email       sql.NullString `db:"email"`
		TwitterUID  sql.NullString `db:"twitter_uid"`
		TwitterName sql.NullString `db:"twitter_name"`
	}
	// Execute the query
	err := r.db.Get(&result, query, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No user found for UID: %s\n", uid)
			return port.UserInfo{}, nil // Returning empty user info without error
		}
		log.Printf("Error fetching user info for UID: %s - %v\n", uid, err)
		return port.UserInfo{}, err
	}
	// Return the result

	// Return the result with null values converted to empty strings
	return port.UserInfo{
		Email:       result.Email.String,
		TwitterUID:  result.TwitterUID.String,
		TwitterName: result.TwitterName.String,
	}, nil
}

func (r *AuthRepo) ExistTwitterUID(uid string) (bool, error) {
	query := `
	SELECT COUNT(twitter_uid)
	FROM crypto_user
	WHERE twitter_uid = $1
	`
	var count []int
	err := r.db.Select(&count, query, uid)
	return count[0] > 0, err
}
