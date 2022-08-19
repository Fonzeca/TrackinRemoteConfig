/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"os"
	"strings"

	"github.com/Fonzeca/TrackinRemoteConfig/cmd"
	"github.com/Fonzeca/TrackinRemoteConfig/server"
	"github.com/c-bata/go-prompt"
	cobraprompt "github.com/stromland/cobra-prompt"
)

var advancedPrompt = &cobraprompt.CobraPrompt{
	RootCmd:                  cmd.RootCmd,
	PersistFlagValues:        true,
	ShowHelpCommandAndFlags:  true,
	DisableCompletionCommand: true,
	AddDefaultExitCommand:    false,
	GoPromptOptions:          []prompt.Option{},
	DynamicSuggestionsFunc: func(annotationValue string, document *prompt.Document) []prompt.Suggest {
		if suggestions := cmd.GetDevicesConnectedDynamic(annotationValue); suggestions != nil {
			return suggestions
		}

		return []prompt.Suggest{}
	},
	OnErrorFunc: func(err error) {
		if strings.Contains(err.Error(), "unknown command") {
			cmd.RootCmd.PrintErrln(err)
			return
		}

		cmd.RootCmd.PrintErr(err)
		os.Exit(1)
	},
}

func main() {
	go server.Serve()
	advancedPrompt.Run()
}
