package main

import (
	"os"

	"github.com/CaiJinKen/mydocker/commands"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = "mydocker is a docker like container"

	//register command
	app.Commands = commands.Commands

	//setup logrus
	app.Before = func(ctx *cli.Context) error {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Errorf("app run error %v", err)
	}
}
