/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/Fonzeca/TrackinRemoteConfig/server"
	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	cobraprompt "github.com/stromland/cobra-prompt"
)

var configDeviceAnnotationValue = "configDevice"

var GetDevicesConnectedDynamic = func(annotationValue string) []prompt.Suggest {
	if annotationValue != configDeviceAnnotationValue {
		return nil
	}

	suggestions := []prompt.Suggest{}

	for k := range server.ConnectionPool {
		suggestions = append(suggestions, prompt.Suggest{Text: k, Description: "Device connected"})
	}

	return suggestions
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Annotations: map[string]string{
		cobraprompt.DynamicSuggestionsAnnotation: configDeviceAnnotationValue,
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 || len(args) > 2 {
			fmt.Println("Se necesita dos argumentos")
			return
		}

		imei := args[0]
		command := args[1]

		if server.ConnectionPool[imei] == nil {
			fmt.Println("La conexion con imei " + imei + " no existe")
			return
		}

		pipeIn := server.ConnectionPool[imei][0]
		pipeOut := server.ConnectionPool[imei][1]

		pipeIn <- command

		respuesta := <-pipeOut
		fmt.Println(respuesta)

	},
}

func init() {
	RootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
