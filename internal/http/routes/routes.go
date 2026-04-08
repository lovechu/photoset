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

	// 初始化服务和处理器
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	authHandler := handlers.NewAuthHandler(userService)

	photosetRepo := repository.NewPhotoSetRepository(database.GetMySQL())
	photosetService := service.NewPhotoSetService(photosetRepo)
	photosetHandler := handlers.NewPhotoSetHandler(photosetService)
	tagHandler := handlers.NewTagHandler(photosetService)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.Auth(), authHandler.Me)
		}

		// 套图路由
		photosets := api.Group("/photosets")
		{
			photosets.GET("", photosetHandler.List)
			photosets.GET("/:id", photosetHandler.Detail)
			photosets.POST("", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Create)
		}

		// 标签路由
		api.GET("/tags", tagHandler.List)

		// 用户路由（后续实现）
		// user := api.Group("/user")
		// {
		// 	user.GET("/profile", middleware.Auth(), userHandler.GetProfile)
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


