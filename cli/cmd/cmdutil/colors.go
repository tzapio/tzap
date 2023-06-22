package cmdutil

import "github.com/fatih/color"

type colorFunction func(a ...interface{}) string

var Cyan colorFunction = color.New(color.FgCyan).SprintFunc()
var Bold colorFunction = color.New(color.Bold).SprintFunc()
var Yellow colorFunction = color.New(color.FgYellow).SprintFunc()
var Black colorFunction = color.New(color.FgBlack).SprintFunc()
