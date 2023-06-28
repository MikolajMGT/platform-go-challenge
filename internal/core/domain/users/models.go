package users_dm

import (
	"github.com/google/uuid"
	"time"
)

/*
 * UserEntity
 */

type User struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required,max=32" json:"password"`
}

type UserEntity struct {
	User
	Id         string    `validate:"required,uuid" json:"id"`
	CreateTime time.Time `validate:"required" json:"create_time"`
	UpdateTime time.Time `validate:"required" json:"update_time"`
}

func NewUserEntity() UserEntity {
	now := time.Now()

	return UserEntity{
		Id:         uuid.NewString(),
		CreateTime: now,
		UpdateTime: now,
	}
}
