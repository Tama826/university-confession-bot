package main

import (
	"strings"
)

// ---------------- KEYWORD & TOXICITY ----------------

func containsBannedWord(text string) bool {
	text = strings.ToLower(text)
	for _, w := range BannedWords {
		if strings.Contains(text, w) {
			return true
		}
	}
	return false
}

func toxicityScore(text string) int {
	score := 0
	text = strings.ToLower(text)
	badWords := []string{"kill", "rape", "die", "bomb", "terror"}
	for _, w := range badWords {
		if strings.Contains(text, w) {
			score += 30
		}
	}
	return score
}
