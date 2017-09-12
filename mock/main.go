package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Server that simulate a running spark streaming 1.5.1 dashboard
func main() {

	router := gin.Default()
	router.StaticFS("/", http.Dir("myApp"))

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8088")

}