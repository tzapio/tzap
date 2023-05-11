package tzap

import (
	"errors"
	"fmt"
	"html"
	"os"
	"os/exec"
	"strings"
)

type GraphVizNode struct {
	Id       string
	Label    string
	Tooltip  string
	Style    string
	Children []*GraphVizNode
}

type GraphVizEdge struct {
	FromNode *GraphVizNode
	ToNode   *GraphVizNode
	Tooltip  string
	Style    string
}

type GraphVizSubgraph struct {
	Id        string
	Label     string
	Tooltip   string
	Nodes     []*GraphVizNode
	Edges     []*GraphVizEdge
	SubGraphs []*GraphVizSubgraph
}

type GraphVizGraph struct {
	Label     string
	Tooltip   string
	Nodes     []*GraphVizNode
	Edges     []*GraphVizEdge
	SubGraphs []*GraphVizSubgraph
}

func ConvertGraphvizToSVG(inputFile string, outputFile string) error {
	// Check if the input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return errors.New("input file does not exist")
	}

	// Run the dot command to convert the input to SVG
	cmd := exec.Command("dot", "-Tsvg", inputFile)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	// Write the SVG output to the output file
	err = os.WriteFile(outputFile, output, 0644)
	if err != nil {
		return err
	}

	return nil
}
func GenerateGraphvizDotFile(filename string, graph *GraphVizGraph) error {
	println("Saving to ", filename)
	var dotBuilder strings.Builder

	dotBuilder.WriteString("digraph G {\n")
	layout := getLayout()
	dotBuilder.WriteString(layout)

	for _, node := range graph.Nodes {
		dotBuilder.WriteString(nodeString(node))

		for _, child := range node.Children {
			dotBuilder.WriteString(fmt.Sprintf("\t\"%s\" -> \"%s\";\n", node.Id, child.Id))
		}
	}
	for _, edge := range graph.Edges {
		dotBuilder.WriteString(edgeString(edge))
	}

	for _, subgraph := range graph.SubGraphs {
		writeSubgraph(&dotBuilder, subgraph)
	}

	dotBuilder.WriteString("}\n")

	return os.WriteFile(filename, []byte(dotBuilder.String()), 0644)
}
func getLayout() string {
	return `
	graph [bgcolor="#222222", fontcolor="white", fontname="Arial", fontsize=10];
    node [shape=box, style=filled, fillcolor="#3a3a3a", fontcolor="white", fontname="Arial", fontsize=10, color="#888888"];
    edge [color="#ffffff", fontcolor="white", fontname="Arial", fontsize=10];
`
}
func nodeString(node *GraphVizNode) string {
	bracket := ""
	if node.Tooltip != "" {
		bracket += fmt.Sprintf(", tooltip=<%s> ", node.Tooltip)
	}

	if node.Style == "invis" {
		bracket += fmt.Sprintf(", style=%s, width = 5, height = 5", node.Style)
	}
	return fmt.Sprintf("\"%s\" [label=<%s> %s];\n", node.Id, node.Label, bracket)
}

func edgeString(edge *GraphVizEdge) string {
	if edge.Style != "" {
		return fmt.Sprintf("\"%s\" -> \"%s\" [style=%s];\n", edge.FromNode.Id, edge.ToNode.Id, edge.Style)
	}
	if edge.FromNode == nil {
		return "\n"
	}
	return fmt.Sprintf("\"%s\" -> \"%s\";\n", edge.FromNode.Id, edge.ToNode.Id)
}

func writeSubgraph(dotBuilder *strings.Builder, subgraph *GraphVizSubgraph) {
	dotBuilder.WriteString(fmt.Sprintf("\n\tsubgraph %s {\n", subgraph.Id))
	dotBuilder.WriteString(fmt.Sprintf("\t\tlabel = \"%s\";\n", subgraph.Label))
	dotBuilder.WriteString("\t\tbgcolor = \"#333333\";\n")

	for _, node := range subgraph.Nodes {
		dotBuilder.WriteString("\t\t" + nodeString(node))
	}

	for _, edge := range subgraph.Edges {
		dotBuilder.WriteString("\t\t" + edgeString(edge))
	}

	for _, childSubgraph := range subgraph.SubGraphs {
		writeSubgraph(dotBuilder, childSubgraph)
	}

	dotBuilder.WriteString("\t}\n")
}

