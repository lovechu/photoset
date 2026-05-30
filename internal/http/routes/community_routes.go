package routes

import (
	"photoset/internal/http/handlers"
	"photoset/internal/http/handlers/admin"
	"photoset/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterCommunityRoutes registers community module routes
func RegisterCommunityRoutes(
	r *gin.Engine,
	communityHandler *handlers.CommunityHandler,
	followHandler *handlers.FollowHandler,
	adminHandler *admin.AdminCommunityHandler,
) {
	// Public routes (with optional auth)
	public := r.Group("/api/community")
	{
		// Posts
		public.GET("/posts", communityHandler.GetPosts)
		public.GET("/posts/:id", communityHandler.GetPostDetail)
		public.GET("/posts/:id/replies", communityHandler.GetReplies)
		public.GET("/categories", communityHandler.GetCategories)
		public.GET("/hot", communityHandler.GetHotPosts)
	}

	// Protected routes (require login)
	protected := r.Group("/api/community")
	protected.Use(middleware.Auth())
	{
		// Posts
		protected.POST("/posts", communityHandler.CreatePost)

		// Replies
		protected.POST("/posts/:id/reply", communityHandler.CreateReply)

		// Likes
		protected.POST("/posts/:id/like", communityHandler.TogglePostLike)
		protected.POST("/replies/:id/like", communityHandler.ToggleReplyLike)

		// Favorites
		protected.POST("/posts/:id/favorite", communityHandler.TogglePostFavorite)
		protected.GET("/posts/:id/favorite/check", communityHandler.CheckPostFavorite)
		protected.GET("/me/favorites", communityHandler.GetMyFavorites)

		// Reports
		protected.POST("/posts/:id/report", communityHandler.ReportPost)
		protected.POST("/replies/:id/report", communityHandler.ReportReply)

		// My content
		protected.GET("/me/posts", communityHandler.GetMyPosts)
		protected.GET("/me/replies", communityHandler.GetMyReplies)
		protected.GET("/me/points", communityHandler.GetMyPoints)

		// Follow
		protected.POST("/users/:id/follow", followHandler.Follow)
		protected.DELETE("/users/:id/follow", followHandler.Unfollow)
		protected.GET("/users/:id/follow/check", followHandler.CheckFollowing)
		protected.POST("/users/batch-follow-check", followHandler.BatchCheckFollowing)
		protected.GET("/users/:id/following", followHandler.GetFollowingList)
		protected.GET("/users/:id/followers", followHandler.GetFollowerList)
	}

	// Admin routes
	adminGroup := r.Group("/api/admin/community")
	adminGroup.Use(middleware.Auth(), middleware.AdminOnly())
	{
		// Post management
		adminGroup.GET("/posts", adminHandler.GetPosts)
		adminGroup.PUT("/posts/:id/pin", adminHandler.TogglePin)
		adminGroup.PUT("/posts/:id/status", adminHandler.UpdateStatus)
		adminGroup.DELETE("/posts/:id", adminHandler.DeletePost)

		// Reply management
		adminGroup.GET("/replies", adminHandler.GetReplies)
		adminGroup.DELETE("/replies/:id", adminHandler.DeleteReply)

		// Sensitive words management
		adminGroup.GET("/keywords", adminHandler.GetKeywords)
		adminGroup.POST("/keywords", adminHandler.AddKeyword)
		adminGroup.PUT("/keywords/:id", adminHandler.UpdateKeyword)
		adminGroup.DELETE("/keywords/:id", adminHandler.DeleteKeyword)
		adminGroup.PUT("/keywords/reload", adminHandler.ReloadKeywords)

		// Reports management
		adminGroup.GET("/reports", adminHandler.GetReports)
		adminGroup.PUT("/reports/:id/resolve", adminHandler.ResolveReport)

		// User points management
		adminGroup.GET("/users", adminHandler.GetUsers)
		adminGroup.PUT("/users/:id/points", adminHandler.AdjustPoints)

		// Statistics
		adminGroup.GET("/stats", adminHandler.GetStats)

		// Category management
		adminGroup.GET("/categories", adminHandler.ListCategories)
		adminGroup.POST("/categories", adminHandler.CreateCategory)
		adminGroup.PUT("/categories/sort", adminHandler.SortCategories)
		adminGroup.PUT("/categories/:id", adminHandler.UpdateCategory)
		adminGroup.DELETE("/categories/:id", adminHandler.DeleteCategory)
	}
}
