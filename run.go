/*
3.1中是通过 /proc/self/exec init args 来传递参数，在 init 子命令内去解析参数内容。
旧版本缺点是，如果用户输入的参数很长，或者有一些特殊自负，这种方案可能会失败。
runC 实现的方案是通过匿名管道来实现父子进程之间通信
*/
package main

import (
	"os"
	"strings"

	"github.com/fimreal/mydocker/container"
	log "github.com/sirupsen/logrus"
)

func Run(tty bool, comArray []string) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	sendInitCommand(comArray, writePipe)
	parent.Wait()
	os.Exit(-1)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, "")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
