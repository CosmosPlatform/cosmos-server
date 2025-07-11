package middleware

import "github.com/gin-gonic/gin"

func ErrorMiddleware(translator Translator) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if callError := c.Errors.Last(); callError != nil {
			apiError := translator.ToApiError(callError.Err)
			c.JSON(apiError.StatusCode, apiError)
		}
	}
}
