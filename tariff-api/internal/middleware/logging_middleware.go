package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := logrus.Fields{
			"status":   status,
			"method":   c.Request.Method,
			"path":     path,
			"clientIP": c.ClientIP(),
			"latency":  latency.String(),
		}
		if rawQuery != "" {
			fields["query"] = rawQuery
		}

		entry := log.WithFields(fields)

		if len(c.Errors) > 0 {
			entry = entry.WithField("errors", c.Errors.String())
			entry.Error("request completed with errors")
			return
		}

		entry.Info("request completed")
	}
}
