package action

import (
	"encoding/json"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

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
			"fileout": {
				Type:        jsonschema.String,
				Description: "The full file path to use as output",
			},
			"code": {
				Type:        jsonschema.String,
				Description: "The code change.",
			},
		},
		Required: []string{"suggestion", "plan", "filein", "code"},
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
var refactors openai.FunctionDefinition = openai.FunctionDefinition{
	Name:        "refactors",
	Description: "Instructions for refactoring several refactors or creating code.",
	Parameters: jsonschema.Definition{
		Type: jsonschema.Object,
		Properties: map[string]jsonschema.Definition{
			"refactors": {
				Type: jsonschema.Array,
				Items: &jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"filein": {
							Type:        jsonschema.String,
							Description: "The full file path to use as input",
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
							Description: "{the format of fileout, e.g. 'go code', 'github md', 'mermaid js'}",
						},
						"inspirationfiles": {
							Type:        jsonschema.Array,
							Description: "A list of files that can be used as inspiration. Eg the files containing types related to the plan",
							Items: &jsonschema.Definition{
								Type:        jsonschema.String,
								Description: "File that can be used as inspiration. Eg the file containing types related to the plan",
							},
						},
					},

					Required: []string{"filein", "mission", "task", "plan", "outputformat", "inspirationfiles"},
				},
			},
		},
		Required: []string{"refactors"},
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
				Description: "The code file to use as input",
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
		Required: []string{"filein", "fileout", "mission", "task", "plan", "outputformat", "example"},
	},
}
var implement openai.FunctionDefinition = openai.FunctionDefinition{
	Name:        "implement",
	Description: "Handles implement changes.",
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
			"changes": {
				Type: jsonschema.Array,
				Items: &jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"filein": {
							Type:        jsonschema.String,
							Description: "The full file path to use as input",
						},
						"fileout": {
							Type:        jsonschema.String,
							Description: "The full file path to use as output",
						},
						"task": {
							Type:        jsonschema.String,
							Description: "Explain what should be done",
						},
						"code": {
							Type:        jsonschema.String,
							Description: "The code change",
						},
						"inspirationfiles": {
							Type:        jsonschema.Array,
							Description: "A list of files that can be used as inspiration. Eg the files containing types related to the plan",
							Items: &jsonschema.Definition{
								Type:        jsonschema.String,
								Description: "File that can be used as inspiration. Eg the file containing types related to the plan",
							},
						},
					},

					Required: []string{"filein", "task", "code", "inspirationfiles"},
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
	test,
	implement,
}
var Tfs string = toString(tzapFunctions)

func toString(o []openai.FunctionDefinition) string {
	bytes, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
