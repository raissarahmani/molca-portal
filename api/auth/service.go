package auth

import (
	"github.com/golang-jwt/jwt/v4"
	userModel "github.com/molca-id/portal-app-api/api/user/model"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/molca-id/portal-app-api/config"
)

type Service interface {
	VerifyToken(tokenString string, keyFunc jwt.Keyfunc) (*jwt.Token, error)
	ValidateClaims(claims *jwt.Token) (*userModel.User, bool)
}

type service struct {
	network.BaseService
	env *config.Env
}

func NewService(env *config.Env) Service {
	return &service{
		BaseService: network.NewBaseService(),
		env:         env,
	}
}

func (s *service) VerifyToken(tokenString string, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	return jwt.Parse(tokenString, keyFunc)
}

func (s *service) ValidateClaims(claims *jwt.Token) (*userModel.User, bool) {
	user, ok := claims.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}

	userID, ok := user["sub"].(string)
	if !ok {
		return nil, false
	}

	userObj, err := userModel.NewUser(userID)
	if err != nil {
		return nil, false
	}

	return userObj, true
}
