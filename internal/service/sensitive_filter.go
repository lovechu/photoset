package service

import (
	"strings"
	"sync"

	"photoset/internal/repository"
)

var (
	// SensitiveWordsMap stores sensitive words in memory (lowercase word -> replacement)
	SensitiveWordsMap sync.Map
)

// ClearSensitiveWordsMap clears all sensitive words from memory
func ClearSensitiveWordsMap() {
	SensitiveWordsMap.Range(func(key, value interface{}) bool {
		SensitiveWordsMap.Delete(key)
		return true
	})
}

// SensitiveFilterService provides sensitive word filtering
type SensitiveFilterService struct {
	wordRepo *repository.SensitiveWordRepository
}

// NewSensitiveFilterService creates a new SensitiveFilterService
func NewSensitiveFilterService(wordRepo *repository.SensitiveWordRepository) *SensitiveFilterService {
	return &SensitiveFilterService{wordRepo: wordRepo}
}

// LoadSensitiveWords loads all active sensitive words from database to memory
func (s *SensitiveFilterService) LoadSensitiveWords() error {
	words, err := s.wordRepo.LoadAllActive()
	if err != nil {
		return err
	}

	// Clear existing map by creating a new map (Range+Delete can miss entries)
	SensitiveWordsMap = sync.Map{}

	// Load new words (store lowercase)
	for _, w := range words {
		SensitiveWordsMap.Store(strings.ToLower(w.Word), w.Replacement)
	}

	return nil
}

// FilterText filters sensitive words from text (case-insensitive, preserves original case)
// Returns the filtered text and number of replacements made
// Deprecated: Use FilterTextAdvanced for better case preservation
func (s *SensitiveFilterService) FilterText(text string) (string, int) {
	return s.FilterTextAdvanced(text)
}

// FilterTextAdvanced filters sensitive words while preserving original case
func (s *SensitiveFilterService) FilterTextAdvanced(text string) (string, int) {
	if text == "" {
		return text, 0
	}

	result := text
	replacementCount := 0

	SensitiveWordsMap.Range(func(key, value interface{}) bool {
		word := key.(string)
		replacement := value.(string)

		// Case-insensitive replacement
		count := strings.Count(strings.ToLower(result), word)
		if count > 0 {
			replacementCount += count
			// Use case-insensitive replacement
			result = replaceAllCaseInsensitive(result, word, replacement)
		}

		return true
	})

	return result, replacementCount
}

// replaceAllCaseInsensitive replaces all occurrences of old with new (case-insensitive)
func replaceAllCaseInsensitive(s, old, new string) string {
	result := s
	lowerS := strings.ToLower(s)
	lowerOld := strings.ToLower(old)

	for {
		idx := strings.Index(lowerS, lowerOld)
		if idx == -1 {
			break
		}

		// Replace preserving case of surrounding text
		result = result[:idx] + new + result[idx+len(old):]
		lowerS = strings.ToLower(result)
	}

	return result
}

// InitSensitiveWords initializes sensitive words on service startup
func InitSensitiveWords(wordRepo *repository.SensitiveWordRepository) error {
	service := NewSensitiveFilterService(wordRepo)
	return service.LoadSensitiveWords()
}
