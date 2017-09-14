package analyser

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"strings"
	"golang.org/x/net/html"

	"fmt"
	"io/ioutil"
	"os"
)

func TestFindTagWithId(t *testing.T) {

	doc, _ := html.Parse(strings.NewReader(HTML1))
	res, err := FindTagWithId(doc, "table", "completed-batches-table")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	fmt.Println(res)
	assert.Equal(t, "id", res.Attr[0].Key)
	assert.Equal(t, "completed-batches-table", res.Attr[0].Val)
}

func TestFindFirstChild(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(HTML1))
	res, _ := FindTagWithId(doc, "table", "completed-batches-table")

	res, err := FindFirstChild(res, "td")

	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, "id", res.Attr[0].Key)
	assert.Equal(t, "batch-1504876700000", res.Attr[0].Val)
}

func TestFindTagWithContent(t *testing.T) {

	pageContent, _ := readFile("/mock/mainPage/mainpage.html")
	doc, _ := html.Parse(strings.NewReader(pageContent))

	node, err := FindTagWithContent(doc, "h4", "<h4> Running Applications </h4>")

	assert.Nil(t, err)
	assert.NotNil(t, node)

	assert.Equal(t, " Running Applications ", node.FirstChild.Data)
}


func TestRenderNode(t *testing.T) {

	doc, _ := html.Parse(strings.NewReader(HTML1))
	td, _ := FindTagWithId(doc, "td", "batch-1504876700000")

	res := renderNode(td)

	assert.NotNil(t, res)

	expected := "<td id=\"batch-1504876700000\" sorttable_customkey=\"1504876700000\">\n<a href=\"http://cops-fco-spark-worker-a-11.cloud.alt:4041/streaming/batch?id=1504876700000\">\n2017/09/08 15:18:20\n</a>\n</td>"
	assert.Equal(t, expected, res)
}

func TestFindWorkerForApp(t *testing.T) {

	pageContent, _ := readFile("/mock/mainPage/mainpage.html")

	res, err := FindWorkerForApp("colis360", pageContent)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}

// TODO test following functions
// ParseSparkDashboard
// browseTr
// browseTd

func TestReadFile(t *testing.T) {
	pageContent, err := readFile("/mock/example.test")

	expected := "Example file to test readContent function\n\nThere are also new lines !\n"
	assert.Nil(t, err)
	assert.NotNil(t, pageContent)
	assert.Equal(t, expected, pageContent)

}
func readFile(pathInProject string) (string, error) {
	pwd, _ := os.Getwd()
	pwd = strings.Replace(pwd, "/analyser", "", 1)

	path := pwd + pathInProject
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}
	return string(data), nil
}

const HTML1 = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-type" content="text/html; charset=UTF-8">
<link rel="stylesheet" href="colis360%20-%20Streaming%20Statistics_fichiers/bootstrap.css" type="text/css">
<link rel="stylesheet" href="colis360%20-%20Streaming%20Statistics_fichiers/vis.css" type="text/css">
<link rel="stylesheet" href="colis360%20-%20Streaming%20Statistics_fichiers/webui.css" type="text/css">
<link rel="stylesheet" href="colis360%20-%20Streaming%20Statistics_fichiers/timeline-view.css" type="text/css">
<script src="colis360%20-%20Streaming%20Statistics_fichiers/sorttable.js"></script>
<script src="colis360%20-%20Streaming%20Statistics_fichiers/jquery-1.js"></script>
<script src="colis360%20-%20Streaming%20Statistics_fichiers/vis.js"></script>
<script src="colis360%20-%20Streaming%20Statistics_fichiers/bootstrap-tooltip.js"></script>
<script src="colis360%20-%20Streaming%20Statistics_fichiers/initialize-tooltips.js"></script>
<script src="colis360%20-%20Streaming%20Statistics_fichiers/table.js"></script>
<script src="colis360%20-%20Streaming%20Statistics_fichiers/additional-metrics.js"></script>
<script src="colis360%20-%20Streaming%20Statistics_fichiers/timeline-view.js"></script>

<title>colis360 - Streaming Statistics</title>
</head>
<body>
<div class="navbar navbar-static-top">
<div class="navbar-inner">
<div class="brand">
<a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/" class="brand">
<img src="colis360%20-%20Streaming%20Statistics_fichiers/spark-logo-77x50px-hd.png">
<span class="version">1.5.1</span>
</a>
</div>
<ul class="nav">
<li class="">
<a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/jobs/">Jobs</a>
</li>
<li class="">
<a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/stages/">Stages</a>
</li>
<li class="">
<a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/storage/">Storage</a>
</li>
<li class="">
<a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/environment/">Environment</a>
</li>
<li class="">
<a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/executors/">Executors</a>
</li>
<li class="active">
<a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/streaming/">Streaming</a>
</li>
</ul>
<p class="navbar-text pull-right">
<strong title="colis360">colis360</strong> application UI
</p>
</div>
</div>
<div class="container-fluid">
<div class="row-fluid">
<div class="span12">
<h3 style="vertical-align: bottom; display: inline-block;">
Streaming Statistics

</h3>
</div>
</div>
<script src="colis360%20-%20Streaming%20Statistics_fichiers/d3.js"></script>
<link rel="stylesheet" href="colis360%20-%20Streaming%20Statistics_fichiers/streaming-page.css" type="text/css">
<script src="colis360%20-%20Streaming%20Statistics_fichiers/streaming-page.js"></script>
<div>Running batches of
<strong>
10 seconds
</strong>
for
<strong>
2 days 21 hours 24 minutes
</strong>
since
<strong>
2017/09/05 17:54:18
</strong>
(<strong>24978</strong>
completed batches, <strong>21278433</strong> records)
</div>
<br>

<table id="completed-batches-table" class="table table-bordered table-striped table-condensed sortable">
<thead>
<tr>
<th>Batch Time</th>
<th>Input Size</th>
<th>Scheduling Delay
<sup>
(<a data-toggle="tooltip" data-placement="top" title=""
data-original-title="Time taken by Streaming scheduler to submit jobs of a batch">?</a>)
</sup>
</th>
<th>Processing Time
<sup>
(<a data-toggle="tooltip" data-placement="top" title=""
data-original-title="Time taken to process all jobs of a batch">?</a>)
</sup></th>
<th>Total Delay
<sup>
(<a data-toggle="tooltip" data-placement="top" title=""
data-original-title="Total time taken to handle a batch">?</a>)
</sup></th>
</tr>
</thead>
<tbody>
<tr>
<td id="batch-1504876700000" sorttable_customkey="1504876700000">
<a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/streaming/batch?id=1504876700000">
2017/09/08 15:18:20
</a>
</td>
<td sorttable_customkey="1112">1112 events</td>
<td sorttable_customkey="3">
3 ms
</td>
<td sorttable_customkey="3670">
4 s
</td>
<td sorttable_customkey="3674">
5 s
</td>
</tr>
</tbody>
<tfoot></tfoot>
</table>


</div>
</body>
</html>
`
