package main

import (
	"log"
	"myapp/config"
	"myapp/route"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	err = config.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	route.RegisterRouter(app)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "9090"
	}

	err = app.Run(":" + port)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}

}
