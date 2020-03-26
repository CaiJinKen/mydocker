package commands

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/CaiJinKen/mydocker/cgroups"
	"github.com/CaiJinKen/mydocker/cgroups/subsystems"

	"github.com/CaiJinKen/mydocker/container"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

var startCommand = cli.Command{
	Name:  "start",
	Usage: "start a container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container identification")
		}

		containerNameOrID := ctx.Args().Get(0)
		startContainer(containerNameOrID)
		return nil
	},
}

func startContainer(containerNameOrID string) {
	info, err := container.GetContainerInfoByIdentification(containerNameOrID)
	if err != nil {
		logrus.Errorf("start container %s error %v", err)
		return
	}
	if info.Status == container.Running {
		logrus.Errorf("container %s is running", containerNameOrID)
		return
	}

	parent, writePipe := container.NewParentProcess(info.Tty, info.Volumes, info.ID)
	if parent == nil {
		logrus.Errorf("new parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		logrus.Errorf("parent process start error %v", err)
	}

	info.Status = container.Running
	info.Pid = strconv.Itoa(parent.Process.Pid)
	info.Save()

	var res subsystems.ResourceConfig
	if info.Resource != "" {
		json.Unmarshal([]byte(info.Resource), &res)
	}

	//set cgroup limit
	cgroupManager := cgroups.NewCgroupManager(info.ID)
	defer cgroupManager.Destroy()
	cgroupManager.Set(&res)
	cgroupManager.Apply(parent.Process.Pid)

	//send command with args to init process
	sendInitCommand(strings.Split(info.Command, " "), writePipe)

	//tty waiting exit: tty container exit / mydocker stop
	if info.Tty {
		parent.Wait()

		//change container status
		info, err := container.GetContainerInfoByIdentification(info.ID)
		if err != nil {
			logrus.Errorf("get container info by identification error %v", err)
		}
		if info != nil {
			info.Status = container.Exit
			info.Save()
		}

		container.CleanUpWorkspace(info.ID)
	}

}
