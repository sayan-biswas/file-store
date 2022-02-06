package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sayan-biswas/file-store/pkg/server/config"
)

func CORS(config *config.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     config.CORS.AllowOrigins,
		AllowMethods:     config.CORS.AllowMethods,
		AllowHeaders:     config.CORS.AllowHeaders,
		ExposeHeaders:    config.CORS.ExposeHeaders,
		AllowCredentials: config.CORS.AllowCredentials,
		MaxAge:           time.Duration(config.CORS.MaxAge),
	})
}
