package repository

import (
	"database/sql"
	"errors"
	"go-api/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var u models.User
	err := r.DB.QueryRow("SELECT id, name, email FROM users WHERE id=$1", id).Scan(&u.ID, &u.Name, &u.Email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, string, error) {
	// return user WITHOUT password and passwordHash separately
	var u models.User
	var pw string
	err := r.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email=$1", email).Scan(&u.ID, &u.Name, &u.Email, &pw)
	if err == sql.ErrNoRows {
		return nil, "", nil
	}
	if err != nil {
		return nil, "", err
	}
	return &u, pw, nil
}

func (r *UserRepository) Create(u *models.User, passwordHash string) error {
	err := r.DB.QueryRow(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id",
		u.Name, u.Email, passwordHash,
	).Scan(&u.ID)
	return err
}

func (r *UserRepository) Update(id int, u models.User) error {
	res, err := r.DB.Exec("UPDATE users SET name=$1, email=$2 WHERE id=$3", u.Name, u.Email, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *UserRepository) DeleteById(id int) error {
	res, err := r.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *UserRepository) DeleteByName(name string) error {
	res, err := r.DB.Exec("DELETE FROM users WHERE name=$1", name)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return errors.New("not found")
	}
	return nil
}
