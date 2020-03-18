package commands

import (
	"github.com/CaiJinKen/mydocker/container"
	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "init container process",
	Action: func(ctx *cli.Context) error {
		return container.RunContainerInitProcess()
	},
}
