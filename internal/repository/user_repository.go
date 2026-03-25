package repository

import "url_shortener/internal/models"

type userRepository struct {
	users []models.User
}

func NewUserRepository() UserRepository {
	return &userRepository{
		users: make([]models.User, 0),
	}
}
func (ur userRepository) FindAll() {

}
func (ur userRepository) Create() {

}
func (ur userRepository) Delete() {

}
func (ur userRepository) Update() {

}
