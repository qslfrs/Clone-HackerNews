package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// init cache and client
	InitCache(5 * time.Minute)              // default TTL 5m
	client := NewHNClient(10 * time.Second) // http client timeout

	r := gin.Default()

	// Simple CORS for dev (adjust for prod)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")
		c.Header("Access-Control-Allow-Methods", "GET,OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/topstories", func(c *gin.Context) {
			GetTopStoriesHandler(c, client)
		})
		api.GET("/item/:id", func(c *gin.Context) {
			GetItemHandler(c, client)
		})
		api.GET("/user/:id", func(c *gin.Context) {
			GetUserHandler(c, client)
		})
	}

	addr := ":8080"
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
