package routes

import (
	"net/http"
	"photoset/internal/config"
	"photoset/internal/database"
	"photoset/internal/http/handlers"
	"photoset/internal/http/handlers/admin"
	"photoset/internal/http/middleware"
	"photoset/internal/repository"
	"photoset/internal/service"
	"photoset/internal/storage"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, cfg *config.Config) {
	stor, err := storage.NewStorage(&cfg.Storage)
	if err != nil {
		panic("存储初始化失败: " + err.Error())
	}

	healthHandler := handlers.NewHealthHandler()

	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 静态文件服务（付费图片需要签名验证）
	uploadsGroup := r.Group("/uploads", middleware.SignVerify(cfg))
	uploadsGroup.Any("/*path", gin.WrapH(http.StripPrefix("/uploads", http.FileServer(http.Dir("./uploads")))))

	r.GET("/api/health", healthHandler.Check)

	// 初始化服务和处理器
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	captchaService := service.NewCaptchaService()
	captchaHandler := handlers.NewCaptchaHandler(captchaService)
	siteSettingRepo := repository.NewSiteSettingRepository()
	authHandler := handlers.NewAuthHandler(userService, captchaService, siteSettingRepo)

	// 页面服务（新模块）
	pageRepo := repository.NewPageRepository(database.GetMySQL())
	pageService := service.NewPageService(pageRepo)
	pageHandler := handlers.NewPageHandler(pageService)

	photosetRepo := repository.NewPhotoSetRepository(database.GetMySQL())
	orderRepo := repository.NewOrderRepository(database.GetMySQL())
	photosetService := service.NewPhotoSetService(photosetRepo, orderRepo)
	photosetHandler := handlers.NewPhotoSetHandler(photosetService)
	tagHandler := handlers.NewTagHandler(photosetService)
	categoryHandler := handlers.NewCategoryHandler(photosetService)

	// 收藏路由
	favRepo := repository.NewFavoriteRepository(database.GetMySQL())
	favHandler := handlers.NewFavoriteHandler(favRepo)

	// 上传路由
	uploadHandler := handlers.NewUploadHandler(stor)

	// === 社区功能初始化 ===
	postRepo := repository.NewPostRepository(database.GetMySQL())
	replyRepo := repository.NewPostReplyRepository(database.GetMySQL())
	likeRepo := repository.NewPostLikeRepository(database.GetMySQL())
	replyLikeRepo := repository.NewPostReplyLikeRepository(database.GetMySQL())
	pointRepo := repository.NewUserPointRepository(database.GetMySQL())
	wordRepo := repository.NewSensitiveWordRepository(database.GetMySQL())
	reportRepo := repository.NewPostReportRepository(database.GetMySQL())
	categoryRepo := repository.NewPostCategoryRepository(database.GetMySQL())

	pointService := service.NewPointService(pointRepo)
	filterService := service.NewSensitiveFilterService(wordRepo)
	communityService := service.NewCommunityService(
		postRepo,
		replyRepo,
		likeRepo,
		replyLikeRepo,
		pointRepo,
		reportRepo,
		categoryRepo,
		pointService,
		filterService,
	)
	hotPostsService := service.NewHotPostsService(postRepo)

	communityHandler := handlers.NewCommunityHandler(
		database.GetMySQL(),
		communityService,
		pointService,
		hotPostsService,
	)
	adminCommunityHandler := admin.NewAdminCommunityHandler(database.GetMySQL())

	// 加载敏感词到内存
	service.InitSensitiveWords(wordRepo)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/captcha", middleware.CaptchaRateLimit(), captchaHandler.Generate)
			auth.POST("/register", middleware.RegisterRateLimit(), authHandler.Register)
			auth.POST("/login", middleware.LoginRateLimit(), authHandler.Login)
			auth.GET("/me", middleware.OptionalAuth(), authHandler.Me)
			auth.PUT("/password", middleware.Auth(), authHandler.ChangePassword)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPasswordByToken)
			auth.GET("/email-config", authHandler.CheckEmailConfig)
		}

		// 套图路由
		photosets := api.Group("/photosets")
		{
			photosets.GET("", middleware.OptionalAuth(), photosetHandler.List)
			photosets.GET("/advanced", middleware.OptionalAuth(), photosetHandler.AdvancedList)
			photosets.GET("/:id", middleware.OptionalAuth(), photosetHandler.Detail)
			photosets.POST("", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Create)
			photosets.PUT("/:id", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Update)
			photosets.DELETE("/:id", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Delete)
			photosets.GET("/:id/download", middleware.Auth(), photosetHandler.Download)
		}

		// 标签路由
		api.GET("/tags", tagHandler.List)

		// 分类公开路由
		api.GET("/categories", categoryHandler.List)

		// 收藏路由
		favorites := api.Group("/favorites")
		{
			favorites.Use(middleware.Auth())
			favorites.POST("/:photosetId", favHandler.Add)
			favorites.DELETE("/:photosetId", favHandler.Remove)
			favorites.GET("", favHandler.List)
		}

		// 评论路由
		commentRepo := repository.NewCommentRepository(database.GetMySQL())
		commentHandler := handlers.NewCommentHandler(commentRepo)
		comments := api.Group("/photosets/:id/comments")
		{
			comments.GET("", middleware.OptionalAuth(), commentHandler.List)
			comments.POST("", middleware.Auth(), commentHandler.Create)
			comments.GET("/:commentId/replies", middleware.OptionalAuth(), commentHandler.GetReplies)
			comments.DELETE("/:commentId", middleware.Auth(), commentHandler.Delete)
			comments.POST("/:commentId/like", middleware.Auth(), commentHandler.ToggleLike)
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
			orders.GET("", orderHandler.List)
			orders.POST("", orderHandler.Create)
			orders.POST("/:id/pay", orderHandler.Pay)
			orders.POST("/:id/refund", orderHandler.Refund)
		}

		// 管理后台路由（需 admin 权限）
		adminHandler := handlers.NewAdminHandler(photosetRepo, orderRepo, orderService)
		admin := api.Group("/admin")
		{
			admin.Use(middleware.Auth(), middleware.RequireRoles("admin"))
			admin.GET("/users", adminHandler.GetUsers)
			admin.GET("/users/export", adminHandler.ExportUsers)
			admin.GET("/users/:id", adminHandler.GetUserDetail)
			admin.GET("/photosets", adminHandler.GetPhotoSetsByStatus)
			admin.POST("/photosets/:id/approve", adminHandler.ApprovePhotoSet)
			admin.POST("/photosets/:id/reject", adminHandler.RejectPhotoSet)
			admin.POST("/photosets/batch/approve", adminHandler.BatchApprovePhotoSets)
			admin.POST("/photosets/batch/reject", adminHandler.BatchRejectPhotoSets)
			admin.POST("/photosets/batch/delete", adminHandler.BatchDeletePhotoSets)
			admin.PUT("/users/:id/ban", adminHandler.BanUser)
			admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
			admin.PUT("/users/:id/password", adminHandler.ResetUserPassword)
			admin.GET("/stats", adminHandler.Stats)
			admin.GET("/stats/trend", adminHandler.StatsTrend)
			admin.GET("/logs", adminHandler.GetAdminLogs)

			// 订单管理
			admin.GET("/orders", adminHandler.GetOrders)
			admin.GET("/orders/export", adminHandler.ExportOrders)
			admin.POST("/orders/:id/refund", adminHandler.AdminRefund)

			// 标签管理 CRUD
			admin.GET("/tags", tagHandler.AdminList)
			admin.POST("/tags", tagHandler.Create)
			admin.PUT("/tags/:id", tagHandler.Update)
			admin.DELETE("/tags/:id", tagHandler.Delete)

			// 分类管理 CRUD
			admin.GET("/categories", categoryHandler.AdminList)
			admin.POST("/categories", categoryHandler.Create)
			admin.PUT("/categories/:id", categoryHandler.Update)
			admin.DELETE("/categories/:id", categoryHandler.Delete)

			// 站点设置
			admin.GET("/settings", adminHandler.GetSettings)
			admin.PUT("/settings", adminHandler.UpdateSettings)
			// 系统管理
			admin.POST("/system/restart", adminHandler.RestartServer)
			// 邮件配置
			admin.POST("/mail/test-connection", adminHandler.TestMailConnection)
			admin.GET("/mail/config", adminHandler.GetMailConfig)
			admin.POST("/mail/send-test", adminHandler.SendTestMail)
			// 水印配置
			admin.GET("/watermark/info", adminHandler.GetWatermarkInfo)
			// 存储配置
			admin.POST("/storage/test", adminHandler.TestStorageConnection)
			admin.GET("/storage/status", adminHandler.GetStorageStatus)

			// 页面管理 CRUD
			admin.GET("/pages", pageHandler.AdminList)
			admin.POST("/pages", pageHandler.AdminCreate)
			admin.GET("/pages/:id", pageHandler.AdminGet)
			admin.PUT("/pages/:id", pageHandler.AdminUpdate)
			admin.DELETE("/pages/:id", pageHandler.AdminDelete)

			// 会员套餐管理 CRUD
			admin.GET("/memberships", membershipHandler.AdminList)
			admin.POST("/memberships", membershipHandler.AdminCreate)
			admin.PUT("/memberships/:id", membershipHandler.AdminUpdate)
			admin.DELETE("/memberships/:id", membershipHandler.AdminDelete)

			// 开发者中心
			admin.GET("/dev/api-keys", adminHandler.ListApiKeys)
			admin.POST("/dev/api-keys", adminHandler.CreateApiKey)
			admin.DELETE("/dev/api-keys/:id", adminHandler.DeleteApiKey)
			admin.GET("/dev/api-docs", adminHandler.GetApiDocs)
			admin.GET("/dev/sign-url-docs", adminHandler.GetSignUrlDocs)
		}

		// 公开路由 - 站点设置（不需要认证）
		api.GET("/settings", adminHandler.GetPublicSettings)

		// 公开页面路由
		api.GET("/pages/:slug", pageHandler.GetBySlug)
		api.GET("/pages", pageHandler.ListPublished)
	}

	// === 注册社区路由 ===
	RegisterCommunityRoutes(r, communityHandler, adminCommunityHandler)
}
