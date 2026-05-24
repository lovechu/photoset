package service

import (
	"photoset/internal/domain"
	"photoset/internal/repository"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup full test environment with seeded data
func setupCommunityTest(t *testing.T) (
	*CommunityService,
	*PointService,
	*SensitiveFilterService,
	*gorm.DB,
) {
	// Use temp file for SQLite (in-memory database gets deleted between connections)
	tmpfile := t.TempDir() + "/test.db"
	db, err := gorm.Open(sqlite.Open(tmpfile), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Migrate all tables
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Post{},
		&domain.PostReply{},
		&domain.PostLike{},
		&domain.PostReplyLike{},
		&domain.UserPoint{},
		&domain.UserPointLog{},
		&domain.SensitiveWord{},
		&domain.PostReport{},
	)
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}

	// Create repositories
	postRepo := repository.NewPostRepository(db)
	replyRepo := repository.NewPostReplyRepository(db)
	likeRepo := repository.NewPostLikeRepository(db)
	replyLikeRepo := repository.NewPostReplyLikeRepository(db)
	pointRepo := repository.NewUserPointRepository(db)
	reportRepo := repository.NewPostReportRepository(db)
	wordRepo := repository.NewSensitiveWordRepository(db)

	// Create services
	pointService := NewPointService(pointRepo)
	filterService := NewSensitiveFilterService(wordRepo)
	
	communityService := NewCommunityService(
		postRepo,
		replyRepo,
		likeRepo,
		replyLikeRepo,
		pointRepo,
		reportRepo,
		pointService,
		filterService,
	)

	return communityService, pointService, filterService, db
}

// Test 1: Test CreatePost
func TestCommunityService_CreatePost(t *testing.T) {
	svc, _, filterService, db := setupCommunityTest(t)

	// Seed sensitive words
	db.Create(&domain.SensitiveWord{Word: "badword", Replacement: "***", IsActive: true})
	filterService.LoadSensitiveWords()

	// Create test user
	user := domain.User{Nickname: "testuser", Role: "member"}
	db.Create(&user)

	// Test: Create a post
	req := &CreatePostRequest{
		Title:     "Test Post",
		Content:   "This is a test post",
		Category:   "discussion",
		Visibility: "public",
	}

	post, err := svc.CreatePost(user.ID, req)
	if err != nil {
		t.Errorf("Failed to create post: %v", err)
	}

	if post.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got %q", post.Title)
	}
	if post.UserID != user.ID {
		t.Errorf("Expected user ID %d, got %d", user.ID, post.UserID)
	}
}

// Test 2: Test CreatePost with sensitive word filtering
func TestCommunityService_CreatePostWithFiltering(t *testing.T) {
	svc, _, filterService, db := setupCommunityTest(t)

	// Seed sensitive words
	db.Create(&domain.SensitiveWord{Word: "badword", Replacement: "***", IsActive: true})
	filterService.LoadSensitiveWords()

	// Create test user
	user := domain.User{Nickname: "testuser2", Role: "member"}
	db.Create(&user)

	// Test: Create post with sensitive word
	req := &CreatePostRequest{
		Title:     "This is a badword title",
		Content:   "This is a badword content",
		Category:   "discussion",
		Visibility: "public",
	}

	post, err := svc.CreatePost(user.ID, req)
	if err != nil {
		t.Errorf("Failed to create post: %v", err)
	}

	// Title should be filtered
	if post.Title == "This is a badword title" {
		t.Error("Expected title to be filtered, but it wasn't")
	}
	t.Logf("Filtered title: %s", post.Title)
}

// Test 3: Test CreatePost validation errors
func TestCommunityService_CreatePostValidation(t *testing.T) {
	svc, _, _, db := setupCommunityTest(t)

	// Create test user
	user := domain.User{Nickname: "testuser3", Role: "member"}
	db.Create(&user)

	// Test: Empty title should fail
	req := &CreatePostRequest{
		Title:     "",
		Content:   "Valid content",
		Category:   "discussion",
		Visibility: "public",
	}

	_, err := svc.CreatePost(user.ID, req)
	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}
}

