package main

import (
	"log"
	"net/http"
	"time"

	"chatvui/backend/handlers"
	"chatvui/backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new chat hub
	hub := models.NewHub()
	go hub.Run()

	// Create a new handler
	handler := handlers.NewHandler(hub)

	// Set up gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API routes
	api := r.Group("/api")
	{
		// User registration
		api.POST("/users", handler.RegisterUser)

		// Get all users
		api.GET("/users", handler.GetUsers)

		// Get active users
		api.GET("/users/active", handler.GetActiveUsers)

		// WebSocket connection
		api.GET("/ws/:userId", handler.HandleWebSocket)
	}

	// Serve static files from the ui/dist directory
	r.Static("/", "../ui/dist")

	// Handle non-API routes by serving the index.html file
	r.NoRoute(func(c *gin.Context) {
		c.File("../ui/dist/index.html")
	})

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
