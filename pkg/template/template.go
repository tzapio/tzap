package template

import (
	"bytes"
	"encoding/json"
	"text/template"
)

type TriggerEvent struct {
	Name string
	Data interface{}
}

type WorkflowStep struct {
	Name        string
	TemplateSrc string
	Template    *template.Template
	InputParams string
}

type StepExecutionResult struct {
	OutputValues map[string]interface{}
	OutputString string
}

func NewWorkflowStep(name string, templateSrc string) *WorkflowStep {
	step := &WorkflowStep{
		Name:        name,
		TemplateSrc: templateSrc,
	}

	err := step.initialize()
	if err != nil {
		panic(err)
	}

	return step
}

func (step *WorkflowStep) initialize() error {
	tmpl, err := template.New("").Delims("{{", "}}").Parse(step.TemplateSrc)
	if err != nil {
		return err
	}

	step.Template = tmpl
	return nil
}

func (step *WorkflowStep) Execute(input map[string]interface{}) (string, error) {
	var output bytes.Buffer
	if err := step.Template.Execute(&output, input); err != nil {
		return "", err
	}

	return output.String(), nil
}

func parseOutputValues(output string) map[string]interface{} {
	values := make(map[string]interface{})
	json.Unmarshal([]byte(output), &values)
	return values
}
