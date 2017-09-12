package main

import (
	"fmt"
	"htmlparser/analyser"
	"os"
	"github.com/gin-gonic/gin"
	"htmlparser/httpclient"
)

func main() {

	// test connectity to spark
	url := os.Getenv("SPARK_DASHBOARD_URL")
	login := os.Getenv("SPARK_LOGIN")
	pass := os.Getenv("SPARK_PASSWORD")

	content, err := httpclient.RequestPage(url, login, pass)

	if err != nil {
		fmt.Println("FATAL : Cannot request spark dashboard :-(")
		fmt.Println("\tAre the SPARK_DASHBOARD_URL, SPARK_LOGIN and SPARK_PASSWORD env var correctly set ?")
		fmt.Println("\tbye bye")
		os.Exit(-1)
	}

	router := gin.Default()

	router.Use(TransferParamMiddleware(content))

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

	// By default it serves on :8080 unless a PORT environment variable was defined.
	router.Run()

}

func TransferParamMiddleware(htmlContent string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("htmlContent", htmlContent)
		c.Next()
	}
}
