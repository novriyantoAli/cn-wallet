package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	userSecurityDto "github.com/novriyantoAli/cn-wallet/internal/application/user-security/dto"
	userSecurityService "github.com/novriyantoAli/cn-wallet/internal/application/user-security/service"
	jwt "github.com/novriyantoAli/cn-wallet/internal/pkg/jwt"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
		)
	}
}

func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal domain error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-PIN")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func JWTMiddleware(
	jwtManager *jwt.JWTManager,
	logger *zap.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Authorization header missing",
			})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid authorization header format",
			})
			return
		}

		token := authHeader[len(bearerPrefix):]
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Empty token in authorization header",
			})
			return
		}

		_, err := jwtManager.VerifyToken(token)
		if err != nil {
			logger.Error("Failed to verify token", zap.Error(err))
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid token",
			})
			return
		}

		c.Next()
	}
}

func PINMiddleware(
	userSecurity userSecurityService.UserSecurityService,
	jwtManager *jwt.JWTManager,
	logger *zap.Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// get authorization from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Authorization header missing",
			})
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid authorization header format",
			})
			return
		}

		token := authHeader[len(bearerPrefix):]
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Empty token in authorization header",
			})
			return
		}

		claims, err := jwtManager.VerifyToken(token)
		if err != nil {
			logger.Error("Failed to verify token", zap.Error(err))
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Invalid token",
			})
			return
		}

		pin := c.GetHeader("X-PIN")
		if pin == "" {
			c.AbortWithStatusJSON(403, gin.H{
				"error": "PIN required",
			})
			return
		}

		b, err := userSecurity.VerifyPIN(c, &userSecurityDto.VerifyPINRequest{
			UserID: claims.UserID,
			PIN:    pin,
		})
		if err != nil {
			logger.Error("Failed to verify PIN", zap.Error(err))
			c.AbortWithStatusJSON(403, gin.H{
				"error": "Invalid PIN",
			})
			return
		}
		if !b {
			c.AbortWithStatusJSON(403, gin.H{
				"error": "Invalid PIN",
			})
			return
		}

		c.Next()
	}
}
