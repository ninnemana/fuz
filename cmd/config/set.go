package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set Configuration",
	Long:  `Updates the local environment`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Set called")
		fmt.Printf("Arguments: %+v\n", args)
		fmt.Printf("Command: %+v\n", cmd)
	},
}

// NewSet instantiates a new instance of the configuration set command.
func NewSet(root *cobra.Command) {
	root.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