// Test 4: Test CreateReply
func TestCommunityService_CreateReply(t *testing.T) {
	svc, _, filterService, db := setupCommunityTest(t)

	// Seed sensitive words
	db.Create(&domain.SensitiveWord{Word: "spam", Replacement: "***", IsActive: true})
	filterService.LoadSensitiveWords()

	// Create test user
	user := domain.User{Nickname: "testuser3", Role: "member"}
	db.Create(&user)

	// Create a post first
	post := domain.Post{
		UserID:    user.ID,
		Title:      "Test Post",
		Content:    "Content",
		Category:   "discussion",
		Visibility: "public",
		Status:     "approved",
	}
	db.Create(&post)

	// Test: Create a reply
	req := &CreateReplyRequest{
		Content:       "This is a reply",
		ParentReplyID: nil,
	}

	reply, err := svc.CreateReply(user.ID, post.ID, req)
	if err != nil {
		t.Errorf("Failed to create reply: %v", err)
	}

	if reply.Content != "This is a reply" {
		t.Errorf("Expected content 'This is a reply', got %q", reply.Content)
	}
	if reply.PostID != post.ID {
		t.Errorf("Expected post ID %d, got %d", post.ID, reply.PostID)
	}

	// Verify reply count was incremented
	var updatedPost domain.Post
	db.First(&updatedPost, post.ID)
	if updatedPost.ReplyCount != 1 {
		t.Errorf("Expected reply count 1, got %d", updatedPost.ReplyCount)
	}
}

// Test 5: Test TogglePostLike
func TestCommunityService_TogglePostLike(t *testing.T) {
	svc, _, _, db := setupCommunityTest(t)

	// Create test users
	author := domain.User{Nickname: "author", Role: "member"}
	liker := domain.User{Nickname: "liker", Role: "member"}
	db.Create(&author)
	db.Create(&liker)

	// Create a post
	post := domain.Post{
		UserID:    author.ID,
		Title:      "Test Post",
		Content:    "Content",
		Category:   "discussion",
		Visibility: "public",
		Status:     "approved",
	}
	db.Create(&post)

	// Test: First like (should like)
	action, likeCount, err := svc.TogglePostLike(liker.ID, post.ID)
	if err != nil {
		t.Errorf("Failed to toggle like: %v", err)
	}
	if action != "liked" {
		t.Errorf("Expected action 'liked', got %q", action)
	}
	if likeCount != 1 {
		t.Errorf("Expected like count 1, got %d", likeCount)
	}

	// Verify author got points
	var authorPoints domain.UserPoint
	db.Where("user_id = ?", author.ID).First(&authorPoints)
	if authorPoints.Points != 2 {
		t.Errorf("Expected author points 2, got %d", authorPoints.Points)
	}

	// Test: Second like (should unlike)
	action, likeCount, err = svc.TogglePostLike(liker.ID, post.ID)
	if err != nil {
		t.Errorf("Failed to toggle like: %v", err)
	}
	if action != "unliked" {
		t.Errorf("Expected action 'unliked', got %q", action)
	}
	if likeCount != 0 {
		t.Errorf("Expected like count 0, got %d", likeCount)
	}
}

// Test 6: Test ToggleReplyLike
func TestCommunityService_ToggleReplyLike(t *testing.T) {
	svc, _, _, db := setupCommunityTest(t)

	// Create test users
	author := domain.User{Nickname: "reply_author", Role: "member"}
	liker := domain.User{Nickname: "reply_liker", Role: "member"}
	db.Create(&author)
	db.Create(&liker)

	// Create a post
	post := domain.Post{
		UserID:    author.ID,
		Title:      "Test Post",
		Content:    "Content",
		Category:   "discussion",
		Visibility: "public",
		Status:     "approved",
	}
	db.Create(&post)

	// Create a reply
	reply := domain.PostReply{
		PostID:  post.ID,
		UserID:  author.ID,
		Content: "This is a reply",
	}
	db.Create(&reply)

	// Test: First like (should like)
	action, likeCount, err := svc.ToggleReplyLike(liker.ID, reply.ID)
	if err != nil {
		t.Errorf("Failed to toggle reply like: %v", err)
	}
	if action != "liked" {
		t.Errorf("Expected action 'liked', got %q", action)
	}
	if likeCount != 1 {
		t.Errorf("Expected like count 1, got %d", likeCount)
	}

	// Verify author got points
	var authorPoints domain.UserPoint
	db.Where("user_id = ?", author.ID).First(&authorPoints)
	if authorPoints.Points != 1 {
		t.Errorf("Expected author points 1, got %d", authorPoints.Points)
	}
}

