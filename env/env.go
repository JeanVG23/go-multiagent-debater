package env

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Argument struct {
	Author    string
	Role      string
	Content   string
	Round     int
	Timestamp time.Time
}

type Environment struct {
	mu        sync.RWMutex
	topic     string
	arguments []Argument
}

func NewEnvironment(topic string) *Environment {
	return &Environment{
		topic:     topic,
		arguments: make([]Argument, 0),
	}
}

func (e *Environment) Topic() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.topic
}

func (e *Environment) AddArgument(author, role, content string, round int) Argument {
	e.mu.Lock()
	defer e.mu.Unlock()

	arg := Argument{
		Author:    author,
		Role:      role,
		Content:   strings.TrimSpace(content),
		Round:     round,
		Timestamp: time.Now(),
	}

	e.arguments = append(e.arguments, arg)
	return arg
}

func (e *Environment) GetArguments() []Argument {
	e.mu.RLock()
	defer e.mu.RUnlock()

	argsCopy := make([]Argument, len(e.arguments))
	copy(argsCopy, e.arguments)
	return argsCopy
}

func (e *Environment) CountArguments() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.arguments)
}

func (e *Environment) FormatTranscript() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if len(e.arguments) == 0 {
		return "Aucun argument n'a encore été échangé."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Thème du débat : %s\n\n", e.topic))
	for _, arg := range e.arguments {
		sb.WriteString(fmt.Sprintf("[Tour %d] %s (%s) :\n%s\n\n", arg.Round, arg.Author, arg.Role, arg.Content))
	}

	return strings.TrimSpace(sb.String())
}
