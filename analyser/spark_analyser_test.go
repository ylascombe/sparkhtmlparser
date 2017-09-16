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

	// arrange
	doc, _ := html.Parse(strings.NewReader(HTML1))

	// act
	res, err := FindTagWithId(doc, "table", "completed-batches-table")

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
	fmt.Println(res)
	assert.Equal(t, "id", res.Attr[0].Key)
	assert.Equal(t, "completed-batches-table", res.Attr[0].Val)
}

func TestFindFirstChild(t *testing.T) {
	// arrange
	doc, _ := html.Parse(strings.NewReader(HTML1))
	res, _ := FindTagWithId(doc, "table", "completed-batches-table")

	// act
	res, err := FindFirstChild(res, "td")

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, "id", res.Attr[0].Key)
	assert.Equal(t, "batch-1504876700000", res.Attr[0].Val)
}

func TestFindTagWithContent(t *testing.T) {

	// arrange
	pageContent, _ := readFile("/mock/mainPage/sparkMasterActive.html")
	doc, _ := html.Parse(strings.NewReader(pageContent))

	// act
	node, err := FindTagWithContent(doc, "h4", "<h4> Running Applications </h4>")

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, node)

	assert.Equal(t, " Running Applications ", node.FirstChild.Data)
}


func TestRenderNode(t *testing.T) {

	// arrange
	doc, _ := html.Parse(strings.NewReader(HTML1))
	td, _ := FindTagWithId(doc, "td", "batch-1504876700000")

	// act
	res := renderNode(td)

	// assert
	assert.NotNil(t, res)

	expected := "<td id=\"batch-1504876700000\" sorttable_customkey=\"1504876700000\">\n<a href=\"http://cops-fco-spark-worker-a-11.cloud.alt:4041/streaming/batch?id=1504876700000\">\n2017/09/08 15:18:20\n</a>\n</td>"
	assert.Equal(t, expected, res)
}

func TestFindWorkerLinkForApp(t *testing.T) {

	// arrange
	pageContent, _ := readFile("/mock/mainPage/sparkMasterActive.html")

	// act
	res, err := FindWorkerLinkForApp("colis360", pageContent)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "http://localhost:8088/worker9/main.html", res)
}

func TestFindWorkerLinkForAppWhenNotFound(t *testing.T) {

	// arrange
	pageContent, _ := readFile("/mock/mainPage/sparkMasterActive.html")
	appName := "do_not_exist"

	// act
	res, err := FindWorkerLinkForApp(appName, pageContent)

	// assert
	assert.NotNil(t, err)
	assert.Equal(t, "", res)
	assert.Equal(t, "Link not found for application " + appName, err.Error())
}

func TestGenericTRBrowser(t *testing.T) {
	// arrange
	content, _ := readFile("/mock/mainPage/sparkMasterActive.html")
	doc, _ := html.Parse(strings.NewReader(content))
	node, _ := FindTagWithContent(doc, "h4", "<h4> Running Applications </h4>")

	table := node.NextSibling.NextSibling

	tbody, err := FindFirstChild(table, "tbody")

	// act
	res, err := genericTRBrowser(tbody)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, 3, len(*res))
	line := (*res)[0]
	assert.Equal(t, 8, len(line.Cells))
	// cells value will be tested on testGenericTDBrowser
}


func TestGenericTDBrowser(t *testing.T) {
	// arrange
	content := `<html><body><table><tr>
      <td>
        <a href="app?appId=app-20170912104309-5902">app-20170912104309-5902</a>
        <form action="app/kill/" method="POST" style="display:inline">
        <input type="hidden" name="id" value="app-20170912104309-5902"/>
        <input type="hidden" name="terminate" value="true"/>
        <a href="#" onclick="if (window.confirm('Are you sure you want to kill application app-20170912104309-5902 ?')) { this.parentNode.submit(); return true; } else { return false; }" class="kill-link">(kill)</a>
      </form>
      </td>
      <td>
        <a href="http://COPS-FCO-spark-worker-a-09.cloud.alt:4046">ftliv</a>
      </td>
      <td>
        3
      </td>
      <td sorttable_customkey="1024">
        1024.0 MB
      </td>
      <td>2017/09/12 10:43:09</td>
      <td>spark</td>
      <td>RUNNING</td>
      <td>24.7 h</td>
    </tr></table></body></html>`

	html, _ := html.Parse(strings.NewReader(content))
	fmt.Println(html)

	head := html.FirstChild.FirstChild
	body := head.NextSibling
	table := body.FirstChild
	tbody := table.FirstChild
	tr := tbody.FirstChild
	td := tr.FirstChild

	// act
	res, err := genericTDBrowser(td)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, res)

	line := (*res).Cells
	assert.Equal(t, 8, len(line))
	assert.Equal(t, `<a href="app?appId=app-20170912104309-5902">app-20170912104309-5902</a>`, line[0])
	assert.Equal(t, `<a href="http://COPS-FCO-spark-worker-a-09.cloud.alt:4046">ftliv</a>`, line[1])
	assert.Equal(t, "3", line[2])
	assert.Equal(t, "1024.0 MB", line[3])
	assert.Equal(t, "2017/09/12 10:43:09", line[4])
	assert.Equal(t, "spark", line[5])
	assert.Equal(t, "RUNNING", line[6])
	assert.Equal(t, "24.7 h", line[7])
}

func TestIsActiveSparkMasterWhenOK(t *testing.T) {
	// arrange
	content, _ := readFile("/mock/mainPage/sparkMasterActive.html")

	// act
	res := IsActiveSparkMaster(content)

	// assert
	assert.True(t, res)
}


func TestIsActiveSparkMasterWhenKO(t *testing.T) {
	// arrange
	content, _ := readFile("/mock/myApp/appStreamingStatistics.html")

	// act
	res := IsActiveSparkMaster(content)

	// assert
	assert.False(t, res)
}

// TODO test following functions
// ParseSparkDashboard
// browseTr
// browseTd

func TestReadFile(t *testing.T) {
	// arrange
	expected := "Example file to test readContent function\n\nThere are also new lines !\n"

	// act
	pageContent, err := readFile("/mock/example.test")

	// assert
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
