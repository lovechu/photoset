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
	notificationHandler *handlers.NotificationHandler,
	messageHandler *handlers.MessageHandler,
	userLevelHandler *handlers.UserLevelHandler,
) {
	// Public routes (with optional auth)
	public := r.Group("/api/community")
	{
		// Posts
		public.GET("/posts", communityHandler.GetPosts)
		public.GET("/posts/:id", communityHandler.GetPostDetail)
		public.GET("/posts/:id/replies", communityHandler.GetReplies)
		public.GET("/posts/search", communityHandler.SearchPosts)
		public.GET("/categories", communityHandler.GetCategories)
		public.GET("/hot", communityHandler.GetHotPosts)

		// User profile
		public.GET("/users/:id", communityHandler.GetUserProfile)
		public.GET("/users/:id/posts", communityHandler.GetUserPosts)

		// User search for @mention
		public.GET("/users/search", communityHandler.SearchUsersForMention)

		// Share stats (public)
		public.GET("/posts/:id/shares", communityHandler.GetShareStats)

		// Tags
		public.GET("/tags/popular", communityHandler.GetPopularTags)
		public.GET("/tags/search", communityHandler.SearchTags)
		public.GET("/tags/:name/posts", communityHandler.GetPostsByTagName)

		// Topics
		public.GET("/topics/hot", communityHandler.GetHotTopics)
		public.GET("/topics/search", communityHandler.SearchTopics)
		public.GET("/topics/:name/posts", communityHandler.GetPostsByTopicName)

		// Levels (public)
		public.GET("/levels", userLevelHandler.GetAllLevelConfigs)
		public.GET("/users/:id/level", userLevelHandler.GetUserLevelInfoByID)
		public.GET("/users/:id/achievements", userLevelHandler.GetUserAchievementsByID)
		public.GET("/points-mall/items", userLevelHandler.GetPointsMallItems)
		public.GET("/points-mall/categories", userLevelHandler.GetPointsMallCategories)
	}

	// Protected routes (require login)
	protected := r.Group("/api/community")
	protected.Use(middleware.Auth())
	{
		// Posts
		protected.POST("/posts", communityHandler.CreatePost)
		protected.PUT("/posts/:id", communityHandler.UpdatePost)
		protected.DELETE("/posts/:id", communityHandler.DeletePost)

		// Replies
		protected.POST("/posts/:id/reply", communityHandler.CreateReply)
		protected.PUT("/replies/:id", communityHandler.UpdateReply)
		protected.DELETE("/replies/:id", communityHandler.DeleteReply)

		// Likes
		protected.POST("/posts/:id/like", communityHandler.TogglePostLike)
		protected.POST("/replies/:id/like", communityHandler.ToggleReplyLike)

		// Shares
		protected.POST("/posts/:id/share", communityHandler.RecordShare)

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

		// Drafts
		protected.POST("/drafts", communityHandler.SaveDraft)
		protected.GET("/me/drafts", communityHandler.GetMyDrafts)
		protected.DELETE("/drafts/:id", communityHandler.DeleteDraft)

		// Recommendations
		protected.GET("/recommendations", communityHandler.GetRecommendations)
		protected.GET("/interests", communityHandler.GetUserInterestProfile)

		// Follow
		protected.POST("/users/:id/follow", followHandler.Follow)
		protected.DELETE("/users/:id/follow", followHandler.Unfollow)
		protected.GET("/users/:id/follow/check", followHandler.CheckFollowing)
		protected.POST("/users/batch-follow-check", followHandler.BatchCheckFollowing)
		protected.GET("/users/:id/following", followHandler.GetFollowingList)
		protected.GET("/users/:id/followers", followHandler.GetFollowerList)

		// Notifications
		protected.GET("/notifications", notificationHandler.GetNotifications)
		protected.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
		protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
		protected.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
		protected.DELETE("/notifications/:id", notificationHandler.DeleteNotification)
		protected.DELETE("/notifications", notificationHandler.DeleteAllNotifications)

		// Messages
		protected.GET("/messages/conversations", messageHandler.GetConversations)
		protected.GET("/messages/conversations/:user_id", messageHandler.GetConversation)
		protected.POST("/messages", messageHandler.SendMessage)
		protected.GET("/messages/unread-count", messageHandler.GetUnreadCount)
		protected.PUT("/messages/:id/read", messageHandler.MarkAsRead)
		protected.PUT("/messages/conversations/:user_id/read", messageHandler.MarkConversationAsRead)
		protected.DELETE("/messages/:id", messageHandler.DeleteMessage)

		// User Level & Achievements
		protected.GET("/user/level", userLevelHandler.GetUserLevelInfo)
		protected.GET("/user/achievements", userLevelHandler.GetUserAchievements)

		// Points Mall
		protected.POST("/points-mall/exchange", userLevelHandler.ExchangeItem)
		protected.GET("/points-mall/history", userLevelHandler.GetUserExchangeHistory)

		// Points Leaderboard
		protected.GET("/points/leaderboard", userLevelHandler.GetPointsLeaderboard)
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
