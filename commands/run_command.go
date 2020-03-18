package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/CaiJinKen/mydocker/cgroups"

	"github.com/CaiJinKen/mydocker/cgroups/subsystems"

	"github.com/sirupsen/logrus"

	"github.com/CaiJinKen/mydocker/container"

	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `create a container with namespace and cgroups limit
		    mydocker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},

	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing command")
		}

		resConf := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("m"),
			CpuSet:      ctx.String("cpuset"),
			CpuShare:    ctx.String("cpushare"),
		}

		tty := ctx.Bool("ti")
		Run(tty, []string(ctx.Args()), resConf)
		return nil
	},
}

//Run fork a new process to start container
func Run(tty bool, cmdArgs []string, res *subsystems.ResourceConfig) {
	logrus.Infof("Run tty %b, args: %v", tty, cmdArgs)
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("new parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		logrus.Errorf("parent process start error %v", err)
	}

	//set cgroup limit
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	//send command with args to init process
	sendInitCommand(cmdArgs, writePipe)

	parent.Wait()
	os.Exit(0)
}

func sendInitCommand(cmdArgs []string, writePipe *os.File) {
	command := strings.Join(cmdArgs, " ")
	logrus.Infof("command is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
