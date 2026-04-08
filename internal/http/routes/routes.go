package routes

import (
	"photoset/internal/database"
	"photoset/internal/http/handlers"
	"photoset/internal/http/middleware"
	"photoset/internal/repository"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	healthHandler := handlers.NewHealthHandler()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	r.GET("/api/health", healthHandler.Check)

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	authHandler := handlers.NewAuthHandler(userService)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.Auth(), authHandler.Me)
		}

		// 用户路由（后续实现）
		// user := api.Group("/user")
		// {
		// 	user.GET("/profile", middleware.Auth(), userHandler.GetProfile)
		// }

		// 套图路由（后续实现）
		// photoset := api.Group("/photoset")
		// {
		// 	photoset.GET("/list", photosetHandler.List)
		// 	photoset.GET("/:id", photosetHandler.Detail)
		// }

		// 上传路由（后续实现）
		// upload := api.Group("/upload")
		// {
		// 	upload.POST("/image", middleware.Auth(), uploadHandler.UploadImage)
		// }

		// 订单路由（后续实现）
		// order := api.Group("/order")
		// {
		// 	order.POST("/create", middleware.Auth(), orderHandler.Create)
		// }

		// 管理后台路由（后续实现）
		// admin := api.Group("/admin")
		// {
		// 	admin.Use(middleware.Auth())
		// 	admin.GET("/users", adminHandler.GetUsers)
		// }
	}
}

func CloseDB() error {
	database.CloseMySQL()
	database.CloseRedis()
	return nil
}


