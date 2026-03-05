package routes

import (
	"backend/internal/config"
	"backend/internal/routes/api"
	"backend/internal/routes/middlewares"
	"backend/internal/services"
	_ "backend/docs"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

func SetupRouter(cfg *config.Config, svcs *services.Services) *gin.Engine {
	r := gin.New()
	r.Use(middlewares.SlogMiddleware(), gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Update for prod UI
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes
	authGroup := r.Group("/api/auth")
	api.RegisterAuthRoutes(authGroup, cfg)

	// Protected routes
	apiGroup := r.Group("/api")
	apiGroup.Use(middlewares.JwtAuth(cfg))
	{
		propertyGroup := apiGroup.Group("/property")
		api.RegisterPropertyRoutes(propertyGroup, cfg, svcs.Property)

		adminGroup := apiGroup.Group("/admin")
		adminGroup.Use(middlewares.AdminMiddleware())
		api.RegisterAdminRoutes(adminGroup, cfg, svcs.Job)
	}

	return r
}
