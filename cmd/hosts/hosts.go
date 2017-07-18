package hosts

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Root represents the config command
var Root = &cobra.Command{
	Use:   "hosts",
	Short: "Allows acces to the local hosts file",
	Long:  `Gives all available commands to the current user for interacting with their hosts file for
	maintaining network policies.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hosts called")
	},
}

// Init intializes the configuration command and all it's children.
func Init(root *cobra.Command) {
	NewList(Root)
	root.AddCommand(Root)
}