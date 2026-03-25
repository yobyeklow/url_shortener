package handler

import (
	"url_shortener/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service services.UserServices
}

func NewUserHandler(service services.UserServices) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}
func (user_handler *UserHandler) GetAllUser(ctx *gin.Context) {

}
func (user_handler *UserHandler) CreateUser(ctx *gin.Context) {

}
func (user_handler *UserHandler) UpdateUser(ctx *gin.Context) {

}
func (user_handler *UserHandler) DeleteUser(ctx *gin.Context) {

}
