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

func main() {
    // Load configuration from .env
    cfg := configs.LoadConfig()

    // Connect to MongoDB
    client, err := database.Connect(cfg.MongoURI)
    if err != nil {
        log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
    }
    log.Println("✅ Successfully connected to MongoDB")

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
        AllowOrigins: []string{
            "http://localhost:5173",                    // dev frontend
            "https://razorbill-website.vercel.app",    // prod frontend
            "https://muthomivictor.vercel.app",        // prod frontend mirror
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
        // Custom origin function for Capacitor apps
        AllowOriginFunc: func(origin string) bool {
            if origin == "" {
                return true // allow native apps without origin header
            }
            // allow capacitor://localhost specifically
            if origin == "capacitor://localhost" {
                return true
            }
            return false
        },
    }))


    // Register main API routes (Authors, Blogs, Comments, Shares)
    api.RegisterRoutes(r, client)

    // Swagger UI route
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Start server
    log.Printf("🚀 Server running on port %s", cfg.Port)
    if err := r.Run(":" + cfg.Port); err != nil {
        log.Fatalf("❌ Failed to start server: %v", err)
    }
}