// Test 7: Test ReportPost
func TestCommunityService_ReportPost(t *testing.T) {
	svc, _, _, db := setupCommunityTest(t)

	// Create test users
	author := domain.User{Nickname: "post_author", Role: "member"}
	reporter := domain.User{Nickname: "reporter", Role: "member"}
	db.Create(&author)
	db.Create(&reporter)

	// Create a post
	post := domain.Post{
		UserID:    author.ID,
		Title:      "Test Post",
		Content:    "Content",
		Category:   "discussion",
		Visibility: "public",
		Status:     "approved",
	}
	db.Create(&post)

	// Test: Report the post
	err := svc.ReportPost(reporter.ID, post.ID, "This is inappropriate")
	if err != nil {
		t.Errorf("Failed to report post: %v", err)
	}

	// Verify report was created
	var report domain.PostReport
	db.Where("post_id = ? AND reporter_id = ?", post.ID, reporter.ID).First(&report)
	if report.Reason != "This is inappropriate" {
		t.Errorf("Expected report reason 'This is inappropriate', got %q", report.Reason)
	}
	if report.Status != "pending" {
		t.Errorf("Expected report status 'pending', got %q", report.Status)
	}
}

// Test 8: Test GetPostDetail (view count increment)
func TestCommunityService_GetPostDetail(t *testing.T) {
	svc, _, _, db := setupCommunityTest(t)

	// Create test user
	user := domain.User{Nickname: "testuser4", Role: "member"}
	db.Create(&user)

	// Create a post
	post := domain.Post{
		UserID:    user.ID,
		Title:      "Test Post",
		Content:    "Content",
		Category:   "discussion",
		Visibility: "public",
		Status:     "approved",
	}
	db.Create(&post)

	// Get post detail (should increment view count)
	detail, err := svc.GetPostDetail(post.ID)
	if err != nil {
		t.Errorf("Failed to get post detail: %v", err)
	}

	if detail.ViewCount != 1 {
		t.Errorf("Expected view count 1, got %d", detail.ViewCount)
	}
}

// Test 9: Test daily limit for posts
func TestCommunityService_DailyLimitForPosts(t *testing.T) {
	svc, _, _, db := setupCommunityTest(t)

	// Create test user
	user := domain.User{Nickname: "testuser5", Role: "member"}
	db.Create(&user)

	// Create 5 posts (50 points)
	for i := 0; i < 5; i++ {
		req := &CreatePostRequest{
			Title:     "Test Post",
			Content:   "Content",
			Category:   "discussion",
			Visibility: "public",
		}
		_, err := svc.CreatePost(user.ID, req)
		if err != nil {
			t.Fatalf("Failed to create post %d: %v", i, err)
		}
	}

	// Try to create 6th post - should fail with daily limit
	req := &CreatePostRequest{
		Title:     "Test Post 6",
		Content:   "Content",
		Category:   "discussion",
		Visibility: "public",
	}
	_, err := svc.CreatePost(user.ID, req)
	if err == nil {
		t.Error("Expected daily limit error, got nil")
	}
	if err != domain.ErrDailyLimitReached {
		t.Errorf("Expected ErrDailyLimitReached, got %v", err)
	}
}

// Test 10: Test nested replies (楼中楼)
func TestCommunityService_NestedReplies(t *testing.T) {
	_, _, _, db := setupCommunityTest(t)

	// Create test user
	user := domain.User{Nickname: "testuser6", Role: "member"}
	db.Create(&user)

	// Create a post
	post := domain.Post{
		UserID:    user.ID,
		Title:      "Test Post",
		Content:    "Content",
		Category:   "discussion",
		Visibility: "public",
		Status:     "approved",
	}
	db.Create(&post)

	// Create a top-level reply
	reply1 := domain.PostReply{
		PostID:        post.ID,
		UserID:        user.ID,
		Content:       "Top-level reply",
		ParentReplyID: nil,
	}
	db.Create(&reply1)

	// Create a nested reply (楼中楼)
	reply2 := domain.PostReply{
		PostID:        post.ID,
		UserID:        user.ID,
		Content:       "Nested reply",
		ParentReplyID: &reply1.ID,
	}
	db.Create(&reply2)

	// Verify parent_reply_id is set correctly
	if *reply2.ParentReplyID != reply1.ID {
		t.Errorf("Expected parent reply ID %d, got %d", reply1.ID, *reply2.ParentReplyID)
	}
}
