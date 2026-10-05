package agt

import (
	"fmt"
	"log"
	"strings"

	"debater/env"
	net "debater/net"
)

type DebaterAgent struct {
	id     string
	side   string
	model  string
	client *net.ClientGroq
	env    *env.Environment
	inChan chan Message
}

func NewDebaterAgent(id string, side string, model string, client *net.ClientGroq, environment *env.Environment) *DebaterAgent {
	return &DebaterAgent{
		id:     id,
		side:   side,
		model:  model,
		client: client,
		env:    environment,
		inChan: make(chan Message, 5),
	}
}

func (ag *DebaterAgent) ID() string {
	return ag.id
}

func (ag *DebaterAgent) Side() string {
	return ag.side
}

func (ag *DebaterAgent) InChan() chan<- Message {
	return ag.inChan
}

func (ag *DebaterAgent) Start() {
	go func() {
		log.Printf("[%s] Débatteur démarré (Position: %s, Modèle: %s)", ag.id, ag.side, ag.model)

		for msg := range ag.inChan {
			switch msg.Type {
			case GiveFloor:
				// 1. Perception : consultation de l'état actuel du débat dans l'environnement
				transcript := ag.Percept()

				// 2. Délibération : appel LLM Groq avec son prompt et l'historique
				argText, err := ag.Deliberate(transcript)
				if err != nil {
					log.Printf("[%s] Erreur délibération: %v", ag.id, err)
					argText = fmt.Sprintf("En tant que défenseur du %s, je maintiens fermement ma position.", ag.side)
				}

				// 3. Action : renvoi de l'argument au modérateur via le canal Reply
				ag.Act(msg.Reply, argText)

			case Stop:
				log.Printf("[%s] Arrêt de l'agent demandé.", ag.id)
				return
			}
		}
	}()
}

func (ag *DebaterAgent) Percept() string {
	return ag.env.FormatTranscript()
}

func (ag *DebaterAgent) Deliberate(transcript string) (string, error) {
	systemPrompt := fmt.Sprintf(`Tu es %s. Tu participes à un débat contradictoire sur le thème : "%s".
Ta position OBLIGATOIRE est : %s.
Consignes impératives :
- Défends ta position avec conviction, clarté et concision (2 à 4 phrases maximum).
- Si l'autre débatteur a déjà parlé, réfute directement ses arguments en apportant de nouveaux éléments.
- Reste poli mais percutant.
- Ne parle JAMAIS à la place du modérateur ni de ton adversaire.
- Donne uniquement ton argument, sans fioritures méta (pas de "Voici mon argument :").`,
		ag.id, ag.env.Topic(), ag.side)

	userPrompt := fmt.Sprintf(`Voici l'état actuel du débat :
----------------------------------------
%s
----------------------------------------

C'est à ton tour de parler. Formule ton argument en tant que défenseur du %s :`, transcript, ag.side)

	response, err := ag.client.Query(ag.model, systemPrompt, userPrompt)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

func (ag *DebaterAgent) Act(replyChan chan<- Message, content string) {
	if replyChan == nil {
		log.Printf("[%s] Attention: aucun canal de réponse fourni.", ag.id)
		return
	}

	replyChan <- Message{
		From:    ag.id,
		Type:    Argument,
		Content: content,
	}
}
