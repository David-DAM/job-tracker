package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TraceIDMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		c.Set("traceId", traceID)

		c.Writer.Header().Set("X-Trace-ID", traceID)

		logger.Info("incoming request",
			zap.String("traceId", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		c.Next()

		status := c.Writer.Status()
		logger.Info("request completed",
			zap.String("traceId", traceID),
			zap.Int("status", status),
		)
	}
}
