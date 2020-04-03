package commands

import (
	"fmt"
	"os"

	"github.com/CaiJinKen/mydocker/container"

	"github.com/spf13/cobra"

	"github.com/CaiJinKen/mydocker/utils"
	"github.com/sirupsen/logrus"
)

var commitCommand = &cobra.Command{
	Use:   "commit [container] [image]",
	Short: "save a container into image",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			logrus.Errorf("missing container identification")
			return
		}
		commitContainer(args[0], args[1])
	},
}

func commitContainer(containerIdentify, imageName string) {
	info, err := container.GetContainerInfoByIdentification(containerIdentify)
	if err != nil {
		return
	}
	mntURL := fmt.Sprintf(container.MntURL, info.ID)
	currentPath, err := os.Getwd()
	if err != nil {
		logrus.Errorf("get current dir error %v", err)
		return
	}
	imageTar := fmt.Sprintf("%s/%s.tar", currentPath, imageName)
	logrus.Infof("tar file is %s", imageTar)

	_, err = utils.Exec("tar", "-czf", imageTar, "-C", mntURL, ".")
	if err != nil {
		logrus.Errorf("tar dir %s error %v", mntURL, err)
	}
}
