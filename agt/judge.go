package agt

import (
	"fmt"
	"log"
	"strings"

	"debater/env"
	net "debater/net"
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

func (j *JudgeAgent) ID() string {
	return j.id
}

func (j *JudgeAgent) InChan() chan<- Message {
	return j.inChan
}

func (j *JudgeAgent) Start() {
	go func() {
		log.Printf("[%s] Juge démarré (Modèle: %s)", j.id, j.model)

		for msg := range j.inChan {
			switch msg.Type {
			case Verdict:
				transcript := j.Percept()
				verdictText, err := j.Deliberate(transcript)
				if err != nil {
					log.Printf("[%s] Erreur délibération: %v", j.id, err)
					verdictText = "Erreur lors de la délibération du juge."
				}
				j.Act(msg.Reply, verdictText)

			case Stop:
				log.Printf("[%s] Arrêt du juge.", j.id)
				return
			}
		}
	}()
}

func (j *JudgeAgent) Percept() string {
	return j.env.FormatTranscript()
}

func (j *JudgeAgent) Deliberate(transcript string) (string, error) {
	systemPrompt := fmt.Sprintf(`Tu es %s, un juge d'éloquence et de débat contradictoire impartial, rigoureux et neutre.
Thème du débat : "%s".
Critères d'évaluation :
1. Solidité et clarté des arguments.
2. Capacité à réfuter les arguments adverses plutôt que de répéter sa position.
3. Cohérence du raisonnement.

Consignes :
- Fais une analyse brève des prestations de chaque côté (POUR et CONTRE).
- Déclare clairement le gagnant : Gagnant : [Nom ou Camp] (ou Égalité si vraiment impossible de départager).
- Justifie ton choix en 2 à 3 points clés.`, j.id, j.env.Topic())

	userPrompt := fmt.Sprintf(`Voici l'intégralité du débat :
----------------------------------------
%s
----------------------------------------

Donne ton analyse et ton verdict final :`, transcript)

	response, err := j.client.Query(j.model, systemPrompt, userPrompt)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

func (j *JudgeAgent) Act(replyChan chan<- Message, content string) {
	if replyChan == nil {
		log.Printf("[%s] Aucun canal de réponse fourni.", j.id)
		return
	}

	replyChan <- Message{
		From:    j.id,
		Type:    Verdict,
		Content: content,
	}
}
