package repository

type UserRepository interface {
	FindAll()
	Create()
	Update()
	Delete()
}
