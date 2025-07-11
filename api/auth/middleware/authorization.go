package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/molca-id/portal-app-api/common"
)

type authorizationProvider struct {
	network.ResponseSender
	common.ContextPayload
}

func NewAuthorizationProvider() network.AuthorizationProvider {
	return &authorizationProvider{
		ResponseSender: network.NewResponseSender(),
		ContextPayload: common.NewContextPayload(),
	}
}

func (m *authorizationProvider) Middleware(roleNames ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(roleNames) == 0 {
			m.Send(ctx).ForbiddenError("permission denied: role missing", nil)
			return
		}

		user := m.MustGetUser(ctx)

		log.Println(user)

		hasRole := true

		// TODO: check if user has any of the roles

		if !hasRole {
			m.Send(ctx).ForbiddenError("permission denied: does not have sufficient role", nil)
			return
		}

		ctx.Next()
	}
}
