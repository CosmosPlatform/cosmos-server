package healthcheck

import "github.com/gin-gonic/gin"

func AddHealthcheckHandler(e *gin.RouterGroup) {
	e.GET("/healthcheck", handleHealthcheck)
}

func handleHealthcheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
