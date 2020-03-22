package commands

import "github.com/urfave/cli"

//Commands mydocker commands
var Commands = []cli.Command{
	runCommand,
	initCommand,
	commitCommand,
	psCommand,
}
