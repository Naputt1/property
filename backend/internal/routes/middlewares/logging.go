package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func SlogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		slog.Info("Request",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"ip", c.ClientIP(),
			"latency", time.Since(start).String(),
			"user-agent", c.Request.UserAgent(),
			"errors", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}
}
