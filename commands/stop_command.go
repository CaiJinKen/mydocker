package commands

import (
	"strconv"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/CaiJinKen/mydocker/container"
	"github.com/sirupsen/logrus"
)

var stopCommand = &cobra.Command{
	Use:   "stop [container]",
	Short: "stop container",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Errorf("missing container name")
			return
		}

		containerNameOrID := args[0]
		stopContainer(containerNameOrID)
	},
}

func stopContainer(containerNameOrID string) {
	info, err := container.GetContainerInfoByIdentification(containerNameOrID)
	if err != nil {
		logrus.Errorf("get container pid error %v", err)
		return
	}

	if info.Tty {
		logrus.Errorf("container used tty param, please stop it into container")
		return
	}

	if info.Status != container.Running {
		return
	}

	pid, err := strconv.Atoi(info.Pid)
	if err != nil {
		logrus.Errorf("convert pid to int error %v", err)
	}

	if err = syscall.Kill(pid, syscall.SIGTERM); err != nil {
		logrus.Errorf("stop container %s error %v", containerNameOrID, err)
	}

	info.Stop()
	container.CleanUpWorkspace(info.ID)
}
