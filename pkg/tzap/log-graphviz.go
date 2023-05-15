package tzap

import (
	"strings"
)

var GlobalTzaps []*Tzap

func HandleShutdown() {

	GenerateGraphvizDotFile("out/tzap2.dot", FillGraphVizGraph())
	Flush()
}

type GraphVizLogMessage struct {
	Role       string
	Content    string
	TokenCount int
	TzapId     int
	Direction  string
}
type GraphVizLogMessages struct {
	Messages   []GraphVizLogMessage
	TokenCount int
}

var GlobalGraphVizLogThreads []GraphVizLogMessages

func replaceNewLines(s string) string {
	return strings.ReplaceAll(s, "\n", "<br/>")
}

func getMessagesGraphViz(t *Tzap) {
	messages, count := rgetMessagesGraphViz(t.Parent)
	if t.InitialSystemContent != "" {
		c, err := t.TG.CountTokens(t.C, t.Message.Content)
		if err != nil {
			println("WARNING: could not count tokens", err.Error())
		}
		messages = append([]GraphVizLogMessage{{
			Role:      "system",
			Content:   t.InitialSystemContent,
			TzapId:    t.Id,
			Direction: "REQUEST",
		}}, messages...)
		count += c

	}
	c, err := t.TG.CountTokens(t.C, t.Message.Content)
	if err != nil {
		println("WARNING: could not count tokens", err.Error())
	}
	messages = append(messages, GraphVizLogMessage{
		Role:       t.Message.Role,
		Content:    t.Message.Content,
		TzapId:     t.Id,
		Direction:  "RESPONSE",
		TokenCount: c,
	})
	count += c
	GlobalGraphVizLogThreads = append(GlobalGraphVizLogThreads, GraphVizLogMessages{
		Messages:   messages,
		TokenCount: count,
	})

}
func rgetMessagesGraphViz(t *Tzap) ([]GraphVizLogMessage, int) {
	var messages []GraphVizLogMessage
	var count = 0
	if t.Parent != nil {
		m, c := rgetMessagesGraphViz(t.Parent)
		count += c
		messages = m
	}

	if t.Message.Content == "" || t.Message.Role == "" {
		return messages, count
	}
	key, ok := t.Data["memory"].(string)
	if ok && key != "" {
		mV := Mem[key]
		if mV.Content != "" {
			c, err := t.TG.CountTokens(t.C, mV.Content)
			if err != nil {
				println("WARNING: could not count tokens", err.Error())
			}
			message := GraphVizLogMessage{
				Role:       mV.Role,
				Content:    mV.Content,
				TokenCount: c,
			}
			count += c
			messages = append(messages, message)
		}
	}
	c, err := t.TG.CountTokens(t.C, t.Message.Content)
	if err != nil {
		println("WARNING: could not count tokens", err.Error())
	}
	count += c
	return append(messages, GraphVizLogMessage{Direction: "REQUEST", Role: t.Message.Role, Content: t.Message.Content, TzapId: t.Id, TokenCount: count}), count
}
