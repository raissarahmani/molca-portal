package middleware

import (
	"log"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	auth "github.com/molca-id/portal-app-api/api/auth"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/molca-id/portal-app-api/common"
	"github.com/molca-id/portal-app-api/config"
	"github.com/molca-id/portal-app-api/utils"
)

type authenticationProvider struct {
	network.ResponseSender
	common.ContextPayload
	env         *config.Env
	authService auth.Service
}

func NewAuthenticationProvider(env *config.Env, authService auth.Service) network.AuthenticationProvider {
	return &authenticationProvider{
		ResponseSender: network.NewResponseSender(),
		ContextPayload: common.NewContextPayload(),
		env:            env,
		authService:    authService,
	}
}

var jwks *keyfunc.JWKS

func initJWKs(env *config.Env) {

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("[RECOVER] Failed to fetch JWKS: %v", r)
		}
	}()

	var err error
	jwks, err = keyfunc.Get(env.JWKsURL, keyfunc.Options{
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshTimeout:    time.Second * 10,
		RefreshUnknownKID: true,
	})
	if err != nil {
		log.Fatalf("[ERROR] Failed to fetch JWKS: %v", err)
		panic(err)
	}
	log.Println("[INFO] JWKS fetched successfully")
}

func (m *authenticationProvider) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if jwks == nil {
			initJWKs(m.env)
		}

		authHeader := ctx.GetHeader(network.AuthorizationHeader)
		if len(authHeader) == 0 {
			m.Send(ctx).UnauthorizedError("permission denied: missing authorization", nil)
			log.Printf("[DEBUG] Authorization header: %s", authHeader)
			return
		}

		tokenString := utils.ExtractBearerToken(authHeader)
		if tokenString == "" {
			m.Send(ctx).UnauthorizedError("permission denied: invalid authorization", nil)
			return
		}

		log.Printf("[DEBUG] Raw token (first 50 chars): %s...", tokenString[:50])

		// Parse the token
		token, err := m.authService.VerifyToken(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			m.Send(ctx).UnauthorizedError("permission denied: invalid authorization", nil)
			log.Printf("[ERROR] Token verification failed: %v", err)
			return
		}

		// Get the claims
		claims, ok := m.authService.ValidateClaims(token)
		if !ok {
			m.Send(ctx).UnauthorizedError("permission denied: invalid authorization", nil)
			log.Printf("[ERROR] Token claims invalid: %+v", token.Claims)
			return
		}

		m.SetUser(ctx, claims)

		ctx.Next()
	}
}
