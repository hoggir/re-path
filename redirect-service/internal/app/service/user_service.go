package service

import (
	"errors"

	"github.com/hoggir/re-path/redirect-service/internal/app/repository"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}

func (s *UserService) GetAll() ([]string, error) {
	users := s.UserRepo.GetAll()
	if len(users) == 0 {
		return nil, errors.New("no users found")
	}
	return users, nil
}
