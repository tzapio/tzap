// package tzap provides a library to simplify manual workflows when dealing with chatgpt.
package tzap

import "github.com/tzapio/tzap/pkg/types"

// MutationTzap applies the provided function to the current Tzap object and returns a new Tzap object.
func (t *Tzap) MutationTzap(fn func(*Tzap) *Tzap) *Tzap {
	Log(t, "Mutation Tzap")
	return fn(t)
}

// WorkTzap executes the provided function and returns the Tzap object.
func (t *Tzap) WorkTzap(fn func(*Tzap)) *Tzap {
	Log(t, "Work Tzap")
	tb := t.CloneTzap(&Tzap{Name: "Work"})
	fn(tb)
	return t
}

// IsolatedTzap executes the provided function with a new isolated Tzap object.
// It does not modify the current Tzap object but returns it.
func (t *Tzap) IsolatedTzap(fn func(*Tzap)) *Tzap {
	Log(t, "Isolated Tzap")
	isolated := t.New()
	fn(isolated)
	return t
}

// Exit raises a panic to exit from the current Tzap.
func (t *Tzap) Exit() *Tzap {
	Log(t, "Exit Tzap")
	panic("parachute exit")
}

// Map iterates through the children Tzap objects and applies the provided function.
// Returns a new Tzap object containing the result of the function application.
func (t *Tzap) Map(fn func(*Tzap) *Tzap) *Tzap {
	children := t.Data["children"].([]*Tzap)
	Log(t, "START MAP", len(children))
	mappedChildren := make([]*Tzap, len(children))

	for i, child := range children {
		Log(t, "MAP child I", i)
		mappedChildren[i] = fn(child)
	}
	tzmapped := t.AddTzap(&Tzap{Name: "Map", Message: t.Message, Data: types.MappedInterface{"children": mappedChildren}})
	return tzmapped
}

func (t *Tzap) Reduce(fn func(*Tzap, *Tzap) *Tzap) *Tzap {
	children := t.Data["children"].([]*Tzap)
	Log(t, "START REDUCE", len(children))

	reducedTzap := t.AddTzap(&Tzap{Name: "Reduce"})

	for i, child := range children {
		Log(t, "REDUCE child I", i)
		reducedTzap = fn(reducedTzap, child)
	}

	return reducedTzap
}
func (t *Tzap) Accumulate(fn func(*Tzap) *Tzap) *Tzap {
	children := t.Data["children"].([]*Tzap)
	Log(t, "Accumulate start", len(children))
	tzmapped := t.AddTzap(&Tzap{Name: "Accumulate"})
	for _, child := range children {
		child.Parent = tzmapped
		tzmapped = tzmapped.ApplyWorkflowP(fn(child))
	}
	return tzmapped
}

func (t *Tzap) Each(fn func(*Tzap)) *Tzap {
	children := t.Data["children"].([]*Tzap)
	Log(t, "Each start", len(children))

	for _, child := range children {
		fn(child)
	}
	return t
}

// Recursive applies the provided function recursively to the Tzap object
// and its children.
func (t *Tzap) Recursive(tf func(tzapThatCreatesNewChildren *Tzap) *Tzap) *Tzap {
	recursive := func(t *Tzap) *Tzap {
		return t.Recursive(tf)
	}

	Log(t, "Recursive Function called")
	return t.ApplyWorkflowFN(tf). // should error check
					Map(func(t *Tzap) *Tzap {
			return t.ApplyWorkflowFN(recursive)
		})
}
