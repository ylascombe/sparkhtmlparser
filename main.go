package main

import (
	"htmlparser/analyser"
	"github.com/gin-gonic/gin"
)

func main() {


	router := gin.Default()

	prometheus := router.Group("/")
	{
		prometheus.GET("/", analyser.Prometheus)
	}

	csv := router.Group("/csv")
	{
		csv.GET("/", analyser.Csv)
	}
	router.Run()


}



