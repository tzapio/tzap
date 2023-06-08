package main

import (
	"bytes"
	"text/template"

	"github.com/tzapio/tzap/pkg/tzap"
)

var tp = `{{ 
	
}}`

func main() {
	t := tzap.InternalNew()
	t = t.AddAssistantMessage("hello")
	// create template
	tz := template.Must(template.New("template").Parse(tp))
	t.AddUserMessage()
	var output bytes.Buffer
	err := tz.Execute(&output, t)
	if err != nil {
		panic(err)
	}

	println(string(output.Bytes()))
}
