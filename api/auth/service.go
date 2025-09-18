package auth

import (
	"errors"
	"log"
	"time"

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

type CustomClaims struct {
	jwt.RegisteredClaims
}

const tokenLeeway = 10 * time.Second

func (c *CustomClaims) Valid() error {
	now := time.Now().UTC()

	if c.ExpiresAt != nil {
		if now.Add(-tokenLeeway).After(c.ExpiresAt.Time) {
			return errors.New("token is expired")
		}
	}

	if c.NotBefore != nil {
		if c.NotBefore.Time.After(now.Add(tokenLeeway)) {
			return errors.New("token not valid yet (nbf)")
		}
	}

	if c.IssuedAt != nil {
		if c.IssuedAt.Time.After(now.Add(tokenLeeway)) {
			return errors.New("token used before issued (iat)")
		}
	}

	return nil
}

func (s *service) VerifyToken(tokenString string, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err == nil && token.Valid {
		return token, nil
	}

	var rawClaims jwt.MapClaims
	if _, _, perr := new(jwt.Parser).ParseUnverified(tokenString, &rawClaims); perr == nil {
		if v, ok := rawClaims["iat"].(float64); ok {
			log.Printf("[DEBUG] token iat: %v (UTC)", time.Unix(int64(v), 0).UTC())
		}
		if v, ok := rawClaims["nbf"].(float64); ok {
			log.Printf("[DEBUG] token nbf: %v (UTC)", time.Unix(int64(v), 0).UTC())
		}
		if v, ok := rawClaims["exp"].(float64); ok {
			log.Printf("[DEBUG] token exp: %v (UTC)", time.Unix(int64(v), 0).UTC())
		}
		log.Printf("[DEBUG] server now: %v (UTC)", time.Now().UTC())
	}

	return nil, err
}

func (s *service) ValidateClaims(token *jwt.Token) (*userModel.User, bool) {
	if c, ok := token.Claims.(*CustomClaims); ok && c != nil {
		if c.Subject == "" {
			return nil, false
		}
		userObj, err := userModel.NewUser(c.Subject)
		if err != nil {
			return nil, false
		}
		return userObj, true
	}

	if mc, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := mc["sub"].(string); ok && sub != "" {
			userObj, err := userModel.NewUser(sub)
			if err != nil {
				return nil, false
			}
			return userObj, true
		}
	}

	return nil, false
}
