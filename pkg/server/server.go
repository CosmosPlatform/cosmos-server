package server

import (
	"cosmos-server/pkg/routes"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewGinHandler(routes *routes.HTTPRoutes) *gin.Engine {
	e := gin.New()

	unauthenticatedGroup := e.Group("/")
	routes.RegisterUnauthenticatedRoutes(unauthenticatedGroup)

	return e
}

func StartServer(s *http.Server) error {
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
