package utils

import (
	"Gopo/utils/proto"
	"context"
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/xiecat/xhttp"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Poc struct {
	Name       string            `yaml:"name"`
	Transport  string            `yaml:"transport"`
	Set        map[string]string `yaml:"set"`
	Rules      map[string]*Rules `yaml:"rules"`
	Expression string            `yaml:"expression"`
	Details    Details           `yaml:"details"`
}

type Rules struct {
	Name    string `yaml:"-"`
	Request struct {
		Method          string            `yaml:"method"`
		Path            string            `yaml:"path"`
		Headers         map[string]string `yaml:"headers"`
		Body            string            `yaml:"body"`
		FollowRedirects bool              `yaml:"follow_redirects"`
	} `yaml:"request"`
	Expression string `yaml:"expression"`
}

type Details struct {
	Link string `yaml:"link"`
}

type Task struct {
	Poc    *Poc
	Target string
}

func LoadRules(pocName string) ([]string, error) {
	pocDir := "pocs/rules/" + pocName
	return GetAllRules(pocDir)
}

func ParseRules(pocFiles []string) ([]*Poc, error) {
	var pocs []*Poc
	for _, file := range pocFiles {
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		p := &Poc{}
		err = yaml.Unmarshal(yamlFile, p)
		if err != nil {
			return nil, err
		}
		pocs = append(pocs, p)
	}
	return pocs, nil
}

func CheckVuls(pocs []*Poc, targets []string, threads int, needTicker bool) {
	var wg sync.WaitGroup
	var tasks []Task
	taskChan := make(chan Task, threads)
	for _, target := range targets {
		for _, poc := range pocs {
			task := Task{
				poc, target,
			}
			tasks = append(tasks, task)
		}
	}

	worker := func(taskChan <-chan Task, wg *sync.WaitGroup) {
		for task := range taskChan {
			isVul, err := execPoc(task)
			if err != nil {
				Error(err)
			}
			if isVul {
				Green("%v find %v", task.Target, task.Poc.Name)
			}
			wg.Done()
		}
	}

	for i := 0; i < threads; i++ {
		go worker(taskChan, &wg)
	}

	if needTicker {
		ticker := time.NewTicker(time.Second)
		for _, task := range tasks {
			<-ticker.C
			wg.Add(1)
			taskChan <- task
		}
	} else {
		for _, task := range tasks {
			wg.Add(1)
			taskChan <- task
		}
	}
	wg.Wait()
}

func execPoc(task Task) (bool, error) {
	variableMap := make(map[string]interface{})
	target := task.Target
	poc := *task.Poc

	options := newEnvOption()
	options.UpdateCompileOptions(poc.Set)
	env, err := cel.NewEnv(cel.Lib(options))
	if err != nil {
		return false, err
	}

	for k := range poc.Set {
		expression := poc.Set[k]
		if expression == "newReverse()" {
			variableMap[k] = newReverse()
			continue
		}
		out, err := Evaluate(env, expression, variableMap)
		if err != nil {
			Error(err)
			continue
		}
		switch value := out.Value().(type) {
		case *proto.UrlType:
			variableMap[k] = UrlTypeToString(value)
		case int64:
			variableMap[k] = int(value)
		default:
			variableMap[k] = fmt.Sprintf("%v", out)
		}
	}

	for name, rule := range poc.Rules {
		//把Set设置变量的值带入
		for k, v := range variableMap {
			rule.Request.Path = strings.ReplaceAll(rule.Request.Path, "{{"+k+"}}", fmt.Sprintf("%v", v))
			rule.Request.Body = strings.ReplaceAll(rule.Request.Body, "{{"+k+"}}", fmt.Sprintf("%v", v))
		}
		hr, _ := http.NewRequest(rule.Request.Method, target+rule.Request.Path, strings.NewReader(rule.Request.Body))
		//修改header
		for k, v := range rule.Request.Headers {
			hr.Header.Set(k, v)
		}

		req := &xhttp.Request{RawRequest: hr}
		ctx := context.Background()
		oResp, err := Client.Do(ctx, req)

		if err != nil {
			options.UpdateFunctionOptions(name, types.False)
			continue
		}
		resp := ParseResponse(oResp)
		variableMap["response"] = resp
		out, err := Evaluate(env, rule.Expression, variableMap)
		if err != nil {
			return false, err
		}
		//把r0()声明到options里
		options.UpdateFunctionOptions(name, out)
	}
	//重写创建env 加载r0等表达式
	env, err = cel.NewEnv(cel.Lib(options))
	if err != nil {
		return false, err
	}
	out, err := Evaluate(env, poc.Expression, variableMap)
	if err != nil {
		return false, err
	}
	if out.Value() == false {
		return false, nil
	}
	return true, nil
}

func LoadScripts() {

}

/*
94.
*/
