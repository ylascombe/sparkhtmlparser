package main

import (
	"htmlparser/analyser"
	"github.com/gin-gonic/gin"
	"os"
	"net/http"
	"fmt"
)

func main() {


	// test connectity to spark
	url := os.Getenv("SPARK_DASHBOARD_URL")
	_, err := http.Get(url)

	if err != nil {
		fmt.Println("FATAL : Cannot request spark dashboard :-(, bye bye")
		os.Exit(-1)
	}

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



