package middlewhere

import (
	"bankService/internal/logger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := fmt.Sprintf("%x", time.Now().UnixNano())
		logg := logger.GetLogger().WithField("request_id", requestID)

		logg.WithFields(logrus.Fields{
			"method":        c.Request.Method,
			"url":           c.Request.URL.String(),
			"time":          startTime.Format(time.RFC3339),
			"user_agent":    c.Request.UserAgent(),
			"remote_ip":     c.ClientIP(),
			"request_size":  c.Request.ContentLength,
			"response_size": int64(c.Writer.Size()),
			"headers":       c.Request.Header,
		}).Debug("Request started")

		c.Next()

		logg.WithFields(logrus.Fields{
			"method":        c.Request.Method,
			"url":           c.Request.URL.String(),
			"status_code":   c.Writer.Status(),
			"latency":       time.Since(startTime),
			"user_agent":    c.Request.UserAgent(),
			"remote_ip":     c.ClientIP(),
			"request_size":  c.Request.ContentLength,
			"response_size": int64(c.Writer.Size()),
			"headers":       c.Request.Header,
		}).Debug("Request completed")
	}
}
