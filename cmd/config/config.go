package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Root represents the config command
var Root = &cobra.Command{
	Use:   "config",
	Short: "Describes the current configuration of your private cloud platform",
	Long:  `All necessary configuration paramters required to setup your private cloud platform.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

// Init intializes the configuration command and all it's children.
func Init(root *cobra.Command) {
	NewList(Root)
	NewSet(Root)
	NewDescribe(Root)
	root.AddCommand(Root)
}
