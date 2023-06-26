/*
The main package provides the entry point to the application.
The version information is set here and passed to the cmd package 
to be used by the root command.
*/

package main

import "github.com/tzapio/tzap/cli/cmd"

// Version represents the current version of the application.
var Version = "0.0.0-dev"

// The main function is the entry point to the application.
func main() {
	// Set the application version in the root command.
	cmd.RootCmd.Version = Version
	
	// Execute the root command.
	cmd.Execute()
}