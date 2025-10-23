package lightrag

type Message struct {
	Role    string   `json:"role,omitempty"`
	Content string   `json:"content,omitempty"`
	Images  []string `json:"images,omitempty"`
}

type ChatRequest struct {
	Mode                string    `json:"mode"`
	ResponseType        string    `json:"response_type"`
	TopK                int       `json:"top_k"`
	ChunkTopK           int       `json:"chunk_top_k"`
	MaxEntityTokens     int       `json:"max_entity_tokens"`
	MaxRelationTokens   int       `json:"max_relation_tokens"`
	MaxTotalTokens      int       `json:"max_total_tokens"`
	OnlyNeedContext     bool      `json:"only_need_context"`
	OnlyNeedPrompt      bool      `json:"only_need_prompt"`
	Stream              bool      `json:"stream"`
	HistoryTurns        int       `json:"history_turns"`
	UserPrompt          string    `json:"user_prompt"`
	EnableRerank        bool      `json:"enable_rerank"`
	Query               string    `json:"query"`
	ConversationHistory []Message `json:"conversation_history"`
}

type ChatResponse struct {
	Model           string  `json:"model,omitempty"`
	CreatedAt       string  `json:"created_at,omitempty"`
	Message         Message `json:"message,omitempty"`
	Response        string  `json:"response,omitempty"` // for stream response
	Done            bool    `json:"done,omitempty"`
	TotalDuration   int     `json:"total_duration,omitempty"`
	LoadDuration    int     `json:"load_duration,omitempty"`
	PromptEvalCount int     `json:"prompt_eval_count,omitempty"`
	EvalCount       int     `json:"eval_count,omitempty"`
	EvalDuration    int     `json:"eval_duration,omitempty"`
	Error           string  `json:"error,omitempty"`

	Answer  string   `json:"answer"`
	Sources []string `json:"sources,omitempty"` // 来源文档（可选）
}

type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
	// Truncate  bool     `json:"truncate,omitempty"`
	// KeepAlive string   `json:"keep_alive,omitempty"`
}

type EmbeddingResponse struct {
	Error      string      `json:"error,omitempty"`
	Model      string      `json:"model"`
	Embeddings [][]float64 `json:"embeddings"`
}

// LightRAG 响应结构
type LightRAGResponse struct {
	Answer   string   `json:"answer"`
	Sources  []string `json:"sources,omitempty"` // 来源文档（可选）
	Response string   `json:"response"`
	Error    string   `json:"error"`
}
