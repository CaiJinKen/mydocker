package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/CaiJinKen/mydocker/container"

	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `create a container with namespace and cgroups limit
		    mydocker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},

	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing command")
		}

		tty := ctx.Bool("ti")
		Run(tty, []string(ctx.Args()))
		return nil
	},
}

//Run fork a new process to start container
func Run(tty bool, cmdArgs []string) {
	logrus.Infof("Run tty %b, args: %v", tty, cmdArgs)
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("new parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		logrus.Errorf("parent process start error %v", err)
	}

	//send command with args to init process
	sendInitCommand(cmdArgs, writePipe)

	parent.Wait()
	os.Exit(0)
}

func sendInitCommand(cmdArgs []string, writePipe *os.File) {
	command := strings.Join(cmdArgs, " ")
	logrus.Infof("command is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
