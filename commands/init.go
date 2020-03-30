package commands

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "mydocker",
	Short: "mydocker application",
	Long:  "mydocker is a docker-like container just for learn",
}

var (
	//run
	tty, detach, rm                         bool
	containerName, memory, cpushare, cpuset string
	volumeMappings                          []string

	//ps
	all bool
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)

	runCommand.Flags().BoolVar(&tty, "ti", false, "enable tty")
	runCommand.Flags().BoolVarP(&detach, "detach", "d", false, "detach container")
	runCommand.Flags().StringVar(&containerName, "name", "", "container name")
	runCommand.Flags().StringVar(&memory, "m", "", "memory limit")
	runCommand.Flags().StringVar(&cpushare, "cpushare", "", "cpushare limit")
	runCommand.Flags().StringVar(&cpuset, "cpuset", "", "cpuset limit")
	runCommand.Flags().StringSliceVarP(&volumeMappings, "volume", "v", nil, "volume mapping")
	runCommand.Flags().BoolVar(&rm, "rm", false, "delete container after stop")

	psCommand.Flags().BoolVarP(&all, "all", "a", false, "list all container")

	rootCmd.AddCommand(
		runCommand,
		initCommand,
		commitCommand,
		psCommand,
		execCommand,
		stopCommand,
		removeCommand,
		startCommand,
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorf("start mydocker error %v", err)
		os.Exit(1)
	}
}
