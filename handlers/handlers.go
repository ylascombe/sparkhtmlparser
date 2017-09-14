package handlers

import (
	"github.com/montanaflynn/stats"
	"fmt"
	"github.com/gin-gonic/gin"
	"htmlparser/analyser"
	"htmlparser/httpclient"
	"errors"
	"os"
)

func Prometheus(c *gin.Context) {

	content, err := httpclient.RequestSparkDashboard()

	if err != nil {
		c.Error(errors.New("Cannot request spark dashboard"))
	}
	url, err := analyser.FindWorkerLinkForApp("colis360", content)

	if err != nil {
		c.String(503, "Cannot find app url in spark main dashboard")
	}

	// Now, request spark worker dashboard
	login := os.Getenv("SPARK_LOGIN")
	pass := os.Getenv("SPARK_PASSWORD")
	content, err = httpclient.RequestPage(url, login, pass)

	if err != nil {
		c.Error(errors.New("Cannot request spark dashboard"))
	}
	res, err := analyser.ParseSparkDashboard(content)

	if err != nil {
		c.String(503, "Error ")
	}
	var schedulingDelays []float64
	var evtsPerBatchs []float64
	var processingTimes []float64

	for i := 0; i < len(res.Batches); i++ {
		schedulingDelays = append(schedulingDelays, float64(res.Batches[i].SchedulingDelay))
		evtsPerBatchs = append(evtsPerBatchs, float64(res.Batches[i].InputSize))
		processingTimes = append(processingTimes, float64(res.Batches[i].ProcessingTime))
	}

	// actual_events_per_second_processed
	sumEvts, _ := stats.Sum(evtsPerBatchs)
	sumProcessing, _ := stats.Sum(processingTimes)
	c.String(200, fmt.Sprintf("\nactual_events_per_second_processed %v", sumEvts/sumProcessing))

	// actual_events_median_processed
	medianEvts, _ := stats.Median(evtsPerBatchs)
	medianProcessing, _ := stats.Median(processingTimes)
	c.String(200, fmt.Sprintf("\nactual_events_median_processed %v", medianEvts/medianProcessing))

	printPercentiles(c, evtsPerBatchs, "events_per_second", 10)
	printPercentiles(c, schedulingDelays, "scheduling_delay", 1)
	printPercentiles(c, processingTimes, "processing_time", 1)

	c.String(200, "\n\t")
}



func Csv(c *gin.Context) {

	content, err := httpclient.RequestSparkDashboard()

	if err != nil {
		c.Error(errors.New("Cannot request spark dashboard"))
	}
	url, err := analyser.FindWorkerLinkForApp("colis360", content)

	if err != nil {
		c.String(503, "Cannot find app url in spark main dashboard")
	}

	// Now, request spark worker dashboard
	login := os.Getenv("SPARK_LOGIN")
	pass := os.Getenv("SPARK_PASSWORD")
	content, err = httpclient.RequestPage(url, login, pass)

	if err != nil {
		c.Error(errors.New("Cannot request spark dashboard"))
	}
	res, err := analyser.ParseSparkDashboard(content)

	if err != nil {
		c.String(503, "Error")
	}

	c.String(200, "Batch Time,Input Size,Scheduling Delay,Processing Time,Total Delay")

	for i := 0; i < len(res.Batches); i++ {
		c.String(200, fmt.Sprintf("\n%s,%d,%d,%v,%v",
			res.Batches[i].BatchTime,
			res.Batches[i].InputSize,
			res.Batches[i].SchedulingDelay,
			res.Batches[i].ProcessingTime,
			res.Batches[i].TotalDelay,
		))
	}

}

func printPercentiles(c *gin.Context, datas []float64, label string, coeff int) {
	avg, _ := stats.Mean(datas)
	avg = avg / float64(coeff)
	c.String(200, fmt.Sprintf("\n%s_avg %v", label, avg))
	min, _ := stats.Min(datas)
	min = min / float64(coeff)
	c.String(200, fmt.Sprintf("\n%s_min %v", label, min))
	c.String(200, fmt.Sprintf("\n%s_10_percentile %v", label, getPercentile(datas, 10, coeff)))
	c.String(200, fmt.Sprintf("\n%s_20_percentile %v", label, getPercentile(datas, 20, coeff)))
	c.String(200, fmt.Sprintf("\n%s_25_percentile %v", label, getPercentile(datas, 25, coeff)))
	c.String(200, fmt.Sprintf("\n%s_30_percentile %v", label, getPercentile(datas, 30, coeff)))
	c.String(200, fmt.Sprintf("\n%s_40_percentile %v", label, getPercentile(datas, 40, coeff)))
	c.String(200, fmt.Sprintf("\n%s_50_percentile %v", label, getPercentile(datas, 50, coeff)))
	c.String(200, fmt.Sprintf("\n%s_60_percentile %v", label, getPercentile(datas, 60, coeff)))
	c.String(200, fmt.Sprintf("\n%s_70_percentile %v", label, getPercentile(datas, 70, coeff)))
	c.String(200, fmt.Sprintf("\n%s_80_percentile %v", label, getPercentile(datas, 80, coeff)))
	c.String(200, fmt.Sprintf("\n%s_90_percentile %v", label, getPercentile(datas, 90, coeff)))
	c.String(200, fmt.Sprintf("\n%s_95_percentile %v", label, getPercentile(datas, 95, coeff)))
	max, _ := stats.Max(datas)
	max = max / float64(coeff)
	c.String(200, fmt.Sprintf("\n%s_max %v", label, max))
}

func getPercentile(array []float64, percent int, coeff int) float64 {
	percentile, _ := stats.Percentile(array, float64(percent))
	return percentile / float64(coeff)
}

// TODO DEPRECATED to remove after validation
func FindPage(c *gin.Context) {
	content, err := httpclient.RequestSparkDashboard()

	if err != nil {
		c.Error(errors.New("Cannot request spark dashboard"))
	}
	url, err := analyser.FindWorkerLinkForApp("colis360", content)

	if err != nil {
		c.String(503, "Error")
	}
	fmt.Println(url)
	return
}
