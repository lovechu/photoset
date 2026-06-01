package service

import (
	"regexp"
	"strings"

	"photoset/internal/domain"
	"photoset/internal/repository"
)

// MentionService provides @mention business logic
type MentionService struct {
	userRepo         repository.UserRepository
	notificationRepo *repository.NotificationRepository
}

// NewMentionService creates a new MentionService
func NewMentionService(
	userRepo repository.UserRepository,
	notificationRepo *repository.NotificationRepository,
) *MentionService {
	return &MentionService{
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
	}
}

// mentionRegex matches @username patterns in content
// Supports: @username, @"user name", @用户昵称
var mentionRegex = regexp.MustCompile(`@(?:"([^"]+)"|([^\s@]+))`)

// ParseMentions extracts @mentioned usernames from content
func (s *MentionService) ParseMentions(content string) []string {
	matches := mentionRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	// Deduplicate usernames
	seen := make(map[string]bool)
	var usernames []string
	for _, match := range matches {
		var username string
		if match[1] != "" {
			username = match[1] // Quoted username
		} else {
			username = match[2] // Unquoted username
		}
		username = strings.TrimSpace(username)
		if username != "" && !seen[username] {
			seen[username] = true
			usernames = append(usernames, username)
		}
	}

	return usernames
}

// ResolveMentions resolves mentioned usernames to user IDs
func (s *MentionService) ResolveMentions(usernames []string) ([]domain.User, error) {
	if len(usernames) == 0 {
		return nil, nil
	}

	// Search for users by nickname
	var allUsers []domain.User
	for _, username := range usernames {
		users, err := s.userRepo.SearchByNickname(username, 5)
		if err != nil {
			continue // Skip on error
		}
		// Find exact match (case-insensitive)
		for _, user := range users {
			if strings.EqualFold(user.Nickname, username) {
				allUsers = append(allUsers, user)
				break
			}
		}
	}

	return allUsers, nil
}

// SendMentionNotifications sends notifications to mentioned users
func (s *MentionService) SendMentionNotifications(senderID, postID uint, content string) {
	// Parse mentions from content
	usernames := s.ParseMentions(content)
	if len(usernames) == 0 {
		return
	}

	// Resolve usernames to users
	users, err := s.ResolveMentions(usernames)
	if err != nil || len(users) == 0 {
		return
	}

	// Send notification to each mentioned user
	for _, user := range users {
		if user.ID == senderID {
			continue // Don't notify self
		}

		notification := &domain.Notification{
			UserID:     user.ID,
			Type:       domain.NotificationTypeMention,
			Title:      "被提及",
			Content:    "有人在帖子中提到了你",
			SenderID:   &senderID,
			TargetID:   &postID,
			TargetType: "post",
		}

		// Ignore error - notification is non-critical
		s.notificationRepo.Create(notification)
	}
}

// GetMentionableUsers returns users matching the search keyword for @mention autocomplete
func (s *MentionService) GetMentionableUsers(keyword string, limit int) ([]map[string]interface{}, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	users, err := s.userRepo.SearchByNickname(keyword, limit)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(users))
	for i, user := range users {
		result[i] = map[string]interface{}{
			"id":       user.ID,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		}
	}

	return result, nil
}

// FormatMentionedContent converts @mentions in content to clickable links
// This is used for frontend rendering
func (s *MentionService) FormatMentionedContent(content string) string {
	// Replace @username with markdown-style links
	return mentionRegex.ReplaceAllStringFunc(content, func(match string) string {
		submatches := mentionRegex.FindStringSubmatch(match)
		var username string
		if submatches[1] != "" {
			username = submatches[1]
		} else {
			username = submatches[2]
		}
		return "@[" + username + "]"
	})
}
