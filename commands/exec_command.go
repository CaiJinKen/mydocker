package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/CaiJinKen/mydocker/container"
	_ "github.com/CaiJinKen/mydocker/nssetter"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var execCommand = cli.Command{
	Name:  "exec",
	Usage: "enter into a container",
	Action: func(ctx *cli.Context) error {
		if os.Getenv(EnvExecPID) != "" {
			logrus.Infof("pid callback is %s", os.Getgid())
			return nil
		}

		if len(ctx.Args()) < 2 {
			return fmt.Errorf("missing container name or command")
		}

		containerName := ctx.Args().Get(0)
		cmdArgs := ctx.Args()[1:]

		ExecContainer(containerName, cmdArgs)

		return nil
	},
}

const (
	EnvExecPID = "mydocker_pid"
	EnvExecCMD = "mydocker_cmd"
)

func ExecContainer(containerName string, cmdArgs []string) {
	pid, err := container.GetContainerPidByName(containerName)
	if err != nil {
		logrus.Errorf("get container pid by name error %v", err)
		return
	}

	cmdStr := strings.Join(cmdArgs, " ")
	logrus.Infof("exec command is %s pid is %s", cmdStr, pid)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := os.Setenv(EnvExecPID, pid); err != nil {
		logrus.Errorf("set mydocker_pid env error %v", err)
	}
	if err := os.Setenv(EnvExecCMD, cmdStr); err != nil {
		logrus.Errorf("set mydocker_cmd env error %v")
	}

	if err := cmd.Run(); err != nil {
		logrus.Errorf("exec container %s error %v", containerName, err)
	}
}
