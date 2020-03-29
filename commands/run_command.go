package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/CaiJinKen/mydocker/utils"

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

		rm := ctx.Bool("rm")

		resConf := &subsystems.ResourceConfig{
			MemoryLimit: ctx.String("m"),
			CpuSet:      ctx.String("cpuset"),
			CpuShare:    ctx.String("cpushare"),
		}

		volumeURLs := ctx.StringSlice("v")
		logrus.Infof("volumeURLs: %v", volumeURLs)

		containerName := ctx.String("name")
		logrus.Infof("containerName: %s", containerName)

		Run(tty, rm, []string(ctx.Args()), resConf, volumeURLs, containerName)
		return nil
	},
}

//Run fork a new process to start container
func Run(tty, rm bool, cmdArgs []string, res *subsystems.ResourceConfig, volumeURLs []string, containerName string) {
	logrus.Infof("Run tty %b, args: %v", tty, cmdArgs)
	containerID := container.GenerateUUID()
	parent, writePipe := container.NewParentProcess(tty, volumeURLs, containerID)
	if parent == nil {
		logrus.Errorf("new parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		logrus.Errorf("parent process start error %v", err)
	}

	//log container info
	logrus.Infof("container pid is %d", parent.Process.Pid)
	containerInfo := &container.Info{
		Pid:       strconv.Itoa(parent.Process.Pid),
		ID:        containerID,
		Name:      containerName,
		Command:   strings.Join(cmdArgs, " "),
		Status:    container.Running,
		Rm:        rm,
		Tty:       tty,
		Volumes:   volumeURLs,
		Resource:  res.String(),
		CreatedAt: utils.TimeNowString(),
	}

	if err := containerInfo.Save(); err != nil {
		logrus.Errorf("record container info error %v", err)
		return
	}

	//set cgroup limit
	cgroupManager := cgroups.NewCgroupManager(containerID)
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	//send command with args to init process
	sendInitCommand(cmdArgs, writePipe)

	//tty waiting exit: tty container exit / mydocker stop
	if tty {
		parent.Wait()

		//change container status
		info, err := container.GetContainerInfoByIdentification(containerID)
		if err != nil {
			logrus.Errorf("get container info by identification error %v", err)
		}
		if info != nil {
			info.Status = container.Exit
			info.Pid = ""
			info.Save()
		}

		container.CleanUpWorkspace(containerID)
	}
	//parent exit

}

func sendInitCommand(cmdArgs []string, writePipe *os.File) {
	command := strings.Join(cmdArgs, " ")
	logrus.Infof("command is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
