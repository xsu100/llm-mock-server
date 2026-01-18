package main

import (
	"log"

	"llm-mock-server/handlers"
	"llm-mock-server/store"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Store
	s, err := store.NewStore("./mock.db")
	if err != nil {
		log.Fatal("Failed to initialize store:", err)
	}

	if err := s.SeedDefaults(); err != nil {
		log.Println("Warning: Failed to seed default rules:", err)
	}

	// Initialize Handlers
	llm := handlers.NewLLMHandler(s)
	admin := handlers.NewAdminHandler(s)

	r := gin.Default()

	// OpenAI Compatible Routes
	v1 := r.Group("/v1")
	{
		v1.POST("/chat/completions", llm.HandleChatCompletion)
		v1.POST("/images/generations", llm.HandleImageGeneration)
		v1.POST("/videos/generations", llm.HandleVideoGeneration)
	}

	// Vertex AI Compatible Routes
	// v1/projects/:project/locations/:location/publishers/google/models/:model:generateContent
	r.POST("/v1/projects/:project/locations/:location/publishers/google/models/*modelWithMethod", llm.HandleVertexGenerateContent)
	r.POST("/v1beta1/projects/:project/locations/:location/publishers/google/models/*modelWithMethod", llm.HandleVertexGenerateContent)

	// Admin API
	mgmt := r.Group("/admin")
	{
		mgmt.GET("/rules", admin.ListRules)
		mgmt.POST("/rules", admin.CreateRule)
		mgmt.DELETE("/rules/:id", admin.DeleteRule)
	}

	log.Println("LLM Mock Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
