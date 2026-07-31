package splitter

import "strings"

type Chunk struct {
	Content string
	Index   int
}

type TextSplitter interface {
	Split(text string) []Chunk
}

// TokenSplitter splits text by approximate token count (1 token ≈ 4 chars)
type TokenSplitter struct {
	ChunkSize         int
	ChunkOverlap      int
	MinChunkSizeChars int
}

func NewTokenSplitter(chunkSize, overlap, minChars int) *TokenSplitter {
	if chunkSize <= 0 {
		chunkSize = 1000
	}
	if minChars <= 0 {
		minChars = 400
	}
	return &TokenSplitter{ChunkSize: chunkSize, ChunkOverlap: overlap, MinChunkSizeChars: minChars}
}

func (s *TokenSplitter) Split(text string) []Chunk {
	charsPerToken := 4
	maxChars := s.ChunkSize * charsPerToken
	overlapChars := s.ChunkOverlap * charsPerToken
	var chunks []Chunk
	idx := 0
	for i := 0; i < len(text); i += maxChars - overlapChars {
		end := i + maxChars
		if end > len(text) {
			end = len(text)
		}
		chunk := text[i:end]
		if len(chunk) >= s.MinChunkSizeChars {
			chunks = append(chunks, Chunk{Content: chunk, Index: idx})
			idx++
		}
	}
	return chunks
}

// RecursiveSplitter splits by separators recursively
type RecursiveSplitter struct {
	Separators   []string
	ChunkSize    int
	ChunkOverlap int
}

func NewRecursiveSplitter(chunkSize, overlap int) *RecursiveSplitter {
	return &RecursiveSplitter{
		Separators:   []string{"\n\n", "\n", "。", ".", " ", ""},
		ChunkSize:    chunkSize,
		ChunkOverlap: overlap,
	}
}

func (s *RecursiveSplitter) Split(text string) []Chunk {
	return s.splitText(text, s.Separators, 0)
}

func (s *RecursiveSplitter) splitText(text string, seps []string, idx int) []Chunk {
	if len(seps) == 0 {
		return []Chunk{{Content: text, Index: idx}}
	}
	sep := seps[0]
	parts := strings.Split(text, sep)
	var chunks []Chunk
	current := ""
	for _, part := range parts {
		if len(current)+len(part)+len(sep) > s.ChunkSize && current != "" {
			chunks = append(chunks, Chunk{Content: current, Index: idx})
			idx++
			if s.ChunkOverlap > 0 && len(current) > s.ChunkOverlap {
				current = current[len(current)-s.ChunkOverlap:]
			} else {
				current = ""
			}
		}
		if current != "" {
			current += sep
		}
		current += part
	}
	if current != "" {
		chunks = append(chunks, Chunk{Content: current, Index: idx})
	}
	return chunks
}

// SentenceSplitter splits by sentence boundaries
type SentenceSplitter struct {
	ChunkSize       int
	SentenceOverlap int
}

func NewSentenceSplitter(chunkSize, overlap int) *SentenceSplitter {
	return &SentenceSplitter{ChunkSize: chunkSize, SentenceOverlap: overlap}
}

func (s *SentenceSplitter) Split(text string) []Chunk {
	sentences := splitSentences(text)
	var chunks []Chunk
	current := ""
	idx := 0
	for _, sent := range sentences {
		if len(current)+len(sent) > s.ChunkSize && current != "" {
			chunks = append(chunks, Chunk{Content: current, Index: idx})
			idx++
			if s.SentenceOverlap > 0 {
				words := strings.Fields(current)
				start := len(words) - s.SentenceOverlap
				if start < 0 {
					start = 0
				}
				current = strings.Join(words[start:], " ") + " "
			} else {
				current = ""
			}
		}
		current += sent
	}
	if current != "" {
		chunks = append(chunks, Chunk{Content: current, Index: idx})
	}
	return chunks
}

func splitSentences(text string) []string {
	var sentences []string
	current := ""
	for _, ch := range text {
		current += string(ch)
		if ch == '。' || ch == '.' || ch == '！' || ch == '?' || ch == '\n' {
			sentences = append(sentences, current)
			current = ""
		}
	}
	if current != "" {
		sentences = append(sentences, current)
	}
	return sentences
}

// ParagraphSplitter splits by paragraphs
type ParagraphSplitter struct {
	ChunkSize        int
	ParagraphOverlap int
}

func NewParagraphSplitter(chunkSize, overlap int) *ParagraphSplitter {
	return &ParagraphSplitter{ChunkSize: chunkSize, ParagraphOverlap: overlap}
}

func (s *ParagraphSplitter) Split(text string) []Chunk {
	paragraphs := strings.Split(text, "\n\n")
	var chunks []Chunk
	current := ""
	idx := 0
	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}
		if len(current)+len(para)+2 > s.ChunkSize && current != "" {
			chunks = append(chunks, Chunk{Content: current, Index: idx})
			idx++
			current = ""
		}
		if current != "" {
			current += "\n\n"
		}
		current += para
	}
	if current != "" {
		chunks = append(chunks, Chunk{Content: current, Index: idx})
	}
	return chunks
}

// SemanticSplitter splits based on embedding similarity (stub - requires embedder)
type SemanticSplitter struct {
	MinChunkSize int
	MaxChunkSize int
	Threshold    float64
}

func NewSemanticSplitter(minSize, maxSize int, threshold float64) *SemanticSplitter {
	return &SemanticSplitter{MinChunkSize: minSize, MaxChunkSize: maxSize, Threshold: threshold}
}

func (s *SemanticSplitter) Split(text string) []Chunk {
	ps := NewParagraphSplitter(s.MaxChunkSize, 0)
	return ps.Split(text)
}

// Factory
func NewSplitter(name string) TextSplitter {
	switch name {
	case "token":
		return NewTokenSplitter(1000, 0, 400)
	case "recursive":
		return NewRecursiveSplitter(1000, 200)
	case "sentence":
		return NewSentenceSplitter(1000, 1)
	case "semantic":
		return NewSemanticSplitter(200, 1000, 0.5)
	case "paragraph":
		return NewParagraphSplitter(1000, 200)
	default:
		return NewTokenSplitter(1000, 0, 400)
	}
}
