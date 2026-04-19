package handler

import (
	"fmt"
	"net/http"
	"url_shortener/internal/dto"
	"url_shortener/internal/services"
	"url_shortener/internal/utils"
	"url_shortener/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service services.UserServices
}

func NewUserHandler(service services.UserServices) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) Create(ctx *gin.Context) {
	var input dto.UserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userInput := input.MapCreateInputToModel()
	userData, err := uh.service.CreateUser(ctx, userInput)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapToUserDTO(userData)
	utils.ResponseSuccess(ctx, http.StatusOK, "Created user successfully!", userDTO)
}
func (uh *UserHandler) Update(ctx *gin.Context) {
	var input dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(input.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	var userInput dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userRole, _ := ctx.Get("user_role")
	if userRole == 1 && userInput.Role != nil && userInput.Status != nil {
		utils.ResponseError(ctx, fmt.Errorf("User does not have the permission to do this!"))
		return
	}
	userUpdateInput := userInput.MapUpdateInputToModel(userUuid)
	userUpdated, err := uh.service.UpdateUser(ctx, userUpdateInput)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapToUserDTO(userUpdated)
	utils.ResponseSuccess(ctx, http.StatusOK, "Updated user successfully!", userDTO)
}
func (uh *UserHandler) SoftDelteUser(ctx *gin.Context) {
	var input dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(input.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDeleted, err := uh.service.SoftDeleteUser(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapToUserDTO(userDeleted)
	utils.ResponseSuccess(ctx, http.StatusOK, "Deleted user successfully!", userDTO)
}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var input dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(input.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	err = uh.service.CleanSoftDelete(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}
func (uh *UserHandler) RestoreUser(ctx *gin.Context) {
	var input dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(input.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userData, err := uh.service.RestoreUser(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapToUserDTO(userData)
	utils.ResponseSuccess(ctx, http.StatusOK, "Restored user successfully!", userDTO)
}
func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	var input dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	userUuid, err := uuid.Parse(input.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userData, err := uh.service.GetUserByUUID(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := dto.MapToUserDTO(userData)
	utils.ResponseSuccess(ctx, http.StatusOK, "Fetched user data successfully!", userDTO)
}
func (uh *UserHandler) GetAllUser(ctx *gin.Context) {
	var lookupInput dto.GetSearchParams
	if err := ctx.ShouldBindQuery(&lookupInput); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	users, totalRecords, err := uh.service.GetAllUser(ctx, lookupInput.Search, lookupInput.Page, lookupInput.Limit, lookupInput.Order, lookupInput.Sort, false)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	usersDTO := dto.MapUsersToDTO(users)
	paginationResp := utils.NewPaginationResponse(usersDTO, lookupInput.Page, lookupInput.Limit, totalRecords)
	utils.ResponseSuccess(ctx, http.StatusOK, "Fetched all of users succesfully!", paginationResp)
}
func (uh *UserHandler) GetSoftDeleteUsers(ctx *gin.Context) {
	var lookupInput dto.GetSearchParams
	if err := ctx.ShouldBindQuery(&lookupInput); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	users, totalRecords, err := uh.service.GetAllUser(ctx, lookupInput.Search, lookupInput.Page, lookupInput.Limit, lookupInput.Order, lookupInput.Sort, true)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	usersDTO := dto.MapUsersToDTO(users)
	paginationResp := utils.NewPaginationResponse(usersDTO, lookupInput.Page, lookupInput.Limit, totalRecords)
	utils.ResponseSuccess(ctx, http.StatusOK, "Fetched all of delete users succesfully!", paginationResp)
}
