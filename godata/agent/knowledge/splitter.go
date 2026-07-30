package knowledge

import (
	"unicode/utf8"
)

// Split splits a text into overlapping chunks of approximately chunkSize
// characters, each overlapping by overlap characters with the previous chunk.
//
// The algorithm prefers splitting on paragraph boundaries (double newline),
// then sentence boundaries (period + space), then word boundaries (space),
// and finally falls back to hard character-level truncation.
//
// Returns nil for empty input.
func Split(text string, chunkSize, overlap int) []string {
	if text == "" || chunkSize <= 0 {
		return nil
	}
	if overlap >= chunkSize {
		overlap = chunkSize / 2
	}
	if overlap < 0 {
		overlap = 0
	}

	// Fast path: text is shorter than one chunk.
	if utf8.RuneCountInString(text) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	runes := []rune(text)
	start := 0

	for start < len(runes) {
		end := start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}

		chunk := string(runes[start:end])

		// Only try boundary-aware splitting when we are not at the very end.
		if end < len(runes) {
			chunk = splitOnBoundary(chunk, chunkSize)
		}

		chunks = append(chunks, chunk)

		// Advance by (chunkSize - overlap) characters.
		advance := chunkSize - overlap
		if advance <= 0 {
			advance = 1
		}
		start += utf8.RuneCountInString(chunk) - overlap
		if start < 0 {
			start = 0
		}
	}
	return chunks
}

// splitOnBoundary attempts to trim a chunk to a natural text boundary.
func splitOnBoundary(chunk string, maxSize int) string {
	runes := []rune(chunk)

	// Try paragraph break (double newline).
	if idx := lastIndexRune(runes, "\n\n"); idx >= 0 {
		return string(runes[:idx+2])
	}

	// Try single newline.
	if idx := lastIndexRune(runes, "\n"); idx >= 0 {
		return string(runes[:idx+1])
	}

	// Try sentence boundary (. followed by space or end).
	if idx := lastIndexRune(runes, ". "); idx >= 0 {
		return string(runes[:idx+2])
	}
	if idx := lastIndexRune(runes, "。"); idx >= 0 {
		return string(runes[:idx+1])
	}

	// Try space (word boundary).
	if idx := lastIndexRune(runes, " "); idx >= maxSize/2 {
		return string(runes[:idx])
	}

	return chunk
}

func lastIndexRune(runes []rune, sep string) int {
	sepRunes := []rune(sep)
	if len(sepRunes) == 0 || len(runes) < len(sepRunes) {
		return -1
	}
	for i := len(runes) - len(sepRunes); i >= 0; i-- {
		match := true
		for j, r := range sepRunes {
			if runes[i+j] != r {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}
