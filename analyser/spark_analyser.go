package analyser

import (
	"bytes"
	"errors"
	"fmt"
	"htmlparser/httpclient"
	"htmlparser/models"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/montanaflynn/stats"
	"golang.org/x/net/html"
)

func Prometheus(c *gin.Context) {

	url := os.Getenv("SPARK_DASHBOARD_URL")
	content, err := httpclient.RequestPage(url, "spark", "spark")

	if err != nil {
		c.String(503, "error when trying to request spark")
	}
	res, err := parseSparkDashboard(content)

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

	printPercentiles(c, evtsPerBatchs, "events_per_second")
	printPercentiles(c, schedulingDelays, "scheduling_delay")
	printPercentiles(c, processingTimes, "processing_time")
}

func Csv(c *gin.Context) {

	url := os.Getenv("SPARK_DASHBOARD_URL")
	content, err := httpclient.RequestPage(url, "spark", "spark")

	if err != nil {
		c.String(503, "error when trying to request spark")
	}

	res, err := parseSparkDashboard(content)

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

func parseSparkDashboard(content string) (*models.Report, error) {

	doc, _ := html.Parse(strings.NewReader(content))
	bn, err := GetTableBody(doc)
	if err != nil {
		return &models.Report{}, err
	}

	tbody, err := FindFirstChild(bn, "tbody")

	if err != nil {
		return &models.Report{}, err
	}

	res, err := browseTr(tbody)

	if err == nil {
		return &res, nil
	} else {
		return &models.Report{}, err
	}
}

func printPercentiles(c *gin.Context, datas []float64, label string) {
	avg, _ := stats.Mean(datas)
	c.String(200, fmt.Sprintf("\n%s_avg %v", label, avg))
	min, _ := stats.Min(datas)
	c.String(200, fmt.Sprintf("\n%s_min %v", label, min))
	c.String(200, fmt.Sprintf("\n%s_10_percentile %v", label, getPercentile(datas, 10)))
	c.String(200, fmt.Sprintf("\n%s_20_percentile %v", label, getPercentile(datas, 20)))
	c.String(200, fmt.Sprintf("\n%s_25_percentile %v", label, getPercentile(datas, 25)))
	c.String(200, fmt.Sprintf("\n%s_30_percentile %v", label, getPercentile(datas, 30)))
	c.String(200, fmt.Sprintf("\n%s_40_percentile %v", label, getPercentile(datas, 40)))
	c.String(200, fmt.Sprintf("\n%s_50_percentile %v", label, getPercentile(datas, 50)))
	c.String(200, fmt.Sprintf("\n%s_60_percentile %v", label, getPercentile(datas, 60)))
	c.String(200, fmt.Sprintf("\n%s_70_percentile %v", label, getPercentile(datas, 70)))
	c.String(200, fmt.Sprintf("\n%s_80_percentile %v", label, getPercentile(datas, 80)))
	c.String(200, fmt.Sprintf("\n%s_90_percentile %v", label, getPercentile(datas, 90)))
	c.String(200, fmt.Sprintf("\n%s_95_percentile %v", label, getPercentile(datas, 95)))
	max, _ := stats.Max(datas)
	c.String(200, fmt.Sprintf("\n%s_max %v", label, max))

}

func GetTableBody(doc *html.Node) (*html.Node, error) {
	var res *html.Node

	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "table" {

			for i := 0; i < len(node.Attr); i++ {
				if node.Attr[i].Key == "id" && node.Attr[i].Val == "completed-batches-table" {
					res = node
				}
			}

		}

		if res == nil {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				f(child)
			}
		}
	}
	f(doc)
	if res != nil {
		return res, nil
	}
	return nil, errors.New("Missing <table> with id completed-batches-table in the node tree")
}

func FindFirstChild(doc *html.Node, tagName string) (*html.Node, error) {
	var res *html.Node
	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tagName {
			res = node
		}

		if res == nil {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				f(child)

				if res != nil {
					break
				}
			}
		}
	}
	f(doc)
	if res != nil {
		return res, nil
	}
	return nil, errors.New("Missing <" + tagName + "> in the node tree")
}

func browseTr(tr *html.Node) (models.Report, error) {

	lignes := -1
	var batches []models.Batch

	evtsNumber := 0
	processings := float32(0.0)
	schedulingDelays := float32(0.0)

	for child := tr.FirstChild; child != nil; child = child.NextSibling {

		if child.Data == "tr" {

			td, err := FindFirstChild(child, "td")

			if err != nil {
				return models.Report{}, errors.New("Cannot find td")
			}

			batch, err := browseTd(td)

			if err != nil {
				return models.Report{}, nil
			}
			batches = append(batches, batch)

			processings += batch.ProcessingTime
			evtsNumber += batch.InputSize
			schedulingDelays += float32(batch.SchedulingDelay)
			lignes++
			batch = models.Batch{}
		}
	}

	report := models.Report{
		Batches:            batches,
		EventsPerSecondAvg: int(float32(evtsNumber) / processings),
		RowCount:           lignes + 1,
	}
	return report, nil
}

func browseTd(td *html.Node) (models.Batch, error) {
	cols := 0
	batch := models.Batch{}

	for child := td; child != nil; child = child.NextSibling {

		if child.Data == "td" {

			val := strings.Trim(renderNode(child.FirstChild), " ")
			val = strings.Replace(val, " ", "", -1)
			val = strings.Replace(val, "\n", "", -1)
			val = strings.Replace(val, "\r", "", -1)

			switch cols {
			case 0:
				val := strings.Trim(renderNode(child.FirstChild.NextSibling.FirstChild), " ")
				val = strings.Replace(val, "\r", "", -1)
				val = strings.Replace(val, "\n", "", -1)
				val = strings.Replace(val, "  ", "", -1)
				val = strings.Trim(val, "")
				val = strings.Replace(val, " ", "_", -1)

				batch.BatchTime = val
			case 1:
				val = strings.Replace(val, "events", "", 1)
				val, _ := strconv.Atoi(val)

				batch.InputSize = val
			case 2:
				if strings.Index(val, "ms") > 0 {
					val = strings.Replace(val, "ms", "", 1)
					val, _ := strconv.Atoi(val)
					batch.SchedulingDelay = val
				} else {
					val = strings.Replace(val, "s", "", 1)
					val, _ := strconv.Atoi(val)
					batch.SchedulingDelay = (val * 1000)
				}
			case 3:
				val = strings.Replace(val, "s", "", 1)
				val, _ := strconv.ParseFloat(val, 4)

				batch.ProcessingTime = float32(val)
			case 4:
				val = strings.Replace(val, "s", "", 1)
				val, _ := strconv.ParseFloat(val, 4)

				batch.TotalDelay = float32(val)
			}
			cols++
		}
	}
	return batch, nil
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func getPercentile(array []float64, percent int) float64 {
	percentile, _ := stats.Percentile(array, float64(percent))
	return percentile
}
