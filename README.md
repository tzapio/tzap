# Tzap

Tzap is a library for Prompts as Code. It provides a toolkit to build, customize, and extend chatbot prompts in a streamlined and extensible manner. The library is designed to make it easy for developers to create reusable Tzap instances and combinations of Tzaps to quickly and effectively implement desired outcomes in their chatbot-based applications.

## Features

- Create reusable Tzap instances and templates
- Apply templates and functions to existing Tzaps
- Manipulate file paths and directories
- Fetch chat responses and generate content using OpenAI's GPT-4 model
- Provide chat message context in Golang
- Manage Tzap instances with parent and child relationships
- Utility methods to read, write and manipulate files

## Example Usage

Here's an example of how you can use Tzap to read a list of Go files and output a README.md: 

```go
package main

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"
	"github.com/tzapio/tzap/templates/code/documents"
)

func main() {
	tzap.NewWithConnector(tzapconnect.WithConfig(
		config.Configuration{
			AutoMode:    true,
			OpenAIModel: openai.GPT4,
			MD5Rewrites: true,
		})).
		ApplyTemplateFN(documents.ReadmeGithub(
			"Tzap is a library for Prompts as Code.",
			[]string{
				"pkg/types/interfaces.go",
				"pkg/types/structs.go",
				"pkg/tzap/templates.go",
				"pkg/tzap/file.go",
				"pkg/tzap/fetch-chat.go",
				"pkg/tzap/tzap.go",
				"examples/githubdoc/main.go",
				"examples/refactoring/main.go",
			},
			"README.md",
			"",
		))
}
```

Another example uses Tzap to refactor a Go file, improving its readability and adding comments to functions: 

```go
package main

import (
	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types/openai"
	"github.com/tzapio/tzap/pkg/tzap"
	"github.com/tzapio/tzap/pkg/tzapconnect"

	"github.com/tzapio/tzap/pkg/util"
	"github.com/tzapio/tzap/templates/code/codegeneration"
)

func main() {
	tzap.NewWithConnector(tzapconnect.WithConfig(
		config.Configuration{
			AutoMode:    true,
			OpenAIModel: openai.GPT4,
			MD5Rewrites: true,
		})).
		LoadFileDir("/workspaces/goman/tzaps", "*.go").
		Map(func(t *tzap.Tzap) *tzap.Tzap {
			return t.
				ApplyTemplateFN(
					codegeneration.MakeCodeGO(`
You are helping the user writing a library for chatgpt prompting. You primarely write Golang. Most files already exists. Do not create new data structures.
### Current interface: 
//file: /workspaces/goman/tzaps/interfaces.go
`+util.ReadFileP("/workspaces/goman/tzaps/interfaces.go")+`
### General Tzap logic
//file: /workspaces/goman/tzaps/tzap.go
`+util.ReadFileP("/workspaces/goman/tzaps/tzap.go")+`
### Additional types
//file: /workspaces/goman/tzaps/msg/message.go
`+util.ReadFileP("/workspaces/goman/tzaps/msg/message.go"),
						//	"Make unit tests. Use testify go. If needed create tmp files. Use package tzap_test. Use testnames Test_{function}_{givenCamelCase}_{expectCamelCase}."),
						"Analyze what can be improved. Refactor the following file to be more readable. Make comments for the functions. Do not add any new public functions, only rewrite."),
				)
		})
}
```

## Package Structure

The following files and packages make up the Tzap library:

### `pkg/types/interfaces.go`

This file defines the `ITzap` interface that Tzap instances should implement. It includes methods to manage Tzap instances, load files and directories, request chatbot content, and apply templates.

### `pkg/types/structs.go`

This file defines the `Message` and `MappedInterface` structures. The `Message` structure holds information related to individual chatbot messages, while `MappedInterface` is a type alias for a map of string keys with interface values.

### `pkg/tzap/templates.go`

The `templates.go` file contains the `ApplyTemplate` and `ApplyTemplateFN` methods, which allow you to apply a pre-defined template Tzap instance (`ApplyTemplate`) or a custom function that takes a Tzap instance and returns a modified Tzap instance (`ApplyTemplateFN`).

### `pkg/tzap/file.go`

The `file.go` file contains various utility methods to handle Tzap instances for file operations. These include loading a file or directory of files, preparing a file for output, and handling the case where a file does not exist and new content needs to be requested.

### `pkg/tzap/fetch-chat.go`

This file contains methods to send chatbot requests and fetch the chatbot-generated content. These methods mainly work with OpenAI's GPT models, using the Tzap generation API to generate chatbot responses based on given prompt and context data.

### `pkg/tzap/tzap.go`

This file defines the core logic of the Tzap library. The `Tzap` struct is defined in this file with its fields, methods, and interfaces. The struct contains methods to manage Tzap instances, apply templates, and manipulate files.

### Example files

`examples/githubdoc/main.go` contains an example of how to use the Tzap library to generate a README.md file for a GitHub project, based on an input list of Go files.

`examples/refactoring/main.go` contains a sample usage that demonstrates how to refactor a Go file using the Tzap library, with an emphasis on improving code readability and adding comments to functions.

## Usage

To get started with Tzap, create a new Tzap instance using the `NewWithConnector` method:

```go
tzapInstance := tzap.NewWithConnector(tzapconnect.WithConfig(settings))
```

Then, chain various methods to load and manipulate data, or apply templates and functions. The rendered content can be fetched and used as desired.