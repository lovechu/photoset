package routes

import (
	"net/http"
	"photoset/internal/config"
	"photoset/internal/database"
	"photoset/internal/http/handlers"
	"photoset/internal/http/handlers/admin"
	"photoset/internal/http/middleware"
	"photoset/internal/logger"
	"photoset/internal/repository"
	"photoset/internal/service"
	"photoset/internal/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "photoset/docs" // Swagger generated docs
)

func Setup(r *gin.Engine, cfg *config.Config) {
	stor, err := storage.NewStorage(&cfg.Storage)
	if err != nil {
		panic("存储初始化失败: " + err.Error())
	}

	healthHandler := handlers.NewHealthHandler()
	performanceHandler := handlers.NewPerformanceHandler()

	r.Use(middleware.CORS())
	r.Use(middleware.RequestID())
	r.Use(middleware.Metrics()) // 性能指标追踪
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 静态文件服务（付费图片需要签名验证）
	uploadsGroup := r.Group("/uploads", middleware.SignVerify(cfg))
	uploadsGroup.Any("/*path", gin.WrapH(http.StripPrefix("/uploads", http.FileServer(http.Dir("./uploads")))))

	r.GET("/api/health", healthHandler.Check)
	r.GET("/api/v1/health", healthHandler.Check) // 兼容新版本

	// Swagger API 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 初始化服务和处理器
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	captchaService := service.NewCaptchaService()
	captchaHandler := handlers.NewCaptchaHandler(captchaService)
	siteSettingRepo := repository.NewSiteSettingRepository()

	// 邮箱验证码服务（需要在 AuthHandler 之前初始化，因为 AuthHandler 依赖它）
	emailVerificationSvc := service.NewEmailVerificationService(siteSettingRepo)

	// 登录历史服务（需要在 AuthHandler 之前初始化）
	loginHistoryRepo := repository.NewLoginHistoryRepository(database.GetMySQL())
	loginHistoryService := service.NewLoginHistoryService(loginHistoryRepo)
	loginHistoryHandler := handlers.NewLoginHistoryHandler(loginHistoryService)

	// 用户设备管理服务（需要在 AuthHandler 之前初始化）
	userDeviceRepo := repository.NewUserDeviceRepository(database.GetMySQL())
	userDeviceService := service.NewUserDeviceService(userDeviceRepo)
	userDeviceHandler := handlers.NewUserDeviceHandler(userDeviceService)

	authHandler := handlers.NewAuthHandler(userService, captchaService, siteSettingRepo, loginHistoryService, userDeviceService, emailVerificationSvc)

	// 页面服务（新模块）
	pageRepo := repository.NewPageRepository(database.GetMySQL())
	pageService := service.NewPageService(pageRepo)
	pageHandler := handlers.NewPageHandler(pageService)

	photosetRepo := repository.NewPhotoSetRepository(database.GetMySQL())
	orderRepo := repository.NewOrderRepository(database.GetMySQL())
	photosetService := service.NewPhotoSetService(photosetRepo, orderRepo, cfg, stor)
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
	draftRepo := repository.NewDraftRepository(database.GetMySQL())
	shareRepo := repository.NewPostShareRepository(database.GetMySQL())
	tagRepo := repository.NewTagRepository(database.GetMySQL())
	postTagRepo := repository.NewPostTagRepository(database.GetMySQL())
	topicRepo := repository.NewTopicRepository(database.GetMySQL())
	postTopicRepo := repository.NewPostTopicRepository(database.GetMySQL())
	notificationRepo := repository.NewNotificationRepository(database.GetMySQL())

	pointService := service.NewPointService(pointRepo)
	filterService := service.NewSensitiveFilterService(wordRepo)
	mentionService := service.NewMentionService(userRepo, notificationRepo)
	communityService := service.NewCommunityService(
		postRepo,
		replyRepo,
		likeRepo,
		replyLikeRepo,
		shareRepo,
		pointRepo,
		reportRepo,
		categoryRepo,
		draftRepo,
		tagRepo,
		postTagRepo,
		topicRepo,
		postTopicRepo,
		pointService,
		filterService,
		mentionService,
		stor,
	)
	hotPostsService := service.NewHotPostsService(postRepo)
	recommendationService := service.NewRecommendationService(
		database.GetMySQL(),
		postRepo,
		likeRepo,
		shareRepo,
		replyRepo,
		postTagRepo,
		postTopicRepo,
		tagRepo,
		topicRepo,
		hotPostsService,
	)

	// IP地理位置服务（用于社区发帖/回帖时实时解析IP）
	ipGeoService := service.NewIPGeoService()

	communityHandler := handlers.NewCommunityHandler(
		database.GetMySQL(),
		communityService,
		pointService,
		hotPostsService,
		recommendationService,
		ipGeoService,
	)
	adminCommunityHandler := admin.NewAdminCommunityHandler(database.GetMySQL(), stor)

	// WebSocket 实时消息（必须在通知和私信功能之前初始化）
	wsHub := service.NewHub()
	go wsHub.Run()
	wsHandler := handlers.NewWebSocketHandler(wsHub)

	// 通知功能
	notificationService := service.NewNotificationService(notificationRepo, wsHub)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	// 推送通知功能
	pushTokenRepo := repository.NewPushTokenRepository(database.GetMySQL())
	pushNotificationService := service.NewPushNotificationService(pushTokenRepo)
	pushNotificationHandler := handlers.NewPushNotificationHandler(pushNotificationService)
	
	// 将推送服务注入到通知服务中
	notificationService.SetPushService(pushNotificationService)

	// 私信功能
	messageRepo := repository.NewMessageRepository(database.GetMySQL())
	messageService := service.NewMessageService(messageRepo, userRepo)
	messageHandler := handlers.NewMessageHandler(messageService, wsHub)

	// Follow 功能
	followRepo := repository.NewFollowRepository(database.GetMySQL())
	followService := service.NewFollowService(followRepo, userRepo)
	followHandler := handlers.NewFollowHandler(followService)

	// 用户等级系统
	levelRepo := repository.NewUserLevelRepository(database.GetMySQL())
	achievementRepo := repository.NewAchievementRepository(database.GetMySQL())
	pointsMallRepo := repository.NewPointsMallRepository(database.GetMySQL())
	userLevelService := service.NewUserLevelService(levelRepo, achievementRepo, pointsMallRepo, pointRepo, database.GetMySQL())
	userLevelService.Initialize() // 初始化默认数据
	userLevelHandler := handlers.NewUserLevelHandler(userLevelService)
	adminLevelHandler := admin.NewAdminLevelHandler(database.GetMySQL())

	// 加载敏感词到内存
	service.InitSensitiveWords(wordRepo)

// 用户拉黑/屏蔽功能
blockRepo := repository.NewUserBlockRepository(database.GetMySQL())
blockService := service.NewUserBlockService(blockRepo)
blockHandler := handlers.NewUserBlockHandler(blockService)

// 用户隐私设置功能
privacyRepo := repository.NewUserPrivacyRepository(database.GetMySQL())
privacyService := service.NewUserPrivacyService(privacyRepo)
privacyHandler := handlers.NewUserPrivacyHandler(privacyService)

// 浏览历史功能
viewHistoryRepo := repository.NewViewHistoryRepository(database.GetMySQL())
viewHistoryService := service.NewViewHistoryService(viewHistoryRepo)
viewHistoryHandler := handlers.NewViewHistoryHandler(viewHistoryService)

	// 账号注销功能
	accountDeletionService := service.NewAccountDeletionService()
	accountDeletionHandler := handlers.NewAccountDeletionHandler(accountDeletionService)

	// 用户反馈功能
	feedbackHandler := handlers.NewFeedbackHandler()

	// 分享链接功能
	shareHandler := handlers.NewShareHandler(photosetRepo, cfg)

	// 探索/发现页功能
	exploreService := service.NewExploreService(database.GetMySQL(), photosetRepo)
	exploreHandler := handlers.NewExploreHandler(exploreService)

	// 创作者数据统计功能
	creatorStatsService := service.NewCreatorStatsService(database.GetMySQL())
	creatorStatsHandler := handlers.NewCreatorStatsHandler(creatorStatsService)

	// 套图评价功能
	reviewRepo := repository.NewPhotoSetReviewRepository(database.GetMySQL())
	reviewService := service.NewReviewService(reviewRepo, photosetRepo)
	reviewHandler := handlers.NewReviewHandler(reviewService)

// OAuth2 第三方登录授权系统
	oauthClientRepo := repository.NewOAuthClientRepository()
	oauthAuthorizationRepo := repository.NewOAuthAuthorizationRepository()
	oauthTokenRepo := repository.NewOAuthTokenRepository()
	oauthService := service.NewOAuthService(oauthClientRepo, oauthAuthorizationRepo, oauthTokenRepo, userRepo)
	oauthHandler := handlers.NewOAuthHandler(oauthService)
	adminOAuthHandler := admin.NewAdminOAuthHandler(oauthService)

	// 评论路由（需要在 api 和 v1 中使用）
	commentRepo := repository.NewCommentRepository(database.GetMySQL())
	commentHandler := handlers.NewCommentHandler(commentRepo)

	// 合集路由（需要在 api 和 v1 中使用）
	collectionRepo := repository.NewUserCollectionRepository(database.GetMySQL())
	collectionService := service.NewCollectionService(collectionRepo)
	collectionHandler := handlers.NewCollectionHandler(collectionService)

	// 会员套餐路由（需要在 api 和 v1 中使用）
	membershipRepo := repository.NewMembershipRepository(database.GetMySQL())
	membershipHandler := handlers.NewMembershipHandler(membershipRepo)

	// 订单路由（需要在 api 和 v1 中使用）
	orderService := service.NewOrderService(orderRepo, membershipRepo, photosetRepo)

	// 支付宝服务（可选，配置不完整时为 nil）
	var alipayService *service.AlipayService
	if cfg.Alipay.AppID != "" {
		var err error
		alipayService, err = service.NewAlipayService(&cfg.Alipay, orderRepo)
		if err != nil {
			logger.Warn("支付宝服务初始化失败，支付宝支付将不可用", "error", err)
		}
	}

	// 微信支付服务（可选，配置不完整时仍可创建）
	wechatPayService := service.NewWechatPayService(&cfg.WechatPay, orderRepo)

	// 从数据库加载支付配置（如果存在则覆盖环境变量配置）
	settings, _ := siteSettingRepo.GetAll()
	if settings != nil {
		strSettings := make(map[string]string)
		for k, v := range settings {
			strSettings[k] = v
		}
		if alipayService != nil {
			alipayService.ReloadFromSettings(strSettings)
		}
		wechatPayService.ReloadFromSettings(strSettings)
	}

	orderHandler := handlers.NewOrderHandler(orderService, alipayService)

	// 管理后台路由（需要在 api 和 v1 中使用）
	adminHandler := handlers.NewAdminHandler(photosetRepo, orderRepo, orderService, cfg, alipayService, wechatPayService, database.GetMySQL())
	systemHandler := admin.NewSystemHandler()
	backupService := service.NewBackupService(cfg)
	backupHandler := admin.NewBackupHandler(backupService)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/captcha", middleware.CaptchaRateLimit(), captchaHandler.Generate)
			auth.POST("/register", middleware.RegisterRateLimit(), authHandler.Register)
			auth.POST("/login", middleware.LoginRateLimit(), authHandler.Login)
			auth.GET("/me", middleware.OptionalAuth(), authHandler.Me)
			auth.PUT("/password", middleware.Auth(), authHandler.ChangePassword)
			auth.PUT("/profile", middleware.Auth(), authHandler.UpdateProfile)
			auth.POST("/logout", middleware.Auth(), authHandler.Logout)
			auth.POST("/forgot-password", middleware.ForgotPasswordRateLimit(), authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPasswordByToken)
			auth.GET("/email-config", authHandler.CheckEmailConfig)
			// 邮箱验证码和绑定邮箱
			auth.POST("/send-email-code", authHandler.SendEmailCode)
			auth.POST("/verify-email-code", authHandler.VerifyEmailCode)
			auth.PUT("/bind-email", middleware.Auth(), authHandler.BindEmail)
			// 账号注销
			auth.POST("/request-deletion", middleware.Auth(), accountDeletionHandler.RequestDeletion)
			auth.POST("/cancel-deletion", middleware.Auth(), accountDeletionHandler.CancelDeletion)
			auth.GET("/deletion-status", middleware.Auth(), accountDeletionHandler.GetDeletionStatus)
		}

		// 套图路由
		photosets := api.Group("/photosets")
		{
			photosets.GET("", middleware.OptionalAuth(), photosetHandler.List)
			photosets.GET("/advanced", middleware.OptionalAuth(), photosetHandler.AdvancedList)
			photosets.GET("/trash", middleware.Auth(), photosetHandler.GetTrash)
			photosets.GET("/:id", middleware.OptionalAuth(), photosetHandler.Detail)
			photosets.POST("", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Create)
			photosets.PUT("/:id", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Update)
			photosets.DELETE("/:id", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Delete)
			photosets.POST("/:id/restore", middleware.Auth(), photosetHandler.Restore)
			photosets.DELETE("/:id/permanent", middleware.Auth(), photosetHandler.PermanentDelete)
			photosets.GET("/:id/download", middleware.Auth(), photosetHandler.Download)
			// 分享链接
			photosets.GET("/:id/share-link", shareHandler.GenerateShareLink)
			// 套图评价
			photosets.GET("/:id/reviews", middleware.OptionalAuth(), reviewHandler.List)
			photosets.GET("/:id/reviews/summary", middleware.OptionalAuth(), reviewHandler.GetSummary)
			photosets.GET("/:id/reviews/mine", middleware.Auth(), reviewHandler.GetMyReview)
			photosets.POST("/:id/reviews", middleware.Auth(), reviewHandler.Create)
			photosets.PUT("/:id/reviews/:reviewId", middleware.Auth(), reviewHandler.Update)
			photosets.DELETE("/:id/reviews/:reviewId", middleware.Auth(), reviewHandler.Delete)
		}

		// 分享链接验证（公开）
		api.GET("/share/:id", shareHandler.VerifyShareLink)

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
			favorites.POST("/batch-remove", favHandler.BatchRemove)
		}

		// 评论路由
		comments := api.Group("/photosets/:id/comments")
		{
			comments.GET("", middleware.OptionalAuth(), commentHandler.List)
			comments.POST("", middleware.Auth(), commentHandler.Create)
			comments.GET("/:commentId/replies", middleware.OptionalAuth(), commentHandler.GetReplies)
			comments.DELETE("/:commentId", middleware.Auth(), commentHandler.Delete)
			comments.POST("/:commentId/like", middleware.Auth(), commentHandler.ToggleLike)
		}

		// 合集路由
		collections := api.Group("/collections")
		{
			collections.Use(middleware.Auth())
			collections.POST("", collectionHandler.Create)
			collections.GET("", collectionHandler.List)
			collections.GET("/:id", collectionHandler.Get)
			collections.PUT("/:id", collectionHandler.Update)
			collections.DELETE("/:id", collectionHandler.Delete)
			collections.POST("/:id/items", collectionHandler.AddItem)
			collections.DELETE("/:id/items/:photosetId", collectionHandler.RemoveItem)
			collections.POST("/:id/items/batch", collectionHandler.BatchAddItems)
			collections.GET("/by-photoset/:photosetId", collectionHandler.GetCollectionsForPhotoset)
		}

		// 上传路由
		upload := api.Group("/upload")
		{
			upload.Use(middleware.Auth(), middleware.RequireRoles("creator", "admin"))
			upload.POST("/image", uploadHandler.UploadImage)
		}

		// 社区上传路由（所有登录用户可用）
		communityUpload := api.Group("/community/upload")
		{
			communityUpload.Use(middleware.Auth())
			communityUpload.POST("/image", uploadHandler.UploadImage)
			communityUpload.POST("/video", uploadHandler.UploadVideo)
		}

		// 头像上传路由（所有登录用户可用）
		api.POST("/user/avatar", middleware.Auth(), uploadHandler.UploadImage)

		// 用户路由
		api.GET("/users/profile", middleware.Auth(), authHandler.Me)

		// 用户隐私设置路由
		api.GET("/user/privacy-settings", middleware.Auth(), privacyHandler.GetPrivacySettings)
		api.PUT("/user/privacy-settings", middleware.Auth(), privacyHandler.UpdatePrivacySettings)

		// 登录历史路由
		api.GET("/user/login-history", middleware.Auth(), loginHistoryHandler.GetLoginHistory)

		// 用户设备管理路由
		api.GET("/user/devices", middleware.Auth(), userDeviceHandler.GetUserDevices)
		api.DELETE("/user/devices/:deviceId", middleware.Auth(), userDeviceHandler.DeactivateDevice)
		api.DELETE("/user/devices", middleware.Auth(), userDeviceHandler.DeactivateOtherDevices)

		// 浏览历史路由
		api.POST("/user/view-history", middleware.Auth(), viewHistoryHandler.Record)
		api.GET("/user/view-history", middleware.Auth(), viewHistoryHandler.List)
		api.DELETE("/user/view-history/:id", middleware.Auth(), viewHistoryHandler.Delete)
		api.DELETE("/user/view-history", middleware.Auth(), viewHistoryHandler.ClearAll)
		api.POST("/user/view-history/batch-delete", middleware.Auth(), viewHistoryHandler.BatchDelete)

		// 推送通知路由
		api.POST("/user/push-tokens", middleware.Auth(), pushNotificationHandler.RegisterToken)
		api.DELETE("/user/push-tokens", middleware.Auth(), pushNotificationHandler.UnregisterToken)

		// 用户反馈路由
		api.POST("/feedback", middleware.OptionalAuth(), feedbackHandler.Create)
		api.POST("/feedback/image", middleware.OptionalAuth(), feedbackHandler.UploadFeedbackImage)

		// 会员套餐路由（公开接口）
		api.GET("/memberships", membershipHandler.List)

		// 订单路由（需登录）
		orders := api.Group("/orders")
		{
			orders.Use(middleware.Auth())
			orders.GET("", orderHandler.List)
			orders.POST("", orderHandler.Create)
			orders.POST("/:id/pay", orderHandler.Pay)
			orders.POST("/:id/refund", orderHandler.Refund)
		}

		// 支付宝回调（无需认证，由支付宝服务器调用）
		api.POST("/payment/alipay/notify", orderHandler.AlipayNotify)

		// 管理后台路由（需 admin 权限）
		admin := api.Group("/admin")
		{
			admin.Use(middleware.Auth(), middleware.RequireRoles("admin"))
			admin.GET("/users", adminHandler.GetUsers)
			admin.POST("/users", adminHandler.CreateUser)
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

			// 套图导出
			admin.GET("/photosets/export", adminHandler.ExportPhotoSets)

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
			admin.POST("/system/restart", middleware.RestartRateLimit(), adminHandler.RestartServer)
			// 性能分析（管理员）
			admin.GET("/system/metrics", performanceHandler.GetMetrics)
			admin.GET("/system/goroutines", performanceHandler.GetGoroutines)
			admin.GET("/system/profile/cpu", performanceHandler.StartCPUProfile)
			admin.GET("/system/profile/heap", performanceHandler.GetHeapProfile)
			// 邮件配置
			admin.POST("/mail/test-connection", adminHandler.TestMailConnection)
			admin.GET("/mail/config", adminHandler.GetMailConfig)
			admin.POST("/mail/send-test", adminHandler.SendTestMail)
			// 水印配置
			admin.GET("/watermark/info", adminHandler.GetWatermarkInfo)
			// 存储配置
			admin.POST("/storage/test", adminHandler.TestStorageConnection)
			admin.GET("/storage/status", adminHandler.GetStorageStatus)

			// 支付配置
			admin.POST("/payment/alipay/test", adminHandler.TestAlipayConfig)
			admin.POST("/payment/wechat/test", adminHandler.TestWechatPayConfig)

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

			// OAuth2 应用管理
			admin.GET("/oauth/clients", adminOAuthHandler.GetClients)
			admin.POST("/oauth/clients", adminOAuthHandler.CreateClient)
			admin.GET("/oauth/clients/:id", adminOAuthHandler.GetClient)
			admin.PUT("/oauth/clients/:id", adminOAuthHandler.UpdateClient)
			admin.DELETE("/oauth/clients/:id", adminOAuthHandler.DeleteClient)

			// 系统监控
			admin.GET("/system/status", systemHandler.GetSystemStatus)
			admin.GET("/system/health", systemHandler.HealthCheck)

			// 数据备份
			admin.POST("/backups", backupHandler.CreateBackup)
			admin.GET("/backups", backupHandler.ListBackups)
			admin.GET("/backups/:filename/download", backupHandler.DownloadBackup)
			admin.DELETE("/backups/:filename", backupHandler.DeleteBackup)

			// 用户登录历史管理
			admin.GET("/users/:id/login-history", adminHandler.GetUserLoginHistory)
			// 用户设备管理
			admin.GET("/users/:id/devices", adminHandler.GetUserDevices)
			admin.DELETE("/users/:id/devices/:deviceId", adminHandler.DeactivateUserDevice)
			// 用户隐私设置管理
			admin.GET("/users/:id/privacy-settings", adminHandler.GetUserPrivacySettings)
			admin.PUT("/users/:id/privacy-settings", adminHandler.UpdateUserPrivacySettings)

			// IP地理位置管理
			adminIPGeoHandler := handlers.NewAdminIPGeoHandler()
			admin.GET("/ip-geo/config", adminIPGeoHandler.GetConfig)
			admin.PUT("/ip-geo/config", adminIPGeoHandler.UpdateConfig)
			admin.POST("/ip-geo/update", adminIPGeoHandler.UpdateDatabase)
			admin.GET("/ip-geo/test", adminIPGeoHandler.TestIP)
			admin.GET("/ip-geo/logs", adminIPGeoHandler.GetUpdateLogs)
			admin.GET("/ip-geo/status", adminIPGeoHandler.GetStatus)
			admin.GET("/ip-geo/info", adminIPGeoHandler.GetDatabaseInfo)
		}

		// 公开路由 - IP地理位置查询
		// 可选参数: ?ip=1.2.3.4 查询指定IP，不传则查询访问者自身IP
		api.GET("/ip-geo/lookup", func(c *gin.Context) {
			ip := c.Query("ip")
			if ip == "" {
				ip = c.ClientIP()
			}
			location := ipGeoService.GetFullLocation(ip)
			c.JSON(200, gin.H{
				"ip":       ip,
				"location": location,
			})
		})

		// 公开路由 - 健康检查
		api.GET("/system/health", systemHandler.HealthCheck)

		// 公开路由 - 站点设置（不需要认证）
		api.GET("/settings", adminHandler.GetPublicSettings)

		// 公开页面路由
		api.GET("/pages/:slug", pageHandler.GetBySlug)
		api.GET("/pages", pageHandler.ListPublished)

		// 探索/发现页路由
		api.GET("/explore/feed", exploreHandler.GetExploreFeed)
		api.GET("/explore/hot", exploreHandler.GetHotPhotosets)

		// 创作者数据统计路由
		api.GET("/creator/stats", middleware.Auth(), creatorStatsHandler.GetCreatorStats)
		api.GET("/creator/stats/daily", middleware.Auth(), creatorStatsHandler.GetDailyStats)
		api.GET("/creator/stats/photosets", middleware.Auth(), creatorStatsHandler.GetPhotoSetStats)

		// OAuth2 公开端点（第三方应用调用）
		oauth := api.Group("/oauth")
		{
			oauth.GET("/authorize", middleware.OptionalAuth(), oauthHandler.Authorize)
			oauth.POST("/authorize", middleware.OptionalAuth(), oauthHandler.AuthorizeConfirm)
			oauth.POST("/token", oauthHandler.Token)
			oauth.POST("/revoke", oauthHandler.Revoke)
			oauth.GET("/userinfo", oauthHandler.UserInfo)
		}
	}

	// === API v1 路由（与 /api 共享相同处理器，便于版本控制） ===
	v1 := r.Group("/api/v1")
	{
		// 认证路由
		v1Auth := v1.Group("/auth")
		{
			v1Auth.GET("/captcha", middleware.CaptchaRateLimit(), captchaHandler.Generate)
			v1Auth.POST("/register", middleware.RegisterRateLimit(), authHandler.Register)
			v1Auth.POST("/login", middleware.LoginRateLimit(), authHandler.Login)
			v1Auth.GET("/me", middleware.OptionalAuth(), authHandler.Me)
			v1Auth.PUT("/password", middleware.Auth(), authHandler.ChangePassword)
			v1Auth.PUT("/profile", middleware.Auth(), authHandler.UpdateProfile)
			v1Auth.POST("/logout", middleware.Auth(), authHandler.Logout)
			v1Auth.POST("/forgot-password", middleware.ForgotPasswordRateLimit(), authHandler.ForgotPassword)
			v1Auth.POST("/reset-password", authHandler.ResetPasswordByToken)
			v1Auth.GET("/email-config", authHandler.CheckEmailConfig)
			// 邮箱验证码和绑定邮箱
			v1Auth.POST("/send-email-code", authHandler.SendEmailCode)
			v1Auth.POST("/verify-email-code", authHandler.VerifyEmailCode)
			v1Auth.PUT("/bind-email", middleware.Auth(), authHandler.BindEmail)
			// 账号注销
			v1Auth.POST("/request-deletion", middleware.Auth(), accountDeletionHandler.RequestDeletion)
			v1Auth.POST("/cancel-deletion", middleware.Auth(), accountDeletionHandler.CancelDeletion)
			v1Auth.GET("/deletion-status", middleware.Auth(), accountDeletionHandler.GetDeletionStatus)
		}

		// 套图路由
		v1Photosets := v1.Group("/photosets")
		{
			v1Photosets.GET("", middleware.OptionalAuth(), photosetHandler.List)
			v1Photosets.GET("/advanced", middleware.OptionalAuth(), photosetHandler.AdvancedList)
			v1Photosets.GET("/trash", middleware.Auth(), photosetHandler.GetTrash)
			v1Photosets.GET("/:id", middleware.OptionalAuth(), photosetHandler.Detail)
			v1Photosets.POST("", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Create)
			v1Photosets.PUT("/:id", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Update)
			v1Photosets.DELETE("/:id", middleware.Auth(), middleware.RequireRoles("creator", "admin"), photosetHandler.Delete)
			v1Photosets.POST("/:id/restore", middleware.Auth(), photosetHandler.Restore)
			v1Photosets.DELETE("/:id/permanent", middleware.Auth(), photosetHandler.PermanentDelete)
			v1Photosets.GET("/:id/download", middleware.Auth(), photosetHandler.Download)
			// 分享链接
			v1Photosets.GET("/:id/share-link", shareHandler.GenerateShareLink)
			// 套图评价
			v1Photosets.GET("/:id/reviews", middleware.OptionalAuth(), reviewHandler.List)
			v1Photosets.GET("/:id/reviews/summary", middleware.OptionalAuth(), reviewHandler.GetSummary)
			v1Photosets.GET("/:id/reviews/mine", middleware.Auth(), reviewHandler.GetMyReview)
			v1Photosets.POST("/:id/reviews", middleware.Auth(), reviewHandler.Create)
			v1Photosets.PUT("/:id/reviews/:reviewId", middleware.Auth(), reviewHandler.Update)
			v1Photosets.DELETE("/:id/reviews/:reviewId", middleware.Auth(), reviewHandler.Delete)
		}

		// 分享链接验证（公开）
		v1.GET("/share/:id", shareHandler.VerifyShareLink)

		// 标签路由
		v1.GET("/tags", tagHandler.List)

		// 分类公开路由
		v1.GET("/categories", categoryHandler.List)

		// 收藏路由
		v1Fav := v1.Group("/favorites")
		{
			v1Fav.Use(middleware.Auth())
			v1Fav.POST("/:photosetId", favHandler.Add)
			v1Fav.DELETE("/:photosetId", favHandler.Remove)
			v1Fav.GET("", favHandler.List)
			v1Fav.POST("/batch-remove", favHandler.BatchRemove)
		}

		// 评论路由
		v1Comments := v1.Group("/photosets/:id/comments")
		{
			v1Comments.GET("", middleware.OptionalAuth(), commentHandler.List)
			v1Comments.POST("", middleware.Auth(), commentHandler.Create)
			v1Comments.GET("/:commentId/replies", middleware.OptionalAuth(), commentHandler.GetReplies)
			v1Comments.DELETE("/:commentId", middleware.Auth(), commentHandler.Delete)
			v1Comments.POST("/:commentId/like", middleware.Auth(), commentHandler.ToggleLike)
		}

		// 合集路由
		v1Collections := v1.Group("/collections")
		{
			v1Collections.Use(middleware.Auth())
			v1Collections.POST("", collectionHandler.Create)
			v1Collections.GET("", collectionHandler.List)
			v1Collections.GET("/:id", collectionHandler.Get)
			v1Collections.PUT("/:id", collectionHandler.Update)
			v1Collections.DELETE("/:id", collectionHandler.Delete)
			v1Collections.POST("/:id/items", collectionHandler.AddItem)
			v1Collections.DELETE("/:id/items/:photosetId", collectionHandler.RemoveItem)
			v1Collections.POST("/:id/items/batch", collectionHandler.BatchAddItems)
			v1Collections.GET("/by-photoset/:photosetId", collectionHandler.GetCollectionsForPhotoset)
		}

		// 上传路由
		v1Upload := v1.Group("/upload")
		{
			v1Upload.Use(middleware.Auth(), middleware.RequireRoles("creator", "admin"))
			v1Upload.POST("/image", uploadHandler.UploadImage)
		}

		// 社区上传路由
		v1CommunityUpload := v1.Group("/community/upload")
		{
			v1CommunityUpload.Use(middleware.Auth())
			v1CommunityUpload.POST("/image", uploadHandler.UploadImage)
			v1CommunityUpload.POST("/video", uploadHandler.UploadVideo)
		}

		// 头像上传路由（所有登录用户可用）
		v1.POST("/user/avatar", middleware.Auth(), uploadHandler.UploadImage)

		// 用户路由
		v1.GET("/users/profile", middleware.Auth(), authHandler.Me)

		// 用户隐私设置路由
		v1.GET("/user/privacy-settings", middleware.Auth(), privacyHandler.GetPrivacySettings)
		v1.PUT("/user/privacy-settings", middleware.Auth(), privacyHandler.UpdatePrivacySettings)

		// 登录历史路由
		v1.GET("/user/login-history", middleware.Auth(), loginHistoryHandler.GetLoginHistory)

		// 用户设备管理路由
		v1.GET("/user/devices", middleware.Auth(), userDeviceHandler.GetUserDevices)
		v1.DELETE("/user/devices/:deviceId", middleware.Auth(), userDeviceHandler.DeactivateDevice)
		v1.DELETE("/user/devices", middleware.Auth(), userDeviceHandler.DeactivateOtherDevices)

		// 浏览历史路由
		v1.POST("/user/view-history", middleware.Auth(), viewHistoryHandler.Record)
		v1.GET("/user/view-history", middleware.Auth(), viewHistoryHandler.List)
		v1.DELETE("/user/view-history/:id", middleware.Auth(), viewHistoryHandler.Delete)
		v1.DELETE("/user/view-history", middleware.Auth(), viewHistoryHandler.ClearAll)
		v1.POST("/user/view-history/batch-delete", middleware.Auth(), viewHistoryHandler.BatchDelete)

		// 推送通知路由
		v1.POST("/user/push-tokens", middleware.Auth(), pushNotificationHandler.RegisterToken)
		v1.DELETE("/user/push-tokens", middleware.Auth(), pushNotificationHandler.UnregisterToken)

		// 用户反馈路由
		v1.POST("/feedback", middleware.OptionalAuth(), feedbackHandler.Create)
		v1.POST("/feedback/image", middleware.OptionalAuth(), feedbackHandler.UploadFeedbackImage)

		// 会员套餐路由
		v1.GET("/memberships", membershipHandler.List)

		// 订单路由
		v1Orders := v1.Group("/orders")
		{
			v1Orders.Use(middleware.Auth())
			v1Orders.GET("", orderHandler.List)
			v1Orders.POST("", orderHandler.Create)
			v1Orders.POST("/:id/pay", orderHandler.Pay)
			v1Orders.POST("/:id/refund", orderHandler.Refund)
		}

		// 管理后台路由
		v1Admin := v1.Group("/admin")
		{
			v1Admin.Use(middleware.Auth(), middleware.RequireRoles("admin"))
			v1Admin.GET("/users", adminHandler.GetUsers)
			v1Admin.POST("/users", adminHandler.CreateUser)
			v1Admin.GET("/users/export", adminHandler.ExportUsers)
			v1Admin.GET("/users/:id", adminHandler.GetUserDetail)
			v1Admin.GET("/photosets", adminHandler.GetPhotoSetsByStatus)
			v1Admin.POST("/photosets/:id/approve", adminHandler.ApprovePhotoSet)
			v1Admin.POST("/photosets/:id/reject", adminHandler.RejectPhotoSet)
			v1Admin.POST("/photosets/batch/approve", adminHandler.BatchApprovePhotoSets)
			v1Admin.POST("/photosets/batch/reject", adminHandler.BatchRejectPhotoSets)
			v1Admin.POST("/photosets/batch/delete", adminHandler.BatchDeletePhotoSets)
			v1Admin.PUT("/users/:id/ban", adminHandler.BanUser)
			v1Admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)
			v1Admin.PUT("/users/:id/password", adminHandler.ResetUserPassword)
			v1Admin.GET("/stats", adminHandler.Stats)
			v1Admin.GET("/stats/trend", adminHandler.StatsTrend)
			v1Admin.GET("/logs", adminHandler.GetAdminLogs)
			v1Admin.GET("/orders", adminHandler.GetOrders)
			v1Admin.GET("/orders/export", adminHandler.ExportOrders)
			v1Admin.POST("/orders/:id/refund", adminHandler.AdminRefund)
			v1Admin.GET("/photosets/export", adminHandler.ExportPhotoSets)
			v1Admin.GET("/tags", tagHandler.AdminList)
			v1Admin.POST("/tags", tagHandler.Create)
			v1Admin.PUT("/tags/:id", tagHandler.Update)
			v1Admin.DELETE("/tags/:id", tagHandler.Delete)
			v1Admin.GET("/categories", categoryHandler.AdminList)
			v1Admin.POST("/categories", categoryHandler.Create)
			v1Admin.PUT("/categories/:id", categoryHandler.Update)
			v1Admin.DELETE("/categories/:id", categoryHandler.Delete)
			v1Admin.GET("/settings", adminHandler.GetSettings)
			v1Admin.PUT("/settings", adminHandler.UpdateSettings)
			v1Admin.POST("/system/restart", middleware.RestartRateLimit(), adminHandler.RestartServer)
			v1Admin.POST("/mail/test-connection", adminHandler.TestMailConnection)
			v1Admin.GET("/mail/config", adminHandler.GetMailConfig)
			v1Admin.POST("/mail/send-test", adminHandler.SendTestMail)
			v1Admin.GET("/watermark/info", adminHandler.GetWatermarkInfo)
			v1Admin.POST("/storage/test", adminHandler.TestStorageConnection)
			v1Admin.GET("/storage/status", adminHandler.GetStorageStatus)

			// 支付配置
			v1Admin.POST("/payment/alipay/test", adminHandler.TestAlipayConfig)
			v1Admin.POST("/payment/wechat/test", adminHandler.TestWechatPayConfig)
			v1Admin.GET("/pages", pageHandler.AdminList)
			v1Admin.POST("/pages", pageHandler.AdminCreate)
			v1Admin.GET("/pages/:id", pageHandler.AdminGet)
			v1Admin.PUT("/pages/:id", pageHandler.AdminUpdate)
			v1Admin.DELETE("/pages/:id", pageHandler.AdminDelete)
			v1Admin.GET("/memberships", membershipHandler.AdminList)
			v1Admin.POST("/memberships", membershipHandler.AdminCreate)
			v1Admin.PUT("/memberships/:id", membershipHandler.AdminUpdate)
			v1Admin.DELETE("/memberships/:id", membershipHandler.AdminDelete)
			v1Admin.GET("/dev/api-keys", adminHandler.ListApiKeys)
			v1Admin.POST("/dev/api-keys", adminHandler.CreateApiKey)
			v1Admin.DELETE("/dev/api-keys/:id", adminHandler.DeleteApiKey)
			v1Admin.GET("/dev/api-docs", adminHandler.GetApiDocs)
			v1Admin.GET("/dev/sign-url-docs", adminHandler.GetSignUrlDocs)
			v1Admin.GET("/oauth/clients", adminOAuthHandler.GetClients)
			v1Admin.POST("/oauth/clients", adminOAuthHandler.CreateClient)
			v1Admin.GET("/oauth/clients/:id", adminOAuthHandler.GetClient)
			v1Admin.PUT("/oauth/clients/:id", adminOAuthHandler.UpdateClient)
			v1Admin.DELETE("/oauth/clients/:id", adminOAuthHandler.DeleteClient)
			v1Admin.GET("/system/status", systemHandler.GetSystemStatus)
			v1Admin.GET("/system/health", systemHandler.HealthCheck)
			v1Admin.POST("/backups", backupHandler.CreateBackup)
			v1Admin.GET("/backups", backupHandler.ListBackups)
			v1Admin.GET("/backups/:filename/download", backupHandler.DownloadBackup)
			v1Admin.DELETE("/backups/:filename", backupHandler.DeleteBackup)
		}

		// 公开路由
		v1.GET("/settings", adminHandler.GetPublicSettings)
		v1.GET("/pages/:slug", pageHandler.GetBySlug)
		v1.GET("/pages", pageHandler.ListPublished)

		// 探索/发现页路由
		v1.GET("/explore/feed", exploreHandler.GetExploreFeed)
		v1.GET("/explore/hot", exploreHandler.GetHotPhotosets)

		// 创作者数据统计路由
		v1.GET("/creator/stats", middleware.Auth(), creatorStatsHandler.GetCreatorStats)
		v1.GET("/creator/stats/daily", middleware.Auth(), creatorStatsHandler.GetDailyStats)
		v1.GET("/creator/stats/photosets", middleware.Auth(), creatorStatsHandler.GetPhotoSetStats)

		// OAuth2 公开端点
		v1OAuth := v1.Group("/oauth")
		{
			v1OAuth.GET("/authorize", middleware.OptionalAuth(), oauthHandler.Authorize)
			v1OAuth.POST("/authorize", middleware.OptionalAuth(), oauthHandler.AuthorizeConfirm)
			v1OAuth.POST("/token", oauthHandler.Token)
			v1OAuth.POST("/revoke", oauthHandler.Revoke)
			v1OAuth.GET("/userinfo", oauthHandler.UserInfo)
		}
	}

	// === 注册社区路由（/api 和 /api/v1） ===
	RegisterCommunityRoutes(api, communityHandler, followHandler, adminCommunityHandler, adminLevelHandler, notificationHandler, messageHandler, userLevelHandler, wsHandler, blockHandler)
	RegisterCommunityRoutes(v1, communityHandler, followHandler, adminCommunityHandler, adminLevelHandler, notificationHandler, messageHandler, userLevelHandler, wsHandler, blockHandler)
}
