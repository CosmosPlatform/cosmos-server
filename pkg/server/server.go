package server

import (
	"cosmos-server/pkg/routes"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewGinHandler(routes *routes.HTTPRoutes) *gin.Engine {
	e := gin.New()

	e.Use(loggingMiddleware(routes.Logger))
	e.Use(gin.Recovery())
	e.Use(corsMiddleware())
	e.Use(errorMiddleware(NewTranslator()))

	unauthenticatedGroup := e.Group("/")
	routes.RegisterUnauthenticatedRoutes(unauthenticatedGroup)

	adminAuthenticatedGroup := e.Group("/",
		authMiddleware(routes.AuthService),
		adminMiddleware(),
	)
	routes.RegisterAdminAuthenticatedRoutes(adminAuthenticatedGroup)

	return e
}

func StartServer(s *http.Server) error {
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
