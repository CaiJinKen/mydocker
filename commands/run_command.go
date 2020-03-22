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
		cli.StringSliceFlag{
			Name:  "v",
			Usage: "volume mapping",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		cli.BoolFlag{
			Name:  "rm",
			Usage: "delete container after stop",
		},
	},

	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing command")
		}

		tty := ctx.Bool("ti")
		detach := ctx.Bool("d")
		if tty && detach {
			return fmt.Errorf("cannot use ti and d at same time")
		}

		resConf := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("m"),
			CpuSet:      ctx.String("cpuset"),
			CpuShare:    ctx.String("cpushare"),
		}

		volumeURLs := ctx.StringSlice("v")
		logrus.Infof("volumeURLs: %v", volumeURLs)

		containerName := ctx.String("name")
		logrus.Infof("containerName: %s", containerName)

		Run(tty, detach, []string(ctx.Args()), resConf, volumeURLs, containerName)
		return nil
	},
}

//Run fork a new process to start container
func Run(tty, detach bool, cmdArgs []string, res *subsystems.ResourceConfig, volumeURLs []string, containerName string) {
	logrus.Infof("Run tty %b, args: %v", tty, cmdArgs)
	parent, writePipe := container.NewParentProcess(tty, volumeURLs)
	if parent == nil {
		logrus.Errorf("new parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		logrus.Errorf("parent process start error %v", err)
	}

	//log container info
	containerID, err := container.RecordContainerInfo(parent.Process.Pid, containerName, cmdArgs)
	if err != nil {
		logrus.Errorf("record container info error %v", err)
		return
	}

	//set cgroup limit
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	//send command with args to init process
	sendInitCommand(cmdArgs, writePipe)

	if tty {
		parent.Wait()
		container.DeleteContainerInfo(containerID)
	}

	//delete container workspace
	//rootURL := "/root/"
	//mntURL := "/root/mnt"
	//container.DeleteWorkSpace(rootURL, mntURL, volumeURLs)

	//os.Exit(0)
}

func sendInitCommand(cmdArgs []string, writePipe *os.File) {
	command := strings.Join(cmdArgs, " ")
	logrus.Infof("command is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
