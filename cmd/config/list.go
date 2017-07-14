package config

import (
	"context"
	"fmt"
	"strings"

	hostsServ "github.com/ninnemana/fuz/hosts/service"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Configurations",
	Long:  `Provides a list of the available configurations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("list called")
		fmt.Printf("Arguments: %+v\n", args)
		fmt.Printf("Command: %+v\n", cmd)

		s, err := hostsServ.New()
		if err != nil {
			return err
		}

		records, err := s.List(context.Background(), nil)
		if err != nil {
			return err
		}

		for _, rec := range records {
			fmt.Printf("%s\t[%s]\n", rec.LocalPtr, strings.Join(rec.Hosts, "/"))
		}

		return nil
	},
}

// NewList instantiates the configuration list command
func NewList(root *cobra.Command) {
	root.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
