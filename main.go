package main

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"htmlparser/httpclient"
	"htmlparser/handlers"
)

func main() {

	// test connectity to spark
	if ! httpclient.SparkDashboardIsRequestable() {
		fmt.Println("FATAL : Cannot request spark dashboard :-(")
		fmt.Println("\tAre the SPARK_DASHBOARD_URL, SPARK_LOGIN and SPARK_PASSWORD env var correctly set ?")
		fmt.Println("\tbye bye")
		os.Exit(-1)
	}

	router := gin.Default()

	prometheus := router.Group("/")
	{
		prometheus.GET("/", handlers.Prometheus)
		prometheus.GET("/metrics", handlers.Prometheus)
	}

	csv := router.Group("/csv")
	{
		csv.GET("/", handlers.Csv)
	}

	// By default it serves on :8080 unless a PORT environment variable was defined.
	router.Run()

}
