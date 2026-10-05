package agt

import (
	"fmt"
	"log"

	"debater/env"
	net "debater/net"
)

type ModeratorAgent struct {
	id        string
	env       *env.Environment
	debaters  []*DebaterAgent
	rounds    int
	replyChan chan Message
	client    *net.ClientGroq
	model     string
}

func NewModeratorAgent(id string, environment *env.Environment, debaters []*DebaterAgent, rounds int, client *net.ClientGroq, model string) *ModeratorAgent {
	return &ModeratorAgent{
		id:        id,
		env:       environment,
		debaters:  debaters,
		rounds:    rounds,
		replyChan: make(chan Message),
		client:    client,
		model:     model,
	}
}

func (mod *ModeratorAgent) RunDebate() {
	log.Printf("[%s] Début du débat : « %s » (%d rounds)", mod.id, mod.env.Topic(), mod.rounds)

	for round := 1; round <= mod.rounds; round++ {
		log.Printf("\n========== ROUND %d / %d ==========", round, mod.rounds)

		for _, debater := range mod.debaters {
			log.Printf("[%s] Donne la parole à %s (%s)", mod.id, debater.ID(), debater.Side())

			// 1. Donner la parole au débatteur
			debater.InChan() <- Message{
				From:  mod.id,
				Type:  GiveFloor,
				Reply: mod.replyChan,
			}

			// 2. Attendre la réponse du débatteur
			reply := <-mod.replyChan

			// 3. Seul le modérateur écrit dans l'environnement
			mod.env.AddArgument(reply.From, debater.Side(), reply.Content, round)
			fmt.Printf("\n[%s (%s)] :\n%s\n\n", reply.From, debater.Side(), reply.Content)
		}
	}

	// 4. Clôture : arrêt des débatteurs
	log.Printf("[%s] Fin des %d rounds. Arrêt des débatteurs.", mod.id, mod.rounds)
	for _, debater := range mod.debaters {
		debater.InChan() <- Message{
			From: mod.id,
			Type: Stop,
		}
	}
}
