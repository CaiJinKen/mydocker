package commands

import (
	"fmt"
	"os"

	"github.com/CaiJinKen/mydocker/utils"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "save a container into image",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container name")
		}
		//todo container name/id
		imageName := ctx.Args().Get(0)
		commitContainer(imageName)
		return nil
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
