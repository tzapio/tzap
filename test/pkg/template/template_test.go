package template_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/tzapio/tzap/pkg/template"
)

func TestRenderTemplate(t *testing.T) {
	testCases := []struct {
		name     string
		template string
		data     interface{}
		expected string
	}{
		{
			name:     "Render with strings",
			template: "{{.Greeting}} {{.Name}}",
			data:     map[string]string{"Greeting": "Hello", "Name": "World"},
			expected: "Hello World",
		},
		{
			name:     "Render with numbers",
			template: "{{.Num1}} {{.Num2}}",
			data:     map[string]int{"Num1": 42, "Num2": 13},
			expected: "42 13",
		},
		{
			name:     "Render with slices",
			template: "{{range .}}{{.}} {{end}}",
			data:     []string{"Hello", "World"},
			expected: "Hello World ",
		},
		{
			name:     "Render with maps",
			template: "{{range $key, $value := .}}{{ $key }}={{ $value }}, {{end}}",
			data:     map[string]string{"a": "1", "b": "2"},
			expected: "a=1, b=2, ",
		},
		{
			name:     "Render with structs",
			template: "{{.Name}} {{.Age}}",
			data: struct {
				Name string
				Age  int
			}{"Bob", 42},
			expected: "Bob 42",
		},
		{
			name:     "Render with nested fields",
			template: "{{.Name}} {{.City}}, {{.Address.Street}}",
			data: struct {
				Name    string
				Address struct {
					Street string
					Number int
				}
				City string
			}{"Bob", struct {
				Street string
				Number int
			}{"Main St", 12}, "New York"},
			expected: "Bob New York, Main St",
		},
		{
			name:     "Render with custom functions",
			template: "{{upperFirst .Name}}",
			data:     struct{ Name string }{"bob"},
			expected: "Bob",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var tpl bytes.Buffer
			tpl.WriteString(tc.template)

			data := map[string]interface{}{"Greeting": "Hello", "Name": "World"}
			workflowStep := template.WorkflowStep{TemplateSrc: tpl.String()}
			workflowStep.InputParams = data
			err := template.InitializeWorkflowStep(&workflowStep)
			if err != nil {
				panic(err)
			}

			result, err := template.ExecuteWorkflowStep(&workflowStep, &template.TriggerEvent{Name: "Test", Data: data})
			if err != nil {
				panic(err)
			}
			fmt.Println(result)
		})
	}
}
