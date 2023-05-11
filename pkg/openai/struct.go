package openai

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenaiChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type SseResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type OpenaiChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Delta        struct {
		Content string `json:"content"`
	} `json:"delta"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	Complationtokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
