/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/Fonzeca/TrackinRemoteConfig/server"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "A brief description of your command",
	Aliases: []string{"ls"},

	Run: func(cmd *cobra.Command, args []string) {
		if len(server.ConnectionPool) == 0 {
			fmt.Println("No hay conexiones")
			return
		}

		for k := range server.ConnectionPool {
			fmt.Println(k)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
