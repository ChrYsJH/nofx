package mcp

import (
	"net/http"
)

const (
	ProviderOpenRouter       = "openrouter"
	DefaultOpenRouterBaseURL = "https://openrouter.ai/api/v1"
	DefaultOpenRouterModel   = "openai/gpt-5.2" // OpenRouter 使用 provider/model 格式
)

type OpenRouterClient struct {
	*Client
}

// NewOpenRouterClient creates OpenRouter client (backward compatible)
func NewOpenRouterClient() AIClient {
	return NewOpenRouterClientWithOptions()
}

// NewOpenRouterClientWithOptions creates OpenRouter client (supports options pattern)
func NewOpenRouterClientWithOptions(opts ...ClientOption) AIClient {
	// 1. Create OpenRouter preset options
	openRouterOpts := []ClientOption{
		WithProvider(ProviderOpenRouter),
		WithModel(DefaultOpenRouterModel),
		WithBaseURL(DefaultOpenRouterBaseURL),
	}

	// 2. Merge user options (user options have higher priority)
	allOpts := append(openRouterOpts, opts...)

	// 3. Create base client
	baseClient := NewClient(allOpts...).(*Client)

	// 4. Create OpenRouter client
	openRouterClient := &OpenRouterClient{
		Client: baseClient,
	}

	// 5. Set hooks to point to OpenRouterClient (implement dynamic dispatch)
	baseClient.hooks = openRouterClient

	return openRouterClient
}

func (c *OpenRouterClient) SetAPIKey(apiKey string, customURL string, customModel string) {
	c.APIKey = apiKey

	if len(apiKey) > 8 {
		c.logger.Infof("🔧 [MCP] OpenRouter API Key: %s...%s", apiKey[:4], apiKey[len(apiKey)-4:])
	}
	if customURL != "" {
		c.BaseURL = customURL
		c.logger.Infof("🔧 [MCP] OpenRouter using custom BaseURL: %s", customURL)
	} else {
		c.logger.Infof("🔧 [MCP] OpenRouter using default BaseURL: %s", c.BaseURL)
	}
	if customModel != "" {
		c.Model = customModel
		c.logger.Infof("🔧 [MCP] OpenRouter using custom Model: %s", customModel)
	} else {
		c.logger.Infof("🔧 [MCP] OpenRouter using default Model: %s", c.Model)
	}
}

// OpenRouter uses standard Bearer auth (OpenAI-compatible)
func (c *OpenRouterClient) setAuthHeader(reqHeaders http.Header) {
	c.Client.setAuthHeader(reqHeaders)
	// OpenRouter recommends adding additional headers for better rate limits and tracking
	// These are optional but recommended
	reqHeaders.Set("HTTP-Referer", "https://nofx.ai") // Optional: for rankings
	reqHeaders.Set("X-Title", "NOFX AI Trading")      // Optional: for rankings
}
