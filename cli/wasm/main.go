//go:build js && wasm
// +build js,wasm

package main

import (
	"github.com/tzapio/tzap/cli/cmd"
)

var Version = "0.0.0-dev"

func main() {
	cmd.RootCmd.Version = Version
	cmd.Execute()
}
