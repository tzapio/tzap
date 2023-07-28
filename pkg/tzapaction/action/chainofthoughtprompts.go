package action

import (
	"bytes"
	"fmt"
	"text/template"
)

type Step struct {
	Task         string
	Format       string
	OutputFormat string
}

type ChainOfThought struct {
	Steps []Step
}

func FindChainOfThoughtPrompt() string {
	tmpl := template.Must(template.New("findchain").Funcs(template.FuncMap{"inc": func(i int) int {
		return i + 1
	}}).Parse(`Follow these steps to look for files for the user.
The user will write a prompt and your task is to find the most relevant file they might be interested in.
The user prompt is not directly related to the content found.
You will answer briefly but each detail will be discussed. 

Make sure to include #### to separate every step. Do not mention the step descriptions or formats.
{{range $i, $step := .Steps}}
Step {{inc $i}}:#### {{$step.Task}}. Format: {{$step.Format}}
{{end}}

Use the following output format:
{{range $i, $step := .Steps}}
Step {{inc $i}}:####<step {{$step.OutputFormat}}>
{{end}}`))

	steps := []Step{
		{
			Task:         "What the user is asking for",
			Format:       `{"query":"<query>"}`,
			OutputFormat: "query json",
		},
		{
			Task:         "Walk through the results and reason how the found results might be relevant",
			Format:       `{"reasoning":"<reasoning hypothesis>"}`,
			OutputFormat: "reasoning json",
		},
		{
			Task:         "Walk through each file that exists and explain each of their relevance. Then evaluate the file relevance based on the query",
			Format:       `[{"filepath":"<filepath>","relevance":"<how the file is relevance>","evaluation":"<evalutation of how it answers query>"},...]`,
			OutputFormat: "relevance and evaluation json",
		},
		{
			Task:         "Write a reasoning followed by a score 0-100 based on relevance, higher is better",
			Format:       `[{"filepath":"<filepath>","reason":"<reason>", "score":"<score>"},...]`,
			OutputFormat: "scores json",
		},
		{
			Task:         "Criticize your answer so far",
			Format:       `{"critique":"<critique>","correction":[{"filepath":"<filepath>","reason":"<reason>", "score":"<score>"},...]}`,
			OutputFormat: "critique json",
		},
		{
			Task:         "Order files based on score and only include relevant files",
			Format:       `["<filepath>",...]`,
			OutputFormat: "file order json",
		},
	}
	msg := ChainOfThought{Steps: steps}

	var buf bytes.Buffer
	err := tmpl.Execute(&buf, msg)
	if err != nil {
		fmt.Println(err)
	}

	// Get the template output as a string.
	stringOutput := buf.String()
	return stringOutput
}

var FindMessage = `Follow these steps to look for files for the user.
The user will write a prompt and your task is to find the most relevant file they might be interested in.
The user prompt is not directly related to the content found.
You will answer briefly but each detail will be discussed. 

Make sure to include #### to separate every step. Do not mention the step descriptions or formats.

Step 1:#### Clarify what the user is searching for. Format: {"clarification":"<clarification>"}
Step 2:#### Walk through the results and reason how the found results might be relevant. Format: {"reasoning":"<reasoning hypothesis>"}
Step 3:#### Walk through each file that exists and explain each of their relevance. Format: [{"filepath":"<filepath>","relevance":"relevance"},...]
Step 4:#### Evaluate the file relevance. Format: [{"filepath":"<filepath>","evaluation":"<evaluation>"},...]
Step 5:#### Write a reasoning followed by a score 0-100 based on relevance, higher is better. Format: [{"filepath":"<filepath>","reason":"<reason>", "score":"<score>"},...]
Step 6:#### Critize your answer so far. {"critique":"<critique>","correction":[{"filepath":"<filepath>","reason":"<reason>", "score":"<score>"},...]}
Step 7:#### Order files based on score. Format: ["<filepath>",...]
Step 8:#### Is there a definitive match? Format: {"matching":"<single, multiple, none>"}
Step 9:#### State the filepaths in JSON. Format: ["<filepath>",...]

Use the following format:
Step 1:####<step 1 clarification json>
Step 2:####<step 2 reasoning json >
Step 3:####<step 3 explanation json>
Step 4:####<step 4 evaluate json>
Step 5:####<step 5 scores json>
Step 6:####<step 6 critique>
Step 7:####<step 7 file order json>
Step 8:####<step 8 assessment json>
Step 9:####<step 9 filepaths json>
`
var COTTemplate = `{{ }}`
