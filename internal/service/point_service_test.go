package service

import (
	"photoset/internal/domain"
	"photoset/internal/repository"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup test database with all required tables
func setupFullTestDB(t *testing.T) *gorm.DB {
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

	return db
}

// Test 1: Test point calculation for posting
func TestPointService_AddPointsForPost(t *testing.T) {
	db := setupFullTestDB(t)

	// Create test user
	user := domain.User{Nickname: "testuser", Role: "member"}
	db.Create(&user)

	pointRepo := repository.NewUserPointRepository(db)
	pointService := NewPointService(pointRepo)

	// Test: Add points for first post
	err := pointService.AddPointsForPost(user.ID)
	if err != nil {
		t.Errorf("Failed to add points for post: %v", err)
	}

	// Verify points
	userPoint, _ := pointRepo.GetByUserID(user.ID)
	if userPoint.Points != 10 {
		t.Errorf("Expected 10 points, got %d", userPoint.Points)
	}
	if userPoint.Level != 1 {
		t.Errorf("Expected level 1, got %d", userPoint.Level)
	}
}

// Test 2: Test daily limit for posts (max 50 points = 5 posts)
func TestPointService_DailyLimitForPosts(t *testing.T) {
	db := setupFullTestDB(t)

	// Create test user
	user := domain.User{Nickname: "testuser2", Role: "member"}
	db.Create(&user)

	pointRepo := repository.NewUserPointRepository(db)
	pointService := NewPointService(pointRepo)

	// Add 50 points (5 posts worth)
	for i := 0; i < 5; i++ {
		err := pointService.AddPointsForPost(user.ID)
		if err != nil {
			t.Fatalf("Failed to add points on iteration %d: %v", i, err)
		}
	}

	// Try to add 6th post - should fail with daily limit
	err := pointService.AddPointsForPost(user.ID)
	if err == nil {
		t.Error("Expected daily limit error, got nil")
	}
	if err != domain.ErrDailyLimitReached {
		t.Errorf("Expected ErrDailyLimitReached, got %v", err)
	}
}

// Test 3: Test point calculation for reply
func TestPointService_AddPointsForReply(t *testing.T) {
	db := setupFullTestDB(t)

	// Create test user
	user := domain.User{Nickname: "testuser3", Role: "member"}
	db.Create(&user)

	pointRepo := repository.NewUserPointRepository(db)
	pointService := NewPointService(pointRepo)

	// Test: Add points for reply
	err := pointService.AddPointsForReply(user.ID)
	if err != nil {
		t.Errorf("Failed to add points for reply: %v", err)
	}

	// Verify points
	userPoint, _ := pointRepo.GetByUserID(user.ID)
	if userPoint.Points != 5 {
		t.Errorf("Expected 5 points, got %d", userPoint.Points)
	}
}

// Test 4: Test level calculation
func TestPointService_CalculateLevel(t *testing.T) {
	tests := []struct {
		points   int
		expected int
	}{
		{0, 1},
		{50, 1},
		{100, 2},
		{500, 3},
		{2000, 4},
		{5000, 5},
		{10000, 6},
		{20000, 7},
		{30000, 8},
		{40000, 9},
		{50000, 10},
		{100000, 10}, // Max level
	}

	for _, tt := range tests {
		level := CalculateLevel(tt.points)
		if level != tt.expected {
			t.Errorf("CalculateLevel(%d) = %d, want %d", tt.points, level, tt.expected)
		}
	}
}

// Test 5: Test points can go negative
func TestPointService_NegativePoints(t *testing.T) {
	db := setupFullTestDB(t)

	// Create test user with some points
	user := domain.User{Nickname: "testuser4", Role: "member"}
	db.Create(&user)

	pointRepo := repository.NewUserPointRepository(db)
	pointService := NewPointService(pointRepo)

	// Add some points first
	pointService.AddPointsForPost(user.ID)
	pointService.AddPointsForPost(user.ID) // 20 points

	// Deduct points (simulate admin delete)
	err := pointService.DeductPointsForDelete(user.ID)
	if err != nil {
		t.Errorf("Failed to deduct points: %v", err)
	}

	// Verify points (-20)
	userPoint, _ := pointRepo.GetByUserID(user.ID)
	if userPoint.Points != -20 {
		t.Errorf("Expected -20 points, got %d", userPoint.Points)
	}

	// Level should still be 1 (minimum)
	if userPoint.Level != 1 {
		t.Errorf("Expected level 1 (minimum), got %d", userPoint.Level)
	}
}

// Test 6: Test GetLevelInfo
func TestPointService_GetLevelInfo(t *testing.T) {
	db := setupFullTestDB(t)

	// Create test user
	user := domain.User{Nickname: "testuser5", Role: "member"}
	db.Create(&user)

	// Add points to reach level 2
	pointRepo := repository.NewUserPointRepository(db)
	pointService := NewPointService(pointRepo)

	// Add 10 posts (100 points = level 2)
	for i := 0; i < 10; i++ {
		pointService.AddPointsForPost(user.ID)
	}

	// Get level info
	level, levelName, currentPoints, nextLevelPoints, err := pointService.GetLevelInfo(user.ID)
	if err != nil {
		t.Errorf("Failed to get level info: %v", err)
	}

	if level != 2 {
		t.Errorf("Expected level 2, got %d", level)
	}
	if levelName != "Active Member" {
		t.Errorf("Expected 'Active Member', got %s", levelName)
	}
	if currentPoints != 100 {
		t.Errorf("Expected 100 points, got %d", currentPoints)
	}
	if nextLevelPoints != 500 {
		t.Errorf("Expected next level at 500 points, got %d", nextLevelPoints)
	}
}

// Test 7: Test user_point_logs table for daily limit tracking
func TestPointService_LogPointChange(t *testing.T) {
	db := setupFullTestDB(t)

	// Create test user
	user := domain.User{Nickname: "testuser6", Role: "member"}
	db.Create(&user)

	pointRepo := repository.NewUserPointRepository(db)

	// Log a point change
	err := pointRepo.LogPointChange(user.ID, 10, "post_create", 123)
	if err != nil {
		t.Errorf("Failed to log point change: %v", err)
	}

	// Verify log was created
	var count int64
	db.Table("user_point_logs").Where("user_id = ?", user.ID).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 log entry, got %d", count)
	}

	// Check daily points
	todayPoints, err := pointRepo.GetTodayPostPoints(user.ID)
	if err != nil {
		t.Errorf("Failed to get today points: %v", err)
	}
	if todayPoints != 10 {
		t.Errorf("Expected 10 points today, got %d", todayPoints)
	}
}

