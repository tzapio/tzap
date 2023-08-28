package tzap

import (
	"context"

	"github.com/tzapio/tzap/pkg/config"
	"github.com/tzapio/tzap/pkg/types"
)

var count = 1

func addId(t *Tzap) {
	if t.Id > 0 {
		println("Tzap already has an id", t.Id, t.Name)
		panic(t)
	}
	t.Id = count
	count += 1
}
func (t *Tzap) HandleShutdown() {
	Flush()
}

// Tzap is a structure that holds data and methods related to Tzap objects.
type Tzap struct {
	Id                   int
	Name                 string
	InitialSystemContent string
	Message              types.Message
	Data                 types.MappedInterface `json:"-"`
	C                    context.Context       `json:"-"`
	TG                   types.TGenerator      `json:"-"`

	types.ITzap[*Tzap, any] `json:"-"`

	Parent *Tzap
}

// NewTzap creates a new Tzap with default values, and returns its pointer.
// Mainly for mocking purposes. Does not have a connector, will likely crash.
func InternalNew() *Tzap {
	t := &Tzap{
		Name:    "ConnectionLess",
		Message: types.Message{},
		Data:    types.MappedInterface{},
		C:       context.Background(),
	}
	addId(t)
	return t
}

// NewTzap creates a new Tzap with default values, and returns its pointer.
// Mainly for mocking purposes. Does not have a connector, will likely crash.
func InjectNew(tg types.TGenerator, conf config.Configuration) *Tzap {
	t := &Tzap{
		Name:    "Connection",
		Message: types.Message{},
		Data:    types.MappedInterface{},
		C:       config.NewContext(context.Background(), conf),
		TG:      tg,
	}
	addId(t)
	return t
}
func NewWithConnector(connector types.TzapConnector) *Tzap {
	tg, conf := connector()
	return InjectNew(tg, conf)
}

// CopyConnection returns a new Tzap with default values.
func (t *Tzap) CopyConnection() *Tzap {
	tc := &Tzap{
		Name:    "NewConnection",
		Message: types.Message{},
		Data:    types.MappedInterface{},
		C:       t.C,
		TG:      t.TG,
	}
	addId(tc)
	return tc
}

// AppendParentContext assigns the parent's context to the Tzap object, if present.
func (t *Tzap) appendParentContext() *Tzap {
	if t.Parent != nil {
		t.C = t.Parent.C
		t.TG = t.Parent.TG
	}
	return t
}

// onNewTzap is a helper method that appends the parent's context to the Tzap object.
func (t *Tzap) onNewTzap() *Tzap {
	addId(t)
	return t.appendParentContext()
}

// AddContextChange replaces the current Tzap's context with the provided context.
func (t *Tzap) AddContextChange(fn func(context.Context) context.Context) *Tzap {
	newTzap := t.AddTzap(&Tzap{Name: "MutateContext"})
	newTzap.C = fn(newTzap.C)
	return newTzap
}

// AddTzap (mostly internal use) initializes and adds a new Tzap child to the current Tzap object.
func (t *Tzap) AddTzap(tc *Tzap) *Tzap {
	Logf(t, "Add tzap (%s)", tc.Name)
	tc.Parent = t
	return tc.onNewTzap()
}

// CloneTzap (mostly internal use) clones a Tzap object and assigns values based on the provided Tzap object.
func (previousTzap *Tzap) CloneTzap(suggestedTzap *Tzap) *Tzap {
	Logf(previousTzap, "Clone tzap (%s)", suggestedTzap.Name)
	baseTzap := &Tzap{
		Parent:               previousTzap,
		Name:                 previousTzap.Name,
		InitialSystemContent: previousTzap.InitialSystemContent,
		Data:                 previousTzap.Data,
	}

	if suggestedTzap.Parent != nil {
		baseTzap.Parent = suggestedTzap.Parent
	}
	if suggestedTzap.Name != "" {
		baseTzap.Name = suggestedTzap.Name
	}
	if suggestedTzap.InitialSystemContent != "" {
		baseTzap.InitialSystemContent = suggestedTzap.InitialSystemContent
	}
	if suggestedTzap.Message.Role != "" {
		baseTzap.Message.Role = suggestedTzap.Message.Role
	}
	if suggestedTzap.Message.Content != "" {
		baseTzap.Message.Content = suggestedTzap.Message.Content
	}

	if len(suggestedTzap.Data) > 0 {
		baseTzap.Data = suggestedTzap.Data
	}

	return baseTzap.onNewTzap()
}

// HijackTzap (mostly internal use) effectively de-attaches from previous Tzap by changing the own parent to parents parent.
// This can be used AddUserMessage("H").LoadTaskOrRequestNewTask().Hijack() .() Tzap replaces the current Tzap's context and parent with the provided Tzap's context and parent.
func (previousTzap *Tzap) HijackTzap(bypassWith *Tzap) *Tzap {
	Logf(previousTzap, "Hijack tzap (%s)", bypassWith.Name)
	bypassWith.Parent = previousTzap.Parent
	return bypassWith.onNewTzap()
}
