package main

import (
	"github.com/gin-gonic/gin"
)

// Server that simulate a running spark streaming 1.5.1 dashboard
func main() {

	router := gin.Default()
	router.StaticFile("/streaming", "worker9/streaming.html")

	// Listen and serve on 0.0.0.0:4050
	router.Run(":4050")

}