// Test 8: Test hot posts calculation
func TestHotPostsService_GetHotPosts(t *testing.T) {
	db := setupFullTestDB(t)

	// Create test user
	user := domain.User{Nickname: "testuser7", Role: "member"}
	db.Create(&user)

	// Create test posts with different engagement
	now := time.Now()
	posts := []domain.Post{
		{
			UserID:     user.ID,
			Title:      "Post 1",
			Content:    "Content 1",
			Category:   "discussion",
			Visibility: "public",
			Status:     "approved",
			ViewCount:  100,
			ReplyCount: 10,
			LikeCount:  5,
			CreatedAt:  now.AddDate(0, 0, -3), // 3 days ago
		},
		{
			UserID:     user.ID,
			Title:      "Post 2",
			Content:    "Content 2",
			Category:   "discussion",
			Visibility: "public",
			Status:     "approved",
			ViewCount:  50,
			ReplyCount: 5,
			LikeCount:  20,
			CreatedAt:  now.AddDate(0, 0, -5), // 5 days ago
		},
	}

	for _, p := range posts {
		db.Create(&p)
	}

	// Get hot posts
	postRepo := repository.NewPostRepository(db)
	hotPostsService := NewHotPostsService(postRepo)

	result, total, err := hotPostsService.GetHotPosts(1, 10)
	if err != nil {
		t.Errorf("Failed to get hot posts: %v", err)
	}

	if total != 2 {
		t.Errorf("Expected 2 hot posts, got %d", total)
	}

	// Post 1 hotness: 10*2 + 5 + 100/10 = 35
	// Post 2 hotness: 5*2 + 20 + 50/10 = 25
	// So Post 1 should be first
	if len(result) > 0 && result[0].Title != "Post 1" {
		t.Errorf("Expected 'Post 1' to be first, got %s", result[0].Title)
	}
}
