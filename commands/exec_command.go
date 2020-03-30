package commands

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/CaiJinKen/mydocker/container"
	_ "github.com/CaiJinKen/mydocker/nssetter"

	"github.com/sirupsen/logrus"
)

var execCommand = &cobra.Command{
	Use:   "exec [container] [command]",
	Short: "enter into a container",
	Run: func(cmd *cobra.Command, args []string) {
		if os.Getenv(EnvExecPID) != "" {
			logrus.Infof("pid callback is %s", os.Getgid())
			return
		}

		if len(args) < 2 {
			logrus.Errorf("missing container name or command")
			return
		}

		containerNameOrID := args[0]
		cmdArgs := args[1:]

		ExecContainer(containerNameOrID, cmdArgs)
	},
}

const (
	EnvExecPID = "mydocker_pid"
	EnvExecCMD = "mydocker_cmd"
)

func ExecContainer(containerNameOrID string, cmdArgs []string) {
	pid, err := container.GetContainerPidByIdentification(containerNameOrID)
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
		logrus.Errorf("exec container %s error %v", containerNameOrID, err)
	}
}
