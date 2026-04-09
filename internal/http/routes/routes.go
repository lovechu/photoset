package routes

import (
	"photoset/internal/config"
	"photoset/internal/database"
	"photoset/internal/http/handlers"
	"photoset/internal/http/middleware"
	"photoset/internal/repository"
	"photoset/internal/service"
	"photoset/internal/storage"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	cfg := config.Load()

	stor, err := storage.NewStorage(&cfg.Storage)
	if err != nil {
		panic("存储初始化失败: " + err.Error())
	}

	healthHandler := handlers.NewHealthHandler()

	// 静态文件服务（付费图片需要签名验证）
	uploadsGroup := r.Group("/uploads", middleware.SignVerify())
	uploadsGroup.Static("/", "./uploads")

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	r.GET("/api/health", healthHandler.Check)

	// 初始化服务和处理器
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	captchaService := service.NewCaptchaService()
	captchaHandler := handlers.NewCaptchaHandler(captchaService)
	authHandler := handlers.NewAuthHandler(userService, captchaService)

	photosetRepo := repository.NewPhotoSetRepository(database.GetMySQL())
	orderRepo := repository.NewOrderRepository(database.GetMySQL())
	photosetService := service.NewPhotoSetService(photosetRepo, orderRepo)
	photosetHandler := handlers.NewPhotoSetHandler(photosetService)
	tagHandler := handlers.NewTagHandler(photosetService)

	// 收藏路由
	favRepo := repository.NewFavoriteRepository(database.GetMySQL())
	favHandler := handlers.NewFavoriteHandler(favRepo)

	// 上传路由
	uploadHandler := handlers.NewUploadHandler(stor)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/captcha", middleware.CaptchaRateLimit(), captchaHandler.Generate)
			auth.POST("/register", middleware.RegisterRateLimit(), authHandler.Register)
			auth.POST("/login", middleware.LoginRateLimit(), authHandler.Login)
			auth.GET("/me", middleware.Auth(), authHandler.Me)
		}

		// 套图路由
		photosets := api.Group("/photosets")
		{
			photosets.GET("", middleware.OptionalAuth(), photosetHandler.List)
			photosets.GET("/:id", middleware.OptionalAuth(), photosetHandler.Detail)
			photosets.POST("", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Create)
			photosets.PUT("/:id", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Update)
			photosets.DELETE("/:id", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Delete)
		}

		// 标签路由
		api.GET("/tags", tagHandler.List)

		// 收藏路由
		favorites := api.Group("/favorites")
		{
			favorites.Use(middleware.Auth())
			favorites.POST("/:photosetId", favHandler.Add)
			favorites.DELETE("/:photosetId", favHandler.Remove)
			favorites.GET("", favHandler.List)
		}

		// 上传路由
		upload := api.Group("/upload")
		{
			upload.Use(middleware.Auth(), middleware.RequireRoles("creator", "admin"))
			upload.POST("/image", uploadHandler.UploadImage)
		}

		// 用户路由
		api.GET("/users/profile", middleware.Auth(), authHandler.Me)

		// 会员套餐路由（公开接口）
		membershipRepo := repository.NewMembershipRepository(database.GetMySQL())
		membershipHandler := handlers.NewMembershipHandler(membershipRepo)
		api.GET("/memberships", membershipHandler.List)

		// 订单路由（需登录）
		orderService := service.NewOrderService(orderRepo, membershipRepo, photosetRepo)
		orderHandler := handlers.NewOrderHandler(orderService)
		orders := api.Group("/orders")
		{
			orders.Use(middleware.Auth())
			orders.POST("", orderHandler.Create)
			orders.POST("/:id/pay", orderHandler.Pay)
			orders.POST("/:id/refund", orderHandler.Refund)
			orders.GET("", orderHandler.List)
		}

		// 管理后台路由（需 admin 权限）
		adminHandler := handlers.NewAdminHandler(photosetRepo, orderRepo, orderService)
		admin := api.Group("/admin")
		{
			admin.Use(middleware.Auth(), middleware.RequireRoles("admin"))
			admin.GET("/users", adminHandler.GetUsers)
			admin.GET("/photosets", adminHandler.GetPhotoSetsByStatus)
			admin.POST("/photosets/:id/approve", adminHandler.ApprovePhotoSet)
			admin.POST("/photosets/:id/reject", adminHandler.RejectPhotoSet)
			admin.PUT("/users/:id/ban", adminHandler.BanUser)
			admin.GET("/stats", adminHandler.Stats)
			admin.POST("/orders/:id/refund", adminHandler.AdminRefund)
		}
	}
}

func CloseDB() error {
	database.CloseMySQL()
	database.CloseRedis()
	return nil
}
