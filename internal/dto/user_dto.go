package dto

import (
	"url_shortener/internal/database/sqlc"

	"github.com/google/uuid"
)

type UserDTO struct {
	UUID      string `json:"user_uuid"`
	Email     string `json:"user_email"`
	Status    string `json:"user_status"`
	Role      string `json:"user_role"`
	CreatedAt string `json:"CreatedAt"`
	DeletedAt string `json:"DeletedAt"`
}
type UserInput struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email" binding:"required,email,email_advanced"`
	Password string `json:"password" binding:"required"`
	Status   int32  `json:"status" binding:"omitempty,oneof=1 2 3"`
	Role     int32  `json:"role" binding:"omitempty,oneof=1 2 3"`
}
type GetUserByUuidParam struct {
	Uuid string `uri:"uuid" binding:"uuid"`
}
type GetSearchParams struct {
	Search string `form:"search" binding:"omitempty,min=3,max=50,search"`
	Page   int32  `form:"page" binding:"omitempty,gte=1"`
	Limit  int32  `form:"limit" binding:"omitempty,gte=1,lte=500"`
	Order  string `form:"order_by" binding:"omitempty,oneof=user_id user_created_at"`
	Sort   string `form:"sort" binding:"omitempty,oneof=asc desc"`
}
type UpdateUserRequest struct {
	Password *string `json:"password" binding:"omitempty"`
	Status   *int32  `json:"status" binding:"omitempty,oneof=1 2 3"`
	Role     *int32  `json:"role" binding:"omitempty,oneof=1 2 3"`
}

func (input *UpdateUserRequest) MapUpdateInputToModel(userUUID uuid.UUID) sqlc.UpdateUserParams {
	return sqlc.UpdateUserParams{
		UserPassword: input.Password,
		UserRole:     input.Role,
		UserStatus:   input.Status,
		UserUuid:     userUUID,
	}
}

func MapToUserDTO(userData sqlc.User) *UserDTO {
	dto := &UserDTO{
		UUID:      userData.UserUuid.String(),
		Email:     userData.UserEmail,
		Status:    mapStatusText(int(userData.UserStatus)),
		Role:      mapRoleText(int(userData.UserRole)),
		CreatedAt: userData.UserCreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if userData.UserDeletedAt.Valid {
		dto.DeletedAt = userData.UserDeletedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	} else {
		dto.DeletedAt = ""
	}
	return dto
}
func MapUsersToDTO(usersData []sqlc.User) []UserDTO {
	dtos := make([]UserDTO, 0, len(usersData))
	for _, user := range usersData {
		dtos = append(dtos, *MapToUserDTO(user))
	}
	return dtos
}
func (input *UserInput) MapCreateInputToModel() sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		UserEmail:    input.Email,
		UserPassword: input.Password,
		UserStatus:   input.Status,
		UserRole:     input.Role,
	}
}
func mapStatusText(status int) string {
	switch status {
	case 1:
		return "Active"
	case 2:
		return "Inactive"
	case 3:
		return "Banned"
	default:
		return "None"
	}
}
func mapRoleText(status int) string {
	switch status {
	case 1:
		return "User"
	case 2:
		return "Moderator"
	case 3:
		return "Adminstrator"
	default:
		return "None"
	}
}
