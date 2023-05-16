package cmd

import "github.com/fatih/color"

type colorFunction func(a ...interface{}) string

var cyan colorFunction = color.New(color.FgCyan).SprintFunc()
var bold colorFunction = color.New(color.Bold).SprintFunc()
var yellow colorFunction = color.New(color.FgYellow).SprintFunc()
