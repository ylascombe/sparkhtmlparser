package analyser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestgetTableBody(t *testing.T) {

	htm := `<!DOCTYPE html>
<html><head>
        <meta http-equiv="Content-type" content="text/html; charset=UTF-8"><link rel="stylesheet" href="colis360%20-%20Streaming%20Statistics_fichiers/bootstrap.css" type="text/css"><link rel="stylesheet" href="colis360%20-%20Streaming%20Statistics_fichiers/vis.css" type="text/css"><link rel="stylesheet" href="colis360%20-%20Streaming%20Statistics_fichiers/webui.css" type="text/css"><link rel="stylesheet" href="colis360%20-%20Streaming%20Statistics_fichiers/timeline-view.css" type="text/css"><script src="colis360%20-%20Streaming%20Statistics_fichiers/sorttable.js"></script><script src="colis360%20-%20Streaming%20Statistics_fichiers/jquery-1.js"></script><script src="colis360%20-%20Streaming%20Statistics_fichiers/vis.js"></script><script src="colis360%20-%20Streaming%20Statistics_fichiers/bootstrap-tooltip.js"></script><script src="colis360%20-%20Streaming%20Statistics_fichiers/initialize-tooltips.js"></script><script src="colis360%20-%20Streaming%20Statistics_fichiers/table.js"></script><script src="colis360%20-%20Streaming%20Statistics_fichiers/additional-metrics.js"></script><script src="colis360%20-%20Streaming%20Statistics_fichiers/timeline-view.js"></script>

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
            <ul class="nav"><li class="">
        <a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/jobs/">Jobs</a>
      </li><li class="">
        <a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/stages/">Stages</a>
      </li><li class="">
        <a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/storage/">Storage</a>
      </li><li class="">
        <a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/environment/">Environment</a>
      </li><li class="">
        <a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/executors/">Executors</a>
      </li><li class="active">
        <a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/streaming/">Streaming</a>
      </li></ul>
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
          <script src="colis360%20-%20Streaming%20Statistics_fichiers/d3.js"></script><link rel="stylesheet" href="colis360%20-%20Streaming%20Statistics_fichiers/streaming-page.css" type="text/css"><script src="colis360%20-%20Streaming%20Statistics_fichiers/streaming-page.js"></script><div>Running batches of
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
    </div><br>

    <table id="active-batches-table" class="table table-bordered table-striped table-condensed sortable">
      <thead>
        <tr><th>Batch Time</th><th>Input Size</th><th>Scheduling Delay
        <sup>
      (<a data-toggle="tooltip" data-placement="top" title="" data-original-title="Time taken by Streaming scheduler to submit jobs of a batch">?</a>)
    </sup>
      </th><th>Processing Time
        <sup>
      (<a data-toggle="tooltip" data-placement="top" title="" data-original-title="Time taken to process all jobs of a batch">?</a>)
    </sup></th><th>Status</th>
      </tr></thead>
      <tbody>
        <tr><td id="batch-1504876710000" sorttable_customkey="1504876710000">
      <a href="http://cops-fco-spark-worker-a-11.cloud.alt:4041/streaming/batch?id=1504876710000">
        2017/09/08 15:18:30
      </a>
    </td><td sorttable_customkey="1823">1823 events</td><td sorttable_customkey="6">
        6 ms
      </td><td sorttable_customkey="9223372036854775807">
        -
      </td><td>processing</td></tr>
      </tbody>
    <tfoot></tfoot></table>


</body></html>`

	res, err := getTableBody(htm)

	assert.Nil(t, err)
	assert.Equal(t, "id", res.Attr[0].Key)
	assert.Equal(t, "active-batches-table", res.Attr[0].Val)

}
