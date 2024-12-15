package typesgrpc

import (
	user_grpc "github.com/guluzadehh/go_bookstore/protos/gen/go/user"
	"github.com/guluzadehh/go_bookstore/services/user/internal/domain/models"
)

func NewUser(u *models.User) *user_grpc.User {
	return &user_grpc.User{
		Id:        u.Id,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format("2006-01-02"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02"),
		IsActive:  u.IsActive,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
	}
}
