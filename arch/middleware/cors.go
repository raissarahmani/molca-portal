package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/molca-id/portal-app-api/arch/network"
	"github.com/rs/cors"
)

type corsMiddleware struct {
	network.BaseMiddleware
	corsHandler *cors.Cors
}

func NewCorsMiddleware() network.RootMiddleware {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	return &corsMiddleware{
		BaseMiddleware: network.NewBaseMiddleware(),
		corsHandler:    c,
	}
}

func (m *corsMiddleware) Attach(engine *gin.Engine) {
	engine.Use(m.Handler)
}

func (m *corsMiddleware) Handler(ctx *gin.Context) {
	m.corsHandler.HandlerFunc(ctx.Writer, ctx.Request)

	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(204)
		return
	}

	ctx.Next()
}
