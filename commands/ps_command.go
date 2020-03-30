package commands

import (
	"github.com/CaiJinKen/mydocker/container"
	"github.com/spf13/cobra"
)

var psCommand = &cobra.Command{
	Use:   "ps",
	Short: "list container ps",
	Run: func(cmd *cobra.Command, args []string) {
		container.ListContainer(all)
	},
}
