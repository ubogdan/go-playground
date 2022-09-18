package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	// gin.SetMode(gin.ReleaseMode)

	// with recovery handler
	r.Use(gin.Recovery())

	// different ways of serving static files
	r.Static("/static", "./static")
	r.StaticFS("/assets", http.Dir("my_file_system"))
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// Set templates folder
	r.LoadHTMLGlob("templates/*.tmpl")

	// home page using html template
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", struct {
			Title string
		}{
			Title: "Home page",
		})
	})

	// Add sub router /api
	api := r.Group("/api")
	api.Use(func(c *gin.Context) {
		// Request logger middleware example
		log.Printf("got request method:%s path:%s", c.Request.Method, c.Request.URL.Path)
	}, func(c *gin.Context) {
		// route protect middleware
		if c.Request.URL.Path == "/api/denied" {
			c.AbortWithStatusJSON(400, gin.H{"error": "Access denied !"})
		}
	})

	// url: http://localhost:8080/api/ping
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, struct {
			Message string `json:"message"`
		}{
			Message: "ping",
		})
	})

	// url: http://localhost:8080/api/denied
	api.GET("/denied", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to the restricted area."})
	})

	// Listen and serve on 0.0.0.0:8080
	err := http.ListenAndServe(":8080", r) // or r.Run(":8080")
	if err != nil {
		log.Printf("http.ListenAndServe %w", err)
	}
}
