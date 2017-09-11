package main

import (
	"htmlparser/analyser"
	"github.com/gin-gonic/gin"
)

func main() {


	router := gin.Default()

	users := router.Group("/")
	{
		users.GET("/", analyser.ParseSparkDashboard)
	}

	router.Run()


}



