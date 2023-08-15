package action

import (
	"encoding/json"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

var add openai.FunctionDefinition = openai.FunctionDefinition{
	Name:        "add",
	Description: "Add a code to a new file ",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"mission": {
				Type:        jsonschema.String,
				Description: "The mission to solve",
			},
			"plan": {
				Type:        jsonschema.String,
				Description: "The plan to solve the suggestion",
			},
			"fileout": {
				Type:        jsonschema.String,
				Description: "The path to the file to add",
			},
			"code": {
				Type:        jsonschema.String,
				Description: "The change to add",
			},
			"filein": {
				Type:        jsonschema.String,
				Description: "The file to use as input.",
			},
			"outputformat": {
				Type:        jsonschema.String,
				Description: "The output format to use, e.g. golang code, github md, mermaid js",
			},
			"inspirationfiles": {
				Type: jsonschema.Array,
				Items: &jsonschema.Definition{
					Type:        jsonschema.String,
					Description: "A list of files that can be used as inspiration for the refactoring. Eg the files containing types related to the refactoring",
				}},
		},
		Required: []string{"mission", "plan", "fileout"},
	},
}
var edit openai.FunctionDefinition = openai.FunctionDefinition{
	Name:        "edit",
	Description: "Edit the code of a file",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"mission": {
				Type:        jsonschema.String,
				Description: "The mission to solve",
			},
			"plan": {
				Type:        jsonschema.String,
				Description: "The plan to solve the suggestion",
			},
			"task": {
				Type:        jsonschema.String,
				Description: "The current task",
			},
			"filein": {
				Type:        jsonschema.String,
				Description: "The path to the file to edit",
			},
			"code": {
				Type:        jsonschema.String,
				Description: "The change to apply",
			},
			"outputformat": {
				Type:        jsonschema.String,
				Description: "The output format to use, e.g. golang code, github md, mermaid js",
			},
			"inspirationfiles": {
				Type: jsonschema.Array,
				Items: &jsonschema.Definition{
					Type:        jsonschema.String,
					Description: "A list of files that can be used as inspiration for the refactoring. Eg the files containing types related to the refactoring",
				}},
		},
		Required: []string{"mission", "plan", "filein"},
	},
}

var refactor openai.FunctionDefinition = openai.FunctionDefinition{
	Name:        "refactor",
	Description: "Creates a new file or modifies existing file to apply refactor changes.",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"filein": {
				Type:        jsonschema.String,
				Description: "The file to use as input.",
			},
			"fileout": {
				Type:        jsonschema.String,
				Description: "The output file that will be created or modified, usually the same file as input",
			},
			"mission": {
				Type:        jsonschema.String,
				Description: "The main mission describing what the users is aiming to achieve with the refactoring",
			},
			"task": {
				Type:        jsonschema.String,
				Description: "Explain what should be done",
			},
			"plan": {
				Type:        jsonschema.String,
				Description: "Explain how to achieve this",
			},
			"outputformat": {
				Type:        jsonschema.String,
				Description: "The output format to use, e.g. golang code, github md, mermaid js",
			},
			"inspirationfiles": {
				Type: jsonschema.Array,
				Items: &jsonschema.Definition{
					Type:        jsonschema.String,
					Description: "A list of files that can be used as inspiration for the refactoring. Eg the files containing types related to the refactoring",
				}},
		},
		Required: []string{"filein", "fileout", "mission", "task", "plan", "outputformat", "inspirationfiles"},
	},
}

var documentation openai.FunctionDefinition = openai.FunctionDefinition{
	Name:        "documentation",
	Description: "Creates a new documentation file or modifies documentation file for a specific code file.",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"filein": {
				Type:        jsonschema.String,
				Description: "The file to use as input.",
			},
			"fileout": {
				Type:        jsonschema.String,
				Description: "The documentation output file that will be created or modified.",
			},
			"mission": {
				Type:        jsonschema.String,
				Description: "Example: Document code in markdown format",
			},
			"task": {
				Type:        jsonschema.String,
				Description: "Example: Add documentation to the following file. Add comments to the functions. Add a markdown file with the same name as the file and add the documentation to it.",
			},
			"plan": {
				Type:        jsonschema.String,
				Description: "Example: Do not write any text because this file will be saved directly to {fileOut}",
			},
			"outputformat": {
				Type:        jsonschema.String,
				Description: "Example: {markdown documentation}",
			},
			"inspirationfiles": {
				Type: jsonschema.Array,
				Items: &jsonschema.Definition{
					Type:        jsonschema.String,
					Description: "A list of files that can be used as inspiration for the refactoring. Eg the files containing types related to the refactoring",
				}},
		},
		Required: []string{"filein", "fileout", "mission", "task", "plan", "outputformat"},
	},
}
var test openai.FunctionDefinition = openai.FunctionDefinition{
	Name:        "test",
	Description: "Creates a new test file or modifies test file for a specific code file.",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"filein": {
				Type:        jsonschema.String,
				Description: "The code file to use as input",
			},
			"fileout": {
				Type:        jsonschema.String,
				Description: "The test file that will be created or modified.",
			},
			"mission": {
				Type:        jsonschema.String,
				Description: "Example: Write a unit test",
			},
			"task": {
				Type:        jsonschema.String,
				Description: "Example: Add unit tests.",
			},
			"plan": {
				Type:        jsonschema.String,
				Description: "Example: Do not write any text because this file will be saved directly to {fileOut}",
			},
			"outputformat": {
				Type:        jsonschema.String,
				Description: "Example: {the specific language unit test description}",
			},
			"inspirationfiles": {
				Type: jsonschema.Array,
				Items: &jsonschema.Definition{
					Type:        jsonschema.String,
					Description: "A list of files that can be used as inspiration for the refactoring. Eg the files containing types related to the refactoring",
				}},
		},
		Required: []string{"filein", "fileout", "mission", "task", "plan", "outputformat"},
	},
}
var implement openai.FunctionDefinition = openai.FunctionDefinition{
	Name:        "code",
	Description: "Add or edit code in multiple files. Thoroughly plan the code changes to be made and then implement them.",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"mission": {
				Type:        jsonschema.String,
				Description: "The main mission describing what the users is aiming to achieve",
			},
			"plan": {
				Type:        jsonschema.String,
				Description: "Explain how to achieve this",
			},
			"tasks": {
				Type:        jsonschema.Array,
				Description: "The list of implementation details of the additions or changes to be made. Each file should only exist",
				Items: &jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"filein": {
							Type:        jsonschema.String,
							Description: "The full file path to use as input",
						},
						"fileout": {
							Type:        jsonschema.String,
							Description: "The full output file path",
						},
						"task": {
							Type:        jsonschema.String,
							Description: "A description of all changes applicable for this file.",
						},
						"inspirationfiles": {
							Type:        jsonschema.Array,
							Description: "A list of files that can be used as inspiration or contain specific information needed to do implement changs..",
							Items: &jsonschema.Definition{
								Type:        jsonschema.String,
								Description: "File that can be used as inspiration.",
							},
						},
					},

					Required: []string{"fileout", "task", "inspirationfiles"},
				},
			},
		},
		Required: []string{"changes", "mission", "plan"},
	},
}
var tzapFunctions []openai.FunctionDefinition = []openai.FunctionDefinition{
	refactor,
	documentation,
	edit,
	add,
	test,
	//implement,
}
var Tfs string = toString(tzapFunctions)

func toString(o []openai.FunctionDefinition) string {
	bytes, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
