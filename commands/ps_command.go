package commands

import (
	"github.com/CaiJinKen/mydocker/container"
	"github.com/urfave/cli"
)

var psCommand = cli.Command{
	Name:  "ps",
	Usage: "list container ps",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "a",
			Usage: "list all container",
		},
	},
	Action: func(ctx *cli.Context) error {
		all := ctx.Bool("a")
		container.ListContainer(all)
		return nil
	},
}
