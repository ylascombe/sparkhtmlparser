package main

import (
	"fmt"
	"htmlparser/analyser"
	"htmlparser/httpclient"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	// test connectity to spark
	url := os.Getenv("SPARK_DASHBOARD_URL")

	ok := httpclient.IsRequestable(url, "spark", "spark")

	if !ok {
		fmt.Println("FATAL : Cannot request spark dashboard :-(")
		fmt.Println("\tIs the SPARK_DASHBOARD_URL env var is correctly set ?")
		fmt.Println("\tbye bye")
		os.Exit(-1)
	}

	//fmt.Println("tHE RESULT", content)

	router := gin.Default()

	prometheus := router.Group("/")
	{
		prometheus.GET("/", analyser.Prometheus)
	}

	metrics := router.Group("/metrics")
	{
		metrics.GET("/", analyser.Prometheus)
	}

	csv := router.Group("/csv")
	{
		csv.GET("/", analyser.Csv)
	}
	router.Run(":5000")

}