func FillGraphVizGraph() *GraphVizGraph {
	graph := &GraphVizGraph{}
	tzapNodes := make(map[int]*GraphVizNode)
	for _, t := range GlobalTzaps {
		metadataLabel := generateGraphvizDotFile2(t)
		label := fmt.Sprintf("%s (%d) %s", t.Name, t.Id, metadataLabel)
		node := &GraphVizNode{Id: fmt.Sprintf("tzap_%d", t.Id), Label: label}
		tzapNodes[t.Id] = node
		graph.Nodes = append(graph.Nodes, node)

		if t.Parent != nil {
			edge := &GraphVizEdge{FromNode: tzapNodes[t.Parent.Id], ToNode: node}
			graph.Edges = append(graph.Edges, edge)
		}
	}

	for j, thread := range GlobalGraphVizLogThreads {
		messageNodes := make(map[int]*GraphVizNode)
		messageEdges := make([]*GraphVizEdge, 0)
		chatId := fmt.Sprintf("cluster_chat_%d", j)
		requestId := fmt.Sprintf("%s_REQUEST", chatId)
		responseId := fmt.Sprintf("%s_RESPONSE", chatId)

		requestSubgraph := &GraphVizSubgraph{Id: requestId, Label: "REQUEST"}

		responseSubgraph := &GraphVizSubgraph{Id: responseId, Label: "RESPONSE"}
		for i, msg := range thread.Messages {
			msgId := fmt.Sprintf("chat_%d_msg_%d", j, i)

			escapedContent := escapeContent(msg.Content)
			label := fmt.Sprintf("Message Tokens(%d) (%d):<br/>Role: %s<br/>%s", msg.TokenCount, i, msg.Role, htmlEscapeNewLine(escapedContent))
			tooltip := escapeNewLine(escapedContent)
			messageNode := &GraphVizNode{Id: msgId, Label: label, Tooltip: tooltip}
			messageNodes[i] = messageNode

			if msg.Direction == "REQUEST" {
				requestSubgraph.Nodes = append(requestSubgraph.Nodes, messageNode)
				edge := &GraphVizEdge{Style: "dotted", FromNode: tzapNodes[msg.TzapId], ToNode: messageNode}
				messageEdges = append(messageEdges, edge)
			} else if msg.Direction == "RESPONSE" {
				responseSubgraph.Nodes = append(responseSubgraph.Nodes, messageNode)
				edge := &GraphVizEdge{Style: "dotted", FromNode: messageNode, ToNode: tzapNodes[msg.TzapId]}
				messageEdges = append(messageEdges, edge)
			}

			if i > 0 {
				previousMessageNode, ok := messageNodes[i-1]
				if !ok {
					continue
				}
				edge := &GraphVizEdge{Style: "dotted", FromNode: previousMessageNode, ToNode: messageNode}
				messageEdges = append(messageEdges, edge)
			}
		}

		chatSubgraph := &GraphVizSubgraph{Id: chatId, Label: fmt.Sprintf("GPT Chat(%d):", j)}
		chatSubgraph.SubGraphs = append(chatSubgraph.SubGraphs, requestSubgraph)
		chatSubgraph.SubGraphs = append(chatSubgraph.SubGraphs, responseSubgraph)
		graph.SubGraphs = append(graph.SubGraphs, chatSubgraph)

		graph.Edges = append(graph.Edges, messageEdges...)
	}
	return graph
}

func generateGraphvizDotFile2(t *Tzap) string {
	truncMsg := fmt.Sprintf("%.30s", t.Message.Content)
	metadataLabel := ""
	filepath, ok := t.Data["filepath"].(string)
	if ok {
		metadataLabel += fmt.Sprintf("\n<b>File out:</b> %s", filepath)
	}
	if len(truncMsg) > 0 {
		metadataLabel += fmt.Sprintf("\nMessage:\nRole:%s\nContent:\n%s [...]", t.Message.Role, html.EscapeString(truncMsg))
	}
	return replaceNewLines(metadataLabel)
}

func escapeContent(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "<", "\\<")
	s = strings.ReplaceAll(s, ">", "\\>")
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "}", "")
	s = strings.ReplaceAll(s, "&", "&amp;")
	escaped := html.EscapeString(s)
	return escaped
}
func htmlEscapeNewLine(s string) string {
	truncated := fmt.Sprintf("%.300s", s)
	newLined := strings.ReplaceAll(truncated, "\n", "<br/>")
	return newLined
}
func escapeNewLine(s string) string {
	newLined := s //strings.ReplaceAll(s, "\n", "\\n")
	return newLined
}
