package internal

import "strings"

func EstimateTokens(text string) int {
    words := strings.Fields(text)
	return int(float64(len(words)) * 1.3)
}
