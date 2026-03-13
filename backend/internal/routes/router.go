package routes

import (
	_ "backend/docs"
	"backend/internal/config"
	"backend/internal/graph"
	"backend/internal/routes/api"
	"backend/internal/routes/middlewares"
	"backend/internal/services"
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
		c.Set("services", svcs)
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

	r.GET("/ws", middlewares.JwtAuth(cfg), func(c *gin.Context) {
		svcs.Socket.ServeWS(c)
	})

	// GraphQL routes
	r.POST("/api/query", middlewares.OptionalJwtAuth(cfg), graph.GinContextToContextMiddleware(), api.GraphQLHandler(svcs))
	r.GET("/playground", api.PlaygroundHandler())

	// Auth routes
	authGroup := r.Group("/api/auth")
	api.RegisterAuthRoutes(authGroup, cfg, svcs)

	// Public analytics routes
	analyticsGroup := r.Group("/api/analytics")
	api.RegisterAnalyticsRoutes(analyticsGroup, cfg, svcs.Analytics)

	// Protected routes
	apiGroup := r.Group("/api")
	apiGroup.Use(middlewares.JwtAuth(cfg))
	{
		userGroup := apiGroup.Group("/user")
		api.RegisterUserRoutes(userGroup, cfg, svcs)

		adminGroup := apiGroup.Group("/admin")
		adminGroup.Use(middlewares.AdminMiddleware())
		api.RegisterAdminRoutes(adminGroup, cfg, svcs)
	}

	return r
}
