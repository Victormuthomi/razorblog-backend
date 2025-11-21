package main

import (
    "log"
    "razorblog-backend/api"
    "razorblog-backend/configs"
    "razorblog-backend/internal/database"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
    _ "razorblog-backend/swagger"
)

// @title RazorBlog API
// @version 1.0
// @description Backend API for RazorBlog including Authors, Blogs, Comments, Shares
// @termsOfService http://swagger.io/terms/

// @contact.name RazorBlog Support
// @contact.url http://razorblog.io
// @contact.email support@razorblog.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// main is the entry point for the RazorBlog backend
func main() {
    // Load configuration from .env
    cfg := configs.LoadConfig()

    // Connect to MongoDB
    client, err := database.Connect(cfg.MongoURI)
    if err != nil {
        log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
    }
    log.Println("✅ Successfully connected to MongoDB")

    // Ensure MongoDB disconnects on exit
    defer func() {
        if err := client.Disconnect(database.Ctx); err != nil {
            log.Printf("⚠️ Error disconnecting MongoDB: %v", err)
        } else {
            log.Println("MongoDB connection closed")
        }
    }()

    // Initialize Gin router
    r := gin.Default()

    // ⚡ CORS middleware
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "https://razorbill-website.vercel.app", "https://muthomivictor.vercel.app/"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Register application routes (Authors, Blogs, Comments, Shares)
    api.RegisterRoutes(r, client)

    // Swagger UI route
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Start HTTP server
    log.Printf("🚀 Server running on port %s", cfg.Port)
    if err := r.Run(":" + cfg.Port); err != nil {
        log.Fatalf("❌ Failed to start server: %v", err)
    }
}

