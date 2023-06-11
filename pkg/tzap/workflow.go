package tzap

import "github.com/tzapio/tzap/pkg/types"

// ApplyWorkflowP applies a given workflow Tzap instance to the current Tzap instance.
// Returns the applied workflow with its Parent set to the current Tzap instance.
func (t *Tzap) ApplyWorkflowP(workflow *Tzap) *Tzap {
	at := t.CloneTzap(&Tzap{Name: "ApplyWorkflowS"})
	Log(t, "Applying workflow")
	workflow.Parent = at
	return workflow
}

// ApplyWorkflowFN applies a function that takes a Tzap instance and returns a modified Tzap instance.
// Returns the result of the given function applied to the current Tzap instance.
func (t *Tzap) ApplyWorkflowFN(nt func(*Tzap) *Tzap) *Tzap {
	Log(t, "Applying workflow FN")
	return nt(t.CloneTzap(&Tzap{Name: "ApplyWorkflow"}))
}

// WARNING: ApplyWorkflow clones messages from previous Tzap instances. This duplicates the message.
func (t *Tzap) ApplyWorkflow(nt types.NamedWorkflow[*Tzap, *Tzap]) *Tzap {
	Log(t, "Applying workflow")
	workflowResult := nt.Workflow(t.CloneTzap(&Tzap{Name: "ApplyWorkflow (" + nt.Name + ") Start"}))
	endWorkflow := workflowResult.CloneTzap(&Tzap{Name: "ApplyWorkflow (" + nt.Name + ") End"})
	return endWorkflow
}

func (t *Tzap) ApplyErrorWorkflow(nt types.NamedWorkflow[*Tzap, *ErrorTzap], fn func(*ErrorTzap) error) *Tzap {
	et := nt.Workflow(t.CloneTzap(&Tzap{Name: "ApplyErrorWorkflow (" + nt.Name + ")"}))
	err := fn(et)
	if err != nil {
		panic(err)
	}

	return et.Tzap
}
