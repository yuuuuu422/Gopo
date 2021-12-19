package cmd

import (
	"Gopo/pocs/scripts"
	"Gopo/utils"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
	"strings"
)

const (
	AUTHOR  = "Theoyu"
	VERSION = "0.1.0"
	EMAIL = "zeyu.ou@foxmail.com"
)

var (
	pocFile    string
	pocRules   string
	pocScript string
	target     string
	targetFile string
	targets    []string
	proxy      string
	//timeout time.Duration
	num       int
	cookie    string
	httpDebug bool
	debug     bool

)

func Hacking() {
	app := &cli.App{
		Name:    "Gopo",
		//UsageText: "hhh",
		Version: VERSION,
		Authors: []*cli.Author{
			{
				Name:  AUTHOR,
				Email: EMAIL,
			},
		},
		Usage: "Proof of concept",
		Commands: cli.Commands{
			{
				Name: "rule",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "poc-rules",
						Destination: &pocRules,
						Usage:       "load multi pocs from `rules`,eg: struts2 or thinkphp ",
						Aliases:     []string{"r"},
					},
					&cli.StringFlag{
						Name:        "poc-file",
						Destination: &pocFile,
						Usage:       "load single poc from poc-rules ",
						Aliases:     []string{"p"},
					},
					&cli.StringFlag{
						Name:        "target",
						Destination: &target,
						Usage:       "target to scan",
						Aliases:     []string{"t"},
					},
					&cli.StringFlag{
						Name:        "target-file",
						Destination: &targetFile,
						Usage:       "load target `FILE` to scan",
						Aliases:     []string{"f"},
					},
					&cli.StringFlag{
						Name:        "proxy",
						Destination: &proxy,
						Usage:       "http proxy",
					},
					//&cli.DurationFlag{
					//	Name:  "timeout",
					//	Value: 2*time.Second,
					//	Destination: &timeout,
					//	Usage: "scan `TIMEOUT`",
					//},
					&cli.IntFlag{
						Name:        "num",
						Aliases:     []string{"n"},
						Value:       20,
						Destination: &num,
						Usage:       "threats `NUM` to scan",
					},
					&cli.StringFlag{
						Name:        "cookie",
						Destination: &cookie,
						Usage:       "http cookie",
					},
					&cli.BoolFlag{
						Name:        "httpDebug",
						Destination: &httpDebug,
						Value:       false,
						Usage:       "http debug",
					},
					&cli.BoolFlag{
						Name:        "debug",
						Destination: &debug,
						Value:       false,
						Usage:       "set the log debug level",
					},
				},
				Action: func(newContext *cli.Context) error {
					utils.InitLog(debug)
					if !utils.InitCeyeApi(){
						utils.Warning("init ceye platform false")
					}
					switch {
					case targetFile != "":
						targets = utils.ReadingLines(targetFile)
					case target != "":
						if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
						} else {
							target = "http://" + target
						}
						targets = append(targets, target)
					}
					var pocFiles []string
					var err error
					if pocFile != "" {
						pocFiles = append(pocFiles, "pocs/rules/"+pocRules+"/"+pocFile)
					} else {
						pocFiles, err = utils.LoadRules(pocRules)
						if err != nil {
							return err
						}
					}
					pocs, err := utils.ParseRules(pocFiles)
					if err != nil {
						return err
					}
					err = utils.InitHttp(cookie, proxy, httpDebug)
					if err != nil {
						return err
					}
					utils.CheckVuls(pocs, targets, num)
					return nil
				},
			},

			{
				Name: "script",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "poc",
						Destination: &pocScript,
						Usage:       "load poc script from scripts",
						Aliases:     []string{"p"},
					},
					&cli.StringFlag{
						Name:        "proxy",
						Destination: &proxy,
						Usage:       "http proxy",
					},
					&cli.IntFlag{
						Name:        "num",
						Aliases:     []string{"n"},
						Value:       20,
						Destination: &num,
						Usage:       "threats `NUM` to scan",
					},
					&cli.StringFlag{
						Name:        "cookie",
						Destination: &cookie,
						Usage:       "http cookie",
					},
					&cli.StringFlag{
						Name:        "target",
						Destination: &target,
						Usage:       "target to scan",
						Aliases:     []string{"t"},
					},
					&cli.StringFlag{
						Name:        "target-file",
						Destination: &targetFile,
						Usage:       "load target `FILE` to scan",
						Aliases:     []string{"f"},
					},
					&cli.BoolFlag{
						Name:        "httpDebug",
						Destination: &httpDebug,
						Value:       false,
						Usage:       "http debug",
					},
					&cli.BoolFlag{
						Name:        "debug",
						Destination: &debug,
						Value:       false,
						Usage:       "set the log debug level",
					},
				},
				Action: func(context *cli.Context) error {
					utils.InitLog(debug)
					if !utils.InitCeyeApi(){
						utils.Warning("init ceye platform false")
					}
					err := utils.InitHttp(cookie, proxy, httpDebug)
					if err != nil {
						return err
					}
					switch {
					case targetFile != "":
						targets = utils.ReadingLines(targetFile)
					case target != "":
						if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
						} else {
							target = "http://" + target
						}
						targets = append(targets, target)
					}
					scriptFunc:= scripts.ScriptInit(pocScript,num)
					if scriptFunc == nil{
						utils.Yellow("Unsupported script,Please check the list:")
						scripts.ShowRegister()
						return nil
					}
					for _,target=range targets{
						scriptFunc(target)
					}
					return nil
				},
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		utils.Yellow("%v",err)
	}
}
