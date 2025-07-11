package middleware

import (
	"cosmos-server/pkg/log"
	"github.com/gin-gonic/gin"
	"time"
)

func LoggingMiddleware(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		logger.Infow("Request",
			"status", statusCode,
			"method", method,
			"path", path,
			"ip", clientIP,
			"latency", latency,
			"errors", c.Errors.String(),
		)
	}
}
