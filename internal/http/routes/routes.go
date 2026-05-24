package routes

import (
	"photoset/internal/config"
	"photoset/internal/database"
	"photoset/internal/http/handlers"
	"photoset/internal/http/handlers/admin"
	"photoset/internal/repository"
	"photoset/internal/service"

	"github.com/gin-gonic/gin"
)

// Setup sets up all routes for the application
func Setup(r *gin.Engine, cfg *config.Config) {
	// Initialize repositories
	postRepo := repository.NewPostRepository(database.GetMySQL())
	replyRepo := repository.NewPostReplyRepository(database.GetMySQL())
	likeRepo := repository.NewPostLikeRepository(database.GetMySQL())
	replyLikeRepo := repository.NewPostReplyLikeRepository(database.GetMySQL())
	pointRepo := repository.NewUserPointRepository(database.GetMySQL())
	wordRepo := repository.NewSensitiveWordRepository(database.GetMySQL())
	reportRepo := repository.NewPostReportRepository(database.GetMySQL())

	// Initialize services
	pointService := service.NewPointService(pointRepo)
	filterService := service.NewSensitiveFilterService(wordRepo)
	communityService := service.NewCommunityService(
		postRepo,
		replyRepo,
		likeRepo,
		replyLikeRepo,
		pointRepo,
		reportRepo,
		pointService,
		filterService,
	)
	hotPostsService := service.NewHotPostsService(postRepo)

	// Initialize handlers
	communityHandler := handlers.NewCommunityHandler(
		database.GetMySQL(),
		communityService,
		pointService,
		hotPostsService,
	)
	adminCommunityHandler := admin.NewAdminCommunityHandler(database.GetMySQL())

	// Load sensitive words on startup
	service.InitSensitiveWords(wordRepo)

	// Register routes
	RegisterCommunityRoutes(r, communityHandler, adminCommunityHandler)

	// TODO: Register other routes (auth, photoset, etc.)
	// This should be done by importing and calling other route registration functions
}
