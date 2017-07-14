package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

// describeCmd represents the list command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe Configurations",
	Long:  `Provides a way to describe the details of a configuration option`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Describe called")
		fmt.Printf("Arguments: %+v\n", args)
		fmt.Printf("Command: %+v\n", cmd)
	},
}

// NewDescribe instantiates a new instance of the configuration describe command.
func NewDescribe(root *cobra.Command) {
	root.AddCommand(describeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
