package main

import (
	"fmt"

	"github.com/frayer/concourse-resource-tekton-trigger/bridge/commands"
	"github.com/jessevdk/go-flags"
)

func main() {
	parser := flags.NewParser(&commands.BridgeCommand{}, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	_, err := parser.Parse()

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
