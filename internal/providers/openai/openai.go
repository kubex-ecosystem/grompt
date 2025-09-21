// Package openai implements the OpenAI API client.
package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	providers "github.com/kubex-ecosystem/grompt/internal/types"
)

// Provider OpenAI com streaming + tool calling + response_format + usage (stream_options)
type Provider struct {
	name         string
	baseURL      string
	apiKey       string
	defaultModel string
}

func New(name, baseURL, apiKey, defaultModel string) (*Provider, error) {
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY vazio")
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}
	return &Provider{name: name, baseURL: strings.TrimRight(baseURL, "/"), apiKey: apiKey, defaultModel: defaultModel}, nil
}

func (p *Provider) Name() string { return p.name }

func (p *Provider) Chat(ctx context.Context, in providers.ChatRequest) (<-chan providers.ChatChunk, error) {
	model := in.Model
	if model == "" {
		model = p.defaultModel
	}
	body := map[string]any{
		"model":       model,
		"temperature": in.Temp,
		"stream":      true,
		// IMPORTANT: inclui usage no stream.
		"stream_options": map[string]any{"include_usage": true},
		"messages":       toOAIMessages(in.Messages),
	}

	// Tools / tool_choice (se vier no Meta)
	if tools, ok := in.Meta["tools"]; ok {
		body["tools"] = tools // esperado no formato OpenAI: [{"type":"function","function":{name,description,parameters}}]
	}
	if tc, ok := in.Meta["tool_choice"]; ok {
		body["tool_choice"] = tc // "auto" | {"type":"function","function":{"name":"..."}}
	} else if _, ok := in.Meta["tools"]; ok {
		body["tool_choice"] = "auto"
	}

	// Response format (structured output) — ex.: {"type":"json_schema","json_schema":{...}}
	if rf, ok := in.Meta["response_format"]; ok {
		body["response_format"] = rf
	}

	b, _ := json.Marshal(body)
	url := p.baseURL + "/v1/chat/completions"

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	// BYOK tem prioridade
	key := in.Headers["x-external-api-key"]
	if key == "" {
		key = p.apiKey
	}
	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		defer resp.Body.Close()
		var data map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&data)
		return nil, fmt.Errorf("openai %d: %v", resp.StatusCode, data)
	}

	out := make(chan providers.ChatChunk, 16)
	go func() {
		defer close(out)
		defer resp.Body.Close()
		start := time.Now()

		// Acumuladores para tool calls (por índice)
		type toolBuf struct {
			Name string
			Args strings.Builder // args chegam em fragmentos
		}
		tools := map[int]*toolBuf{}
		var usage providers.Usage

		sc := bufio.NewScanner(resp.Body)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" {
				continue
			}
			if !strings.HasPrefix(line, "data:") {
				continue
			}
			payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if payload == "[DONE]" {
				break
			}

			// Dois tipos de eventos: chunks (choices[].delta...) e usage (usage no topo)
			// 1) Tenta usage
			var maybeUsage struct {
				Usage *struct {
					PromptTokens     int `json:"prompt_tokens"`
					CompletionTokens int `json:"completion_tokens"`
					TotalTokens      int `json:"total_tokens"`
				} `json:"usage"`
			}
			if json.Unmarshal([]byte(payload), &maybeUsage) == nil && maybeUsage.Usage != nil {
				usage.Prompt = maybeUsage.Usage.PromptTokens
				usage.Completion = maybeUsage.Usage.CompletionTokens
				usage.Tokens = maybeUsage.Usage.TotalTokens
				continue
			}

			// 2) Chunk de delta
			var ev struct {
				Choices []struct {
					Delta struct {
						Content   string `json:"content"`
						ToolCalls []struct {
							Index    int    `json:"index"`
							ID       string `json:"id,omitempty"`
							Type     string `json:"type,omitempty"`
							Function struct {
								Name      string `json:"name,omitempty"`
								Arguments string `json:"arguments,omitempty"`
							} `json:"function"`
						} `json:"tool_calls,omitempty"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(payload), &ev); err != nil || len(ev.Choices) == 0 {
				continue
			}
			ch := ev.Choices[0]

			// Conteúdo textual
			if txt := ch.Delta.Content; txt != "" {
				out <- providers.ChatChunk{Content: txt}
			}

			// Tool calling (fragmentado)
			if len(ch.Delta.ToolCalls) > 0 {
				for _, t := range ch.Delta.ToolCalls {
					tb := tools[t.Index]
					if tb == nil {
						tb = &toolBuf{Name: t.Function.Name}
						tools[t.Index] = tb
					}
					if t.Function.Name != "" {
						tb.Name = t.Function.Name
					}
					if t.Function.Arguments != "" {
						tb.Args.WriteString(t.Function.Arguments)
					}
				}
			}

			// Se finish_reason == "tool_calls", emite cada tool call agregada
			if ch.FinishReason == "tool_calls" && len(tools) > 0 {
				for _, tb := range tools {
					var args any
					argStr := strings.TrimSpace(tb.Args.String())
					if argStr != "" {
						// tentar parsear JSON; se falhar, manda string crua
						if json.Unmarshal([]byte(argStr), &args) != nil {
							args = argStr
						}
					}
					out <- providers.ChatChunk{
						ToolCall: &providers.ToolCall{Name: tb.Name, Args: args},
					}
				}
				// limpa para não reenviar
				tools = map[int]*toolBuf{}
			}
		}

		usage.Ms = time.Since(start).Milliseconds()
		out <- providers.ChatChunk{Done: true, Usage: &usage}
	}()

	return out, nil
}

func toOAIMessages(ms []providers.Message) []map[string]any {
	out := make([]map[string]any, 0, len(ms))
	for _, m := range ms {
		out = append(out, map[string]any{
			"role":    m.Role,
			"content": m.Content,
		})
	}
	return out
}
