---
title: "SSE Streaming Architecture with Coalescing - Grompt V1"
version: 0.2.1
owner: kubex
audience: dev
languages: [en, pt-BR]
sources: ["internal/gateway/transport/sse_coalescer.go", "internal/gateway/transport/grompt_v1.go"]
assumptions: ["Browser SSE support", "Network stability for real-time streaming"]
---

# SSE Streaming Architecture with Coalescing - Grompt V1

## TL;DR

Grompt V1 implements an advanced Server-Sent Events (SSE) streaming system with intelligent chunk coalescing for optimal user experience. The system buffers micro-chunks from AI providers and delivers content at natural language boundaries, creating smooth, readable streaming while maintaining real-time responsiveness.

## Architecture Overview

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                    Client (Browser)                         │
│                EventSource API                              │
├─────────────────────────────────────────────────────────────┤
│                  Grompt V1 Gateway                          │
│              /v1/generate/stream                            │
├─────────────────────────────────────────────────────────────┤
│                  SSE Coalescer                              │
│              (Intelligent Buffering)                        │
├─────────────────────────────────────────────────────────────┤
│   OpenAI Stream  │  Anthropic Stream  │  Gemini Stream    │
│   (Micro-chunks) │   (Token deltas)   │  (Word fragments) │
└─────────────────────────────────────────────────────────────┘
```

### SSE Coalescer Implementation

#### Core Algorithm (`internal/gateway/transport/sse_coalescer.go`)

The SSE Coalescer implements a sophisticated buffering strategy:

```go
type SSECoalescer struct {
    buffer        strings.Builder
    flushTimer    *time.Timer
    flushFunc     func(content string)
    bufferTimeout time.Duration  // 75ms sweet spot
    maxBufferSize int           // 100 characters max
}
```

#### Flush Trigger Logic

The coalescer flushes content based on multiple conditions:

1. **Natural Language Boundaries**
   - Sentence endings: `.`, `!`, `?`
   - Phrase boundaries: `,`, `;`, `:`
   - Natural pauses: spaces after meaningful content

2. **Buffer Management**
   - Maximum buffer size (100 characters)
   - Timeout-based flushing (75ms)
   - Immediate newline preservation

3. **Performance Optimization**
   - Prevents micro-chunk spam
   - Maintains readability during streaming
   - Balances latency vs. smoothness

```go
func (c *SSECoalescer) AddChunk(content string) {
    c.buffer.WriteString(content)

    shouldFlushNow := c.buffer.Len() >= c.maxBufferSize ||
        c.hasNaturalBreak(content) ||
        strings.Contains(content, "\n")

    if shouldFlushNow {
        c.flushNow()
        return
    }

    c.resetFlushTimer()
}
```

## Streaming Endpoint Implementation

### `/v1/generate/stream` Handler

The streaming endpoint in `internal/gateway/transport/grompt_v1.go` integrates the coalescer:

```go
func (h *GromptV1Handlers) generatePromptStream(w http.ResponseWriter, r *http.Request) {
    // Set SSE headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    // Create SSE coalescer for smooth streaming
    coalescer := NewSSECoalescer(func(content string) {
        data := mustMarshalJSON(map[string]interface{}{
            "event": "generation.chunk",
            "content": content,
        })
        fmt.Fprintf(w, "data: %s\n\n", data)
        flusher.Flush()
    })
    defer coalescer.Close()

    // Stream with coalescence
    for chunk := range responseChannel {
        if chunk.Content != "" {
            coalescer.AddChunk(chunk.Content)
        }
    }
}
```

### Event Format

The streaming endpoint follows a standardized SSE event format:

```javascript
// Streaming chunk
data: {
  "event": "generation.chunk",
  "content": "This is a coalesced chunk of content",
  "timestamp": "2024-01-01T12:00:00Z"
}

// Completion event
data: {
  "event": "generation.complete",
  "usage": {
    "tokens": 150,
    "cost_usd": 0.003
  }
}

// Error event
data: {
  "event": "generation.error",
  "error": "Provider temporarily unavailable"
}
```

## Provider-Specific Streaming

### OpenAI Streaming Integration

OpenAI provides delta-based streaming with fine-grained chunks:

```go
// OpenAI streaming characteristics
- Chunk size: 1-5 characters typically
- Frequency: High (multiple chunks per second)
- Content: Token fragments and partial words
- Coalescing benefit: High (reduces UI jitter)
```

### Anthropic (Claude) Streaming

Anthropic delivers more substantial chunks but still benefits from coalescing:

```go
// Anthropic streaming characteristics
- Chunk size: 5-20 characters typically
- Frequency: Medium (controlled pace)
- Content: Word and phrase boundaries
- Coalescing benefit: Medium (smooths phrases)
```

### Gemini Streaming

Google's Gemini provides varied chunk sizes with safety filtering:

```go
// Gemini streaming characteristics
- Chunk size: Variable (1-50 characters)
- Frequency: Variable (depends on content complexity)
- Content: Complete words with safety filtering
- Coalescing benefit: High (handles variable patterns)
```

## Performance Optimizations

### Buffer Management

```go
const (
    OptimalBufferTimeout = 75 * time.Millisecond
    MaxBufferSize       = 100 // characters
    MinFlushSize        = 10  // minimum chars before timeout flush
)
```

**Rationale**:
- **75ms timeout**: Balances perceived responsiveness with reduced chunk frequency
- **100 char buffer**: Prevents excessive buffering while allowing natural phrases
- **Natural boundaries**: Preserves linguistic coherence

### Memory Efficiency

The coalescer uses several memory optimization techniques:

1. **strings.Builder**: Efficient string concatenation
2. **Timer pooling**: Reuses timer objects
3. **Immediate cleanup**: Buffers reset after each flush
4. **Bounded growth**: Maximum buffer size prevents memory leaks

### Concurrent Request Handling

The streaming system supports multiple concurrent streams:

```go
// Per-request coalescer instances
func (h *GromptV1Handlers) generatePromptStream(w http.ResponseWriter, r *http.Request) {
    // Each request gets its own coalescer
    coalescer := NewSSECoalescer(createFlushFunc(w))

    // Goroutine-safe operation
    go func() {
        defer coalescer.Close()
        // Handle streaming...
    }()
}
```

## Frontend Integration

### EventSource Client Implementation

The frontend consumes the SSE stream using the EventSource API:

```typescript
const eventSource = new EventSource('/v1/generate/stream?' + params)

