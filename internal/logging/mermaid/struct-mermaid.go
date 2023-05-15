package mermaid

type MermaidNode struct {
	Id    string
	Label string
}

type MermaidEdge struct {
	FromNode *MermaidNode
	ToNode   *MermaidNode
}

type MermaidSubgraph struct {
	Id        string
	Label     string
	Nodes     []*MermaidNode
	Edges     []*MermaidEdge
	SubGraphs []*MermaidSubgraph
}

type MermaidGraph struct {
	Nodes     []*MermaidNode
	Edges     []*MermaidEdge
	SubGraphs []*MermaidSubgraph
}
