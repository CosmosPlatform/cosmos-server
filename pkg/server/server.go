package server

import (
	"cosmos-server/pkg/routes/healthcheck"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewGinHandler() *gin.Engine {
	e := gin.New()

	healthcheck.AddHealthcheckHandler(e.Group("/"))

	return e
}

func StartServer(s *http.Server) error {
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
