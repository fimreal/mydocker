// 3.3 版本更新，传参使用pipe，避免过长或者特殊参数无法解析的问题。
package main

import (
	"fmt"

	"github.com/fimreal/mydocker/container"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			mydocker run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
	},
	/*
		这里是 run 命令执行的真正函数。
		1. 判断参数是否包含command
		2. 获取用户指定的 command
		3. 调用 Run function 去准备启动容器
	*/
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		// cmd := context.Args().Get(0)
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("it")
		Run(tty, cmdArray)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	/*
		1. 获取传递过来的 command 参数
		2. 执行容器初始化操作
	*/
	Action: func(context *cli.Context) error {
		log.Infof("init come on")
		// 下面这部分放到 run.go 中来执行
		// cmd := context.Args().Get(0)
		// log.Infof("command %s", cmd)
		// err := container.RunContainerInitProcess(cmd, nil)
		err := container.RunContainerInitProcess()
		return err
	},
}
