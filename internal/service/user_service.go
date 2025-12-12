package service

import (
	"errors"
	"go-api/internal/models"
	"go-api/internal/repository"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{Repo: r}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.Repo.GetByID(id)
}

func (s *UserService) CreateUser(u *models.User, passwordHash string) error {
	if err := u.ValidateForCreate(); err != nil {
		return err
	}
	// чек на уникальность email
	existing, _, err := s.Repo.GetByEmail(u.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already exists")
	}
	return s.Repo.Create(u, passwordHash)
}

func (s *UserService) UpdateUser(id int, u models.User) error {
	if err := u.ValidateForUpdate(); err != nil {
		return err
	}
	return s.Repo.Update(id, u)
}

func (s *UserService) DeleteById(id int) error {
	return s.Repo.DeleteById(id)
}

func (s *UserService) DeleteByName(name string) error {
	return s.Repo.DeleteByName(name)
}
