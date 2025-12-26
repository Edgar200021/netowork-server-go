package dto

import (
	"time"

	"github.com/Edgar200021/netowork-server-go/.gen/netowork/public/model"
)

type UserResponse struct {
	ID         string         `json:"id"`
	Email      string         `json:"email"`
	FirstName  string         `json:"firstName"`
	LastName   string         `json:"lastName"`
	Role       model.UserRole `json:"role"`
	CreatedAt  *time.Time     `json:"createdAt"`
	UpdatedAt  *time.Time     `json:"updatedAt"`
	Balance    int32          `json:"balance"`
	IsVerified bool           `json:"isVerified"`
}
