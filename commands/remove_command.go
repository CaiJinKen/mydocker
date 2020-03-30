package commands

import (
	"github.com/spf13/cobra"

	"github.com/CaiJinKen/mydocker/container"
	"github.com/sirupsen/logrus"
)

var removeCommand = &cobra.Command{
	Use:   "rm [container]",
	Short: "remove unused container",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Errorf("missing container name")
			return
		}

		containerNameOrID := args[0]
		removeContainer(containerNameOrID)
	},
}

func removeContainer(containerNameOrID string) {
	info, err := container.GetContainerInfoByIdentification(containerNameOrID)
	if err != nil {
		logrus.Errorf("remove container %s error %v", containerNameOrID, err)
		return
	}

	if info.Status == container.Running {
		logrus.Errorf("couldn`t remove running container")
		return
	}

	container.DeleteContainerInfo(info.ID)
}
