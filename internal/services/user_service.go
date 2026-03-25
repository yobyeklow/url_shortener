package services

import "url_shortener/internal/repository"

type userService struct {
	user_repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserServices {
	return &userService{
		user_repo: repo,
	}
}

func (us *userService) GetAllUser() {
	us.user_repo.FindAll()
}
func (us *userService) CreateUser() {
	us.user_repo.FindAll()
}
func (us *userService) UpdateUser() {
	us.user_repo.FindAll()
}
func (us *userService) DeleteUser() {
	us.user_repo.FindAll()
}
