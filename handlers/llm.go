package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"llm-mock-server/engine"
	"llm-mock-server/models"
	"llm-mock-server/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LLMHandler struct {
	store *store.Store
}

func NewLLMHandler(s *store.Store) *LLMHandler {
	return &LLMHandler{store: s}
}

func (h *LLMHandler) HandleChatCompletion(c *gin.Context) {
	var req models.ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rules, _ := h.store.GetEnabledRules()
	eng := engine.NewEngine(rules)

	responseStr, found := eng.MatchChat(req)

	var finalResponse models.ChatCompletionResponse
	if found {
		json.Unmarshal([]byte(responseStr), &finalResponse)
	} else {
		// Default fallback
		finalResponse = models.ChatCompletionResponse{
			ID:      "mock-" + uuid.New().String(),
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []struct {
				Index   int `json:"index"`
				Message struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"message"`
				FinishReason string `json:"finish_reason"`
			}{
				{
					Index: 0,
					Message: struct {
						Role    string `json:"role"`
						Content string `json:"content"`
					}{
						Role:    "assistant",
						Content: "This is a default mock response. You can configure rules via the admin API.",
					},
					FinishReason: "stop",
				},
			},
		}
	}

	c.JSON(http.StatusOK, finalResponse)
}

func (h *LLMHandler) HandleImageGeneration(c *gin.Context) {
	var req models.ImageGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rules, _ := h.store.GetEnabledRules()
	eng := engine.NewEngine(rules)

	responseStr, found := eng.MatchImage(req)

	var finalResponse models.ImageGenerationResponse
	if found {
		json.Unmarshal([]byte(responseStr), &finalResponse)
	} else {
		finalResponse = models.ImageGenerationResponse{
			Created: time.Now().Unix(),
			Data: []struct {
				URL string `json:"url"`
			}{
				{URL: "https://via.placeholder.com/1024x1024.png?text=Mock+Image"},
			},
		}
	}

	c.JSON(http.StatusOK, finalResponse)
}

func (h *LLMHandler) HandleVideoGeneration(c *gin.Context) {
	var req models.VideoGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rules, _ := h.store.GetEnabledRules()
	eng := engine.NewEngine(rules)

	// Reuse image matching logic for prompt comparison
	responseStr, found := eng.MatchImage(models.ImageGenerationRequest{Prompt: req.Prompt})

	var finalResponse models.VideoGenerationResponse
	if found {
		json.Unmarshal([]byte(responseStr), &finalResponse)
	} else {
		finalResponse = models.VideoGenerationResponse{
			Created: time.Now().Unix(),
			Data: []struct {
				URL string `json:"url"`
			}{
				{URL: "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_1mb.mp4"},
			},
		}
	}

	c.JSON(http.StatusOK, finalResponse)
}

func (h *LLMHandler) HandleVertexGenerateContent(c *gin.Context) {
	modelWithMethod := c.Param("modelWithMethod")
	if !strings.HasSuffix(modelWithMethod, ":generateContent") {
		c.JSON(http.StatusNotFound, gin.H{"error": "Method not found"})
		return
	}

	var req models.VertexGenerateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Adapt Vertex request to engine matching
	// Extract content for matching
	var fullText strings.Builder
	for _, content := range req.Contents {
		for _, part := range content.Parts {
			fullText.WriteString(part.Text)
			fullText.WriteString(" ")
		}
	}

	// Create a dummy chat request for matching
	matchReq := models.ChatCompletionRequest{
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "user", Content: fullText.String()},
		},
	}

	rules, _ := h.store.GetEnabledRules()
	eng := engine.NewEngine(rules)

	responseStr, found := eng.MatchChat(matchReq)

	var finalResponse models.VertexGenerateContentResponse
	if found {
		// Try to parse the rule response as Vertex format
		err := json.Unmarshal([]byte(responseStr), &finalResponse)
		if err != nil {
			// If not Vertex format, maybe it's OpenAI format?
			// For simplicity, we just wrap the text content if it's not a full JSON
			var openaiResp models.ChatCompletionResponse
			if err := json.Unmarshal([]byte(responseStr), &openaiResp); err == nil && len(openaiResp.Choices) > 0 {
				finalResponse = h.wrapTextInVertex(openaiResp.Choices[0].Message.Content)
			} else {
				finalResponse = h.wrapTextInVertex(responseStr)
			}
		}
	} else {
		finalResponse = h.wrapTextInVertex("This is a default Vertex AI mock response.")
	}

	c.JSON(http.StatusOK, finalResponse)
}

func (h *LLMHandler) wrapTextInVertex(text string) models.VertexGenerateContentResponse {
	resp := models.VertexGenerateContentResponse{}
	resp.Candidates = []struct {
		Content struct {
			Role  string `json:"role"`
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	}{
		{
			Content: struct {
				Role  string `json:"role"`
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			}{
				Role: "model",
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: text},
				},
			},
			FinishReason: "STOP",
		},
	}
	resp.UsageMetadata.PromptTokenCount = 10
	resp.UsageMetadata.CandidatesTokenCount = 10
	resp.UsageMetadata.TotalTokenCount = 20
	return resp
}
