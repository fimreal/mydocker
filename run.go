/*
这里是父进程，也就是当前进程执行的内容
1. 这里的 /proc/self/exe 调用中，/proc/self 指的是当前运行进程自己的环境，exec 其实就是自己调用了自己，使用这种方式对创建出来的进程初始化
2. 后面的 args 是参数，其中 init 是传递给本进程的第一个参数，在本例中，其实就是会去调用 initCommand 去初始化进程的一些环境和资源
3. 下面的 clone 参数是去 fork 出来一个新的进程，并且使用了 namespace 隔离新建进程和外部环境
4. 如果用户制定了 -it 参数，就需要把当前进程的输入输出等导入到系统的标准输入输出上。
*/
package main

import (
	"os"
	"strings"

	"github.com/fimreal/mydocker/cgroups"
	"github.com/fimreal/mydocker/cgroups/subsystems"
	"github.com/fimreal/mydocker/container"
	log "github.com/sirupsen/logrus"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	cgroupManager := cgroups.NewCgroupManager("mydocker")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	parent.Wait()
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

// // 最简单的实现，只包含 --it 克隆终端
// func Run(tty bool, command string) {
// 	parent := container.NewParentProcess(tty, command)
// 	if err := parent.Start(); err != nil {
// 		log.Error(err)
// 	}
// 	parent.Wait()
// 	os.Exit(-1)
// }
