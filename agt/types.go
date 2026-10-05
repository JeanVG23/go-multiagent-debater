package agt

type MessageType int

const (
	GiveFloor MessageType = iota
	Argument
	Summarize
	Verdict
	Stop
)

type Message struct {
	From    string
	Type    MessageType
	Content string
	Reply   chan<- Message
}
