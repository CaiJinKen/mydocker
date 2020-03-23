package commands

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/CaiJinKen/mydocker/container"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

var stopCommand = cli.Command{
	Name:  "stop",
	Usage: "stop container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container name")
		}

		containerNameOrID := ctx.Args().Get(0)
		stopContainer(containerNameOrID)
		return nil
	},
}

func stopContainer(containerNameOrID string) {
	info, err := container.GetContainerInfoByIdentification(containerNameOrID)
	if err != nil {
		logrus.Errorf("get container pid error %v", err)
		return
	}

	pid, err := strconv.Atoi(info.Pid)
	if err != nil {
		logrus.Errorf("convert pid to int error %v", err)
		return
	}

	if err = syscall.Kill(pid, syscall.SIGTERM); err != nil {
		logrus.Errorf("stop container %s error %v", containerNameOrID, err)
	}

	info.Stop()
}
