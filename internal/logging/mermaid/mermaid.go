package mermaid

import (
	"fmt"

	"github.com/tzapio/tzap/pkg/tzap"
)

func FillMermaidGraph(t *tzap.Tzap) *MermaidGraph {
	graph := &MermaidGraph{}
	rgetMessagesMermaid(t, graph)
	return graph
}

func rgetMessagesMermaid(t *tzap.Tzap, graph *MermaidGraph) {
	if t == nil {
		return
	}

	// Call the function recursively for the parent
	rgetMessagesMermaid(t.Parent, graph)

	// Add current Tzap to the graph as a node
	node := &MermaidNode{Id: fmt.Sprintf("tzap_%d", t.Id), Label: t.Name}
	graph.Nodes = append(graph.Nodes, node)

	// If the parent exists, add an edge between the parent and the current Tzap
	if t.Parent != nil {
		parentNode := &MermaidNode{Id: fmt.Sprintf("tzap_%d", t.Parent.Id), Label: t.Parent.Name}
		edge := &MermaidEdge{FromNode: parentNode, ToNode: node}
		graph.Edges = append(graph.Edges, edge)
	}
}
