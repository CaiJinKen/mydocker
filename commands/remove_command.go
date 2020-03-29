package commands

import (
	"fmt"

	"github.com/CaiJinKen/mydocker/container"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli"
)

var removeCommand = cli.Command{
	Name:  "rm",
	Usage: "remove unused container",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container name")
		}

		containerNameOrID := ctx.Args().Get(0)
		removeContainer(containerNameOrID)
		return nil
	},
}

func removeContainer(containerNameOrID string) {
	info, err := container.GetContainerInfoByIdentification(containerNameOrID)
	if err != nil {
		logrus.Errorf("remove container %s error %v", containerNameOrID, err)
		return
	}

	if info.Status != container.Stop {
		logrus.Errorf("couldn`t remove running container")
		return
	}

	container.DeleteContainerInfo(info.ID)
}
