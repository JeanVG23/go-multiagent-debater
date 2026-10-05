package agt

import (
	"fmt"
	"log"
	"strings"

	"debater/env"
	net "debater/net"
)

type SummarizerAgent struct {
	id     string
	model  string
	client *net.ClientGroq
	env    *env.Environment
	inChan chan Message
}

func NewSummarizerAgent(id string, model string, client *net.ClientGroq, environment *env.Environment) *SummarizerAgent {
	return &SummarizerAgent{
		id:     id,
		model:  model,
		client: client,
		env:    environment,
		inChan: make(chan Message, 5),
	}
}

func (s *SummarizerAgent) ID() string {
	return s.id
}

func (s *SummarizerAgent) InChan() chan<- Message {
	return s.inChan
}

func (s *SummarizerAgent) Start() {
	go func() {
		log.Printf("[%s] Synthétiseur démarré (Modèle: %s)", s.id, s.model)

		for msg := range s.inChan {
			switch msg.Type {
			case Summarize:
				transcript := s.Percept()
				summary, err := s.Deliberate(transcript)
				if err != nil {
					log.Printf("[%s] Erreur délibération: %v", s.id, err)
					summary = "Erreur lors de la génération du résumé."
				}
				s.Act(msg.Reply, summary)

			case Stop:
				log.Printf("[%s] Arrêt du synthétiseur.", s.id)
				return
			}
		}
	}()
}

func (s *SummarizerAgent) Percept() string {
	return s.env.FormatTranscript()
}

func (s *SummarizerAgent) Deliberate(transcript string) (string, error) {
	systemPrompt := fmt.Sprintf(`Tu es %s, un agent de synthèse neutre et objectif.
Thème du débat : "%s".
Consignes :
- Résume de façon très concise les points clés soulevés par chaque débatteur.
- Ne prends parti pour aucun camp.
- Utilise des puces courtes pour chaque camp.`, s.id, s.env.Topic())

	userPrompt := fmt.Sprintf(`Voici le transcript du débat :
----------------------------------------
%s
----------------------------------------

Donne une synthèse brève des arguments :`, transcript)

	response, err := s.client.Query(s.model, systemPrompt, userPrompt)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

func (s *SummarizerAgent) Act(replyChan chan<- Message, content string) {
	if replyChan == nil {
		log.Printf("[%s] Aucun canal de réponse fourni.", s.id)
		return
	}

	replyChan <- Message{
		From:    s.id,
		Type:    Summarize,
		Content: content,
	}
}
