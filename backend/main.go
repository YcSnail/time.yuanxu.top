package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/YcSnail/time.yuanxu.top/backend/config"
	"github.com/YcSnail/time.yuanxu.top/backend/database"
	"github.com/YcSnail/time.yuanxu.top/backend/handlers"
	"github.com/YcSnail/time.yuanxu.top/backend/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	authH := &handlers.AuthHandler{DB: db, Secret: cfg.JWTSecret, Expire: cfg.JWTExpire}
	cdH := &handlers.CountdownHandler{DB: db}

	r := gin.Default()
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "time": cfg.JWTExpire})
		})

		api.POST("/enter", authH.Enter)

		authed := api.Group("", middleware.Auth(cfg.JWTSecret))
		{
			authed.GET("/me", authH.Me)
			authed.GET("/countdowns", cdH.List)
			authed.POST("/countdowns", cdH.Create)
			authed.DELETE("/countdowns/:id", cdH.Delete)
		}
	}

	log.Printf("time.yuanxu.top backend listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}

// loadEnv loads .env from the current directory, falling back to the
// executable's directory (useful when the binary is run from elsewhere).
func loadEnv() {
	candidates := []string{}
	if dir, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(dir, ".env"))
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), ".env"))
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		if err := godotenv.Load(path); err != nil {
			log.Printf("warning: failed to load %s: %v", path, err)
		} else {
			log.Printf("loaded env from %s", path)
			return
		}
	}
	log.Println("no .env file found, falling back to environment variables")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if origin := c.GetHeader("Origin"); origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Max-Age", "86400")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
