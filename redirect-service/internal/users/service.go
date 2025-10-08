package users

import "errors"

type UserService struct {
	UserRepo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{UserRepo: repo}
}

func (s *UserService) GetAll() ([]string, error) {
	users := s.UserRepo.GetAll()
	if len(users) == 0 {
		return nil, errors.New("no users found")
	}
	return users, nil
}