eventSource.onmessage = (event) => {
  const data = JSON.parse(event.data)

  switch (data.event) {
    case 'generation.chunk':
      appendContent(data.content)
      break
    case 'generation.complete':
      handleCompletion(data.usage)
      break
    case 'generation.error':
      handleError(data.error)
      break
  }
}
```

### React Integration

The React frontend provides hooks for streaming integration:

```typescript
const useStreamingGeneration = () => {
  const [content, setContent] = useState('')
  const [isStreaming, setIsStreaming] = useState(false)

  const generateStream = async (request: GenerateRequest) => {
    setIsStreaming(true)
    setContent('')

    await enhancedAPI.generatePromptStream(
      request,
      (chunk) => setContent(prev => prev + chunk),
      (usage) => {
        setIsStreaming(false)
        console.log('Stream completed:', usage)
      },
      (error) => {
        setIsStreaming(false)
        console.error('Stream error:', error)
      }
    )
  }

  return { content, isStreaming, generateStream }
}
```

## Quality Assurance

### Natural Language Preservation

The coalescer respects linguistic boundaries:

```go
func (c *SSECoalescer) hasNaturalBreak(content string) bool {
    lastChar := content[len(content)-1]
    return lastChar == '.' || lastChar == '!' || lastChar == '?' ||
           lastChar == ',' || lastChar == ';' || lastChar == ':' ||
           lastChar == ' ' && c.buffer.Len() > 20
}
```

### Error Handling and Recovery

```go
// Graceful error handling in streaming
func (h *GromptV1Handlers) handleStreamError(w http.ResponseWriter, err error) {
    // Flush any pending content before error
    coalescer.Close()

    errorData := mustMarshalJSON(map[string]interface{}{
        "event": "generation.error",
        "error": err.Error(),
    })
    fmt.Fprintf(w, "data: %s\n\n", errorData)
}
```

### Connection Management

The system handles connection drops and client disconnects:

```go
// Detect client disconnect
select {
case <-r.Context().Done():
    // Client disconnected, cleanup resources
    coalescer.Close()
    return
case chunk := <-responseChannel:
    // Continue processing
}
```

## Monitoring and Observability

### Streaming Metrics

The system collects comprehensive streaming metrics:

```go
type StreamingMetrics struct {
    ChunksReceived    int64
    ChunksCoalesced   int64
    AvgCoalescingRate float64
    StreamDuration    time.Duration
    BytesTransferred  int64
}
```

### Performance Monitoring

```go
// Record streaming performance
h.recordMetrics("/v1/generate/stream", provider, model,
                duration, tokens, cost, nil)
```

## How to Run / Repro

### Test Streaming Endpoint

```bash
# Test basic streaming
curl -N -H "Accept: text/event-stream" \
  "http://localhost:3000/v1/generate/stream?provider=openai&ideas=AI&ideas=streaming&purpose=general"

# Expected output:
# data: {"event":"generation.chunk","content":"Artificial Intelligence and streaming"}
# data: {"event":"generation.chunk","content":" technologies are revolutionizing"}
# data: {"event":"generation.complete","usage":{"tokens":45}}
```

### Frontend Streaming Test

```javascript
// Test in browser console
const testStreaming = () => {
  const eventSource = new EventSource('/v1/generate/stream?provider=openai&ideas=test&purpose=general')

  eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data)
    console.log(`Event: ${data.event}, Content: ${data.content}`)
  }

  setTimeout(() => eventSource.close(), 30000) // Close after 30s
}

testStreaming()
```

### Performance Benchmarking

```bash
# Benchmark coalescing performance
go test -bench=BenchmarkSSECoalescer ./internal/gateway/transport/

# Expected results:
# BenchmarkSSECoalescer-8    1000000    1543 ns/op    64 B/op    2 allocs/op
```

## Risks & Mitigations

### Risk: Buffer Memory Leaks
**Mitigation**:
- Bounded buffer size (100 characters max)
- Automatic cleanup with defer statements
- Timer cancellation on early termination

### Risk: Streaming Latency
**Mitigation**:
- 75ms timeout balances latency vs. smoothness
- Natural boundary detection for immediate important content
- Bypass coalescing for newlines and punctuation

### Risk: Client Disconnect Handling
**Mitigation**:
- Context cancellation detection
- Resource cleanup on disconnect
- Graceful error propagation

### Risk: Provider Stream Variability
**Mitigation**:
- Adaptive coalescing based on content patterns
- Provider-specific timeout adjustments
- Fallback to non-streaming on failure

## Next Steps

1. **Adaptive Coalescing**: Dynamic timeout adjustment based on content type
2. **Compression**: Gzip compression for large streams
3. **Multiplexing**: Multiple concurrent streams per connection
4. **Analytics**: User experience metrics for optimal tuning
5. **Internationalization**: Language-specific boundary detection

## Changelog

- **v0.2.1**: Advanced SSE coalescing with natural language boundaries
- **v0.2.0**: Basic SSE streaming implementation
- **v0.1.0**: Foundation streaming architecture