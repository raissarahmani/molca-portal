package common

import (
	"errors"

	"github.com/gin-gonic/gin"
	userModel "github.com/molca-id/portal-app-api/api/user/model"
)

const (
	payloadUser string = "user"
)

type ContextPayload interface {
	SetUser(ctx *gin.Context, value *userModel.User)
	MustGetUser(ctx *gin.Context) *userModel.User
}

type payload struct{}

func NewContextPayload() ContextPayload {
	return &payload{}
}

func (u *payload) SetUser(ctx *gin.Context, value *userModel.User) {
	ctx.Set(payloadUser, value)
}

func (u *payload) MustGetUser(ctx *gin.Context) *userModel.User {
	value, ok := ctx.MustGet(payloadUser).(*userModel.User)
	if !ok {
		panic(errors.New(payloadUser + " missing for context"))
	}
	return value
}
