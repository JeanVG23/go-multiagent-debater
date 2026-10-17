package agt

import (
	"debater/env"
	"debater/net"
	"fmt"
	"log"
	"strings"
)

type JudgeAgent struct {
	id     string
	model  string
	client *net.ClientGroq
	env    *env.Environment
	inChan chan Message
}

func NewJudgeAgent(id string, model string, client *net.ClientGroq, environment *env.Environment) *JudgeAgent {
	return &JudgeAgent{
		id:     id,
		model:  model,
		client: client,
		env:    environment,
		inChan: make(chan Message, 5),
	}
}

func (ja *JudgeAgent) ID() string {
	return ja.id
}

func (ja *JudgeAgent) InChan() chan<- Message {
	return ja.inChan
}

func (ja *JudgeAgent) Start() {
	go func() {
		log.Printf("[%s] Juge démarré (Modèle: %s)", ja.id, ja.model)

		for msg := range ja.inChan {
			switch msg.Type {
			case Verdict:
				transcript := ja.Percept()

				verdictText, err := ja.Deliberate(transcript)
				if err != nil {
					log.Printf("[%s] Erreur délibération: %v", ja.id, err)
					verdictText = "Erreur lors de la délibération du juge."
				}

				ja.Act(msg.Reply, verdictText)

			case Stop:
				log.Printf("[%s] Arrêt du juge.", ja.id)
				return
			}
		}
	}()
}

func (ja *JudgeAgent) Percept() string {
	return ja.env.FormatTranscript()
}

func (ja *JudgeAgent) Deliberate(transcript string) (string, error) {
	systemPrompt := fmt.Sprintf(`Tu es %s, un juge d'éloquence et de débat contradictoire impartial, rigoureux et neutre.
Thème du débat : "%s".
Critères d'évaluation :
1. Solidité et clarté des arguments.
2. Capacité à réfuter les arguments adverses plutôt que de répéter sa position.
3. Cohérence du raisonnement.

Consignes :
- Fais une analyse brève des prestations de chaque côté (POUR et CONTRE).
- Déclare clairement le gagnant : Gagnant : [Nom ou Camp] (ou Égalité si impossible de départager).
- Justifie ton choix en 2 à 3 points clés.`, ja.id, ja.env.Topic())

	userPrompt := fmt.Sprintf(`Voici l'intégralité du débat :
----------------------------------------
%s
----------------------------------------

Donne ton analyse et ton verdict final :`, transcript)

	response, err := ja.client.Query(ja.model, systemPrompt, userPrompt)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response), nil
}

func (ja *JudgeAgent) Act(replyChan chan<- Message, content string) {
	if replyChan == nil {
		log.Printf("[%s] Attention: aucun canal de réponse fourni.", ja.id)
		return
	}
	replyChan <- Message{
		From:    ja.id,
		Type:    Verdict,
		Content: content,
	}
}
