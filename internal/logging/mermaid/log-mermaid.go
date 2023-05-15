package mermaid

import (
	"fmt"
	"os"
	"strings"
)

func GenerateMermaidMarkupFile(filename string, graph *MermaidGraph) error {
	markup := createMermaidMarkup(graph)
	return os.WriteFile(filename, []byte(markup), 0644)
}

func createMermaidMarkup(graph *MermaidGraph) string {
	// Implement logic to generate mermaid.js markup from the graph
	var markupBuilder strings.Builder

	markupBuilder.WriteString("graph TD\n")

	// Add nodes
	for _, node := range graph.Nodes {
		markupBuilder.WriteString(fmt.Sprintf("%s[%s]\n", node.Id, node.Label))
	}

	// Add edges
	for _, edge := range graph.Edges {
		markupBuilder.WriteString(fmt.Sprintf("%s-->%s\n", edge.FromNode.Id, edge.ToNode.Id))
	}

	// Add subgraphs
	for _, subgraph := range graph.SubGraphs {
		writeSubgraphMermaid(&markupBuilder, subgraph)
	}

	return markupBuilder.String()
}

func writeSubgraphMermaid(markupBuilder *strings.Builder, subgraph *MermaidSubgraph) {
	// Implement logic to write a subgraph to the markup builder
	markupBuilder.WriteString(fmt.Sprintf("subgraph %s[%s]\n", subgraph.Id, subgraph.Label))

	for _, node := range subgraph.Nodes {
		markupBuilder.WriteString(fmt.Sprintf("  %s[%s]\n", node.Id, node.Label))
	}

	for _, edge := range subgraph.Edges {
		markupBuilder.WriteString(fmt.Sprintf("  %s-->%s\n", edge.FromNode.Id, edge.ToNode.Id))
	}

	for _, childSubgraph := range subgraph.SubGraphs {
		writeSubgraphMermaid(markupBuilder, childSubgraph)
	}

	markupBuilder.WriteString("end\n")
}
