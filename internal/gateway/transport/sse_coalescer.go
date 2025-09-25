package transport

import (
	"strings"
	"time"
)

// SSECoalescer buffers streaming chunks to improve UX by reducing micro-chunks
type SSECoalescer struct {
	buffer        strings.Builder
	flushTimer    *time.Timer
	flushFunc     func(content string)
	bufferTimeout time.Duration
	maxBufferSize int
}

// NewSSECoalescer creates a new coalescer with optimized defaults
func NewSSECoalescer(flushFunc func(content string)) *SSECoalescer {
	return &SSECoalescer{
		flushFunc:     flushFunc,
		bufferTimeout: 75 * time.Millisecond, // Sweet spot for UX vs efficiency
		maxBufferSize: 100,                   // Max chars to buffer
	}
}

// AddChunk adds content to buffer and triggers flush based on natural boundaries
func (c *SSECoalescer) AddChunk(content string) {
	if content == "" {
		return
	}

	c.buffer.WriteString(content)

	// Immediate flush conditions:
	// 1. Buffer getting too large
	// 2. Natural punctuation boundaries
	// 3. Newlines (preserve line breaks)
	shouldFlushNow := c.buffer.Len() >= c.maxBufferSize ||
		c.hasNaturalBreak(content) ||
		strings.Contains(content, "\n")

	if shouldFlushNow {
		c.flushNow()
		return
	}

	// Set/reset flush timer for accumulated content
	c.resetFlushTimer()
}

// hasNaturalBreak checks if content ends with natural punctuation
func (c *SSECoalescer) hasNaturalBreak(content string) bool {
	if len(content) == 0 {
		return false
	}

	// Natural break points for smooth reading
	lastChar := content[len(content)-1]
	return lastChar == '.' || lastChar == '!' || lastChar == '?' ||
		lastChar == ',' || lastChar == ';' || lastChar == ':' ||
		lastChar == ' ' && c.buffer.Len() > 20 // Space after reasonable content
}

// resetFlushTimer cancels existing timer and sets new one
func (c *SSECoalescer) resetFlushTimer() {
	if c.flushTimer != nil {
		c.flushTimer.Stop()
	}

	c.flushTimer = time.AfterFunc(c.bufferTimeout, func() {
		c.flushNow()
	})
}

// flushNow immediately sends buffered content and resets buffer
func (c *SSECoalescer) flushNow() {
	if c.flushTimer != nil {
		c.flushTimer.Stop()
		c.flushTimer = nil
	}

	if c.buffer.Len() > 0 {
		content := c.buffer.String()
		c.buffer.Reset()
		c.flushFunc(content)
	}
}

// Close flushes any remaining content and cleans up
func (c *SSECoalescer) Close() {
	c.flushNow()
}
