package server

import (
	"cosmos-server/pkg/routes"
	"cosmos-server/pkg/server/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewGinHandler(routes *routes.HTTPRoutes) *gin.Engine {
	e := gin.New()

	e.Use(middleware.LoggingMiddleware(routes.Logger))
	e.Use(gin.Recovery())
	e.Use(middleware.CorsMiddleware())
	e.Use(middleware.ErrorMiddleware(middleware.NewTranslator()))

	unauthenticatedGroup := e.Group("/")
	routes.RegisterUnauthenticatedRoutes(unauthenticatedGroup)

	authenticatedGroup := e.Group("/",
		middleware.AuthMiddleware(routes.AuthService),
	)
	routes.RegisterAuthenticatedRoutes(authenticatedGroup)

	adminAuthenticatedGroup := e.Group("/",
		middleware.AuthMiddleware(routes.AuthService),
		middleware.AdminMiddleware(),
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
