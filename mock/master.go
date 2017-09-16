package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Server that simulate a running spark streaming 1.5.1 dashboard
func main() {

	router := gin.Default()
	router.StaticFS("/", http.Dir("mainPage"))
	// router.Static("/worker9/streaming/", "worker9/streaming.html")
	//router.StaticFS("/", http.Dir("."))

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8088")

}
