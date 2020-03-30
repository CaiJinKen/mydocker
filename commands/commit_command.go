package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/CaiJinKen/mydocker/utils"
	"github.com/sirupsen/logrus"
)

var commitCommand = &cobra.Command{
	Use:   "commit [image]",
	Short: "save a container into image",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			logrus.Errorf("missing container identification")
			return
		}
		commitContainer(args[0])
	},
}

func commitContainer(imageName string) {
	mntURL := "/root/mnt"
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
