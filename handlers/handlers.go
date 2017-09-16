package handlers

import (
	"fmt"
	"htmlparser/analyser"
	"htmlparser/httpclient"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/montanaflynn/stats"
	"htmlparser/models"
)

func Prometheus(c *gin.Context) {

	appSparkMetrics, err := getMetrics()

	if err != nil {
		c.Error(err)
		return
	}
	var schedulingDelays []float64
	var evtsPerBatchs []float64
	var processingTimes []float64

	for i := 0; i < len(appSparkMetrics.Batches); i++ {
		schedulingDelays = append(schedulingDelays, float64(appSparkMetrics.Batches[i].SchedulingDelay))
		evtsPerBatchs = append(evtsPerBatchs, float64(appSparkMetrics.Batches[i].InputSize))
		processingTimes = append(processingTimes, float64(appSparkMetrics.Batches[i].ProcessingTime))
	}

	// actual_events_per_second_processed
	sumEvts, _ := stats.Sum(evtsPerBatchs)
	sumProcessing, _ := stats.Sum(processingTimes)

	str := ""
	str += fmt.Sprintf("\nactual_events_per_second_processed %v", sumEvts/sumProcessing)

	// actual_events_median_processed
	medianEvts, _ := stats.Median(evtsPerBatchs)
	medianProcessing, _ := stats.Median(processingTimes)
	str += fmt.Sprintf("\nactual_events_median_processed %v", medianEvts/medianProcessing)

	str += printPercentiles(evtsPerBatchs, "events_per_second", 10)
	str += printPercentiles(schedulingDelays, "scheduling_delay", 1)
	str += printPercentiles(processingTimes, "processing_time", 1)
	str += "\n\t"
	c.String(200, str)
}

func Csv(c *gin.Context) {

	res, err := getMetrics()

	if err != nil {
		c.Error(err)
		return
	}

	str := "Batch Time,Input Size,Scheduling Delay,Processing Time,Total Delay"

	for i := 0; i < len(res.Batches); i++ {
		str += fmt.Sprintf("\n%s,%d,%d,%v,%v",
			res.Batches[i].BatchTime,
			res.Batches[i].InputSize,
			res.Batches[i].SchedulingDelay,
			res.Batches[i].ProcessingTime,
			res.Batches[i].TotalDelay,
		)
	}

	c.String(200, str)
}

func printPercentiles(datas []float64, label string, coeff int) string {
	avg, _ := stats.Mean(datas)
	avg = avg / float64(coeff)
	str := fmt.Sprintf("\n%s_avg %v", label, avg)
	min, _ := stats.Min(datas)
	min = min / float64(coeff)
	str += fmt.Sprintf("\n%s_min %v", label, min)
	str += fmt.Sprintf("\n%s_10_percentile %v", label, getPercentile(datas, 10, coeff))
	str += fmt.Sprintf("\n%s_20_percentile %v", label, getPercentile(datas, 20, coeff))
	str += fmt.Sprintf("\n%s_25_percentile %v", label, getPercentile(datas, 25, coeff))
	str += fmt.Sprintf("\n%s_30_percentile %v", label, getPercentile(datas, 30, coeff))
	str += fmt.Sprintf("\n%s_40_percentile %v", label, getPercentile(datas, 40, coeff))
	str += fmt.Sprintf("\n%s_50_percentile %v", label, getPercentile(datas, 50, coeff))
	str += fmt.Sprintf("\n%s_60_percentile %v", label, getPercentile(datas, 60, coeff))
	str += fmt.Sprintf("\n%s_70_percentile %v", label, getPercentile(datas, 70, coeff))
	str += fmt.Sprintf("\n%s_80_percentile %v", label, getPercentile(datas, 80, coeff))
	str += fmt.Sprintf("\n%s_90_percentile %v", label, getPercentile(datas, 90, coeff))
	str += fmt.Sprintf("\n%s_95_percentile %v", label, getPercentile(datas, 95, coeff))
	max, _ := stats.Max(datas)
	max = max / float64(coeff)
	str +=  fmt.Sprintf("\n%s_max %v", label, max)

	return str
}

func getPercentile(array []float64, percent int, coeff int) float64 {
	percentile, _ := stats.Percentile(array, float64(percent))
	return percentile / float64(coeff)
}

func getMetrics() (*models.Report, error) {

	content, err := httpclient.RequestActiveSparkMasterContent()

	if err != nil {
		return nil, fmt.Errorf("Cannot request spark dashboard. Error detail : %s", err.Error())
	}

	appName := os.Getenv("SPARK_APP")
	url, err := analyser.FindWorkerLinkForApp(appName, *content)
	url = url + "/streaming/"

	fmt.Println("Link for app name:", url)

	if err != nil {
		return nil, fmt.Errorf("Cannot find app url in spark main dashboard. %s", err)
	}

	// Now, request spark worker dashboard
	login := os.Getenv("SPARK_LOGIN")
	pass := os.Getenv("SPARK_PASSWORD")
	content, err = httpclient.RequestPage(url, login, pass)

	if err != nil {
		return nil, fmt.Errorf("%s. %s", "Cannot request spark dashboard : ", err)
	}
	res, err := analyser.ParseSparkDashboard(*content)

	if err != nil {
		return nil, fmt.Errorf("Error while parsing spark worker streaming page. %s", err)
	}
	return res, nil
}
