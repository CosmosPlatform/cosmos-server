package healthcheck

import (
	"github.com/gin-gonic/gin"
)

type handler struct{}

func AddHealthcheckHandler(e *gin.RouterGroup) {
	handler := &handler{}

	e.GET("/healthcheck", handler.handleHealthcheck)
}

func (handler *handler) handleHealthcheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
