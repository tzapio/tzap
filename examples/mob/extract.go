package main

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
)

func extractFunction(goCode, functionName string) (string, error) {

	return "", nil
}

func main() {
	openai_apikey, err := tzapconnect.LoadOPENAI_APIKEY()
	if err != nil {
		panic(err)
	}
	goCode := "hey"
	tzap.
		NewWithConnector(
			tzapconnect.WithConfig(openai_apikey, config.Configuration{MD5Rewrites: true})).
		AddUserMessage("These are the changes you need to do:\n\n" + goCode).
		AddUserMessage("What are the relevant files to change based on this? \n\n" + goCode).
		RequestChatCompletion()
	//Rewrite File

}
