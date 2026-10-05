package main

import (
	"fmt"
	"log"

	"debater/agt"
	"debater/env"
	net "debater/net"
)

func main() {
	fmt.Println("=== Test de l'agent débatteur (cmd/tst-agt) ===")

	// 1. Initialiser le client Groq et l'environnement partagé
	client, err := net.NewClientGroq()
	if err != nil {
		log.Fatalf("Erreur client Groq : %v", err)
	}

	sujet := "Le télétravail complet est-il bénéfique pour les entreprises et les salariés ?"
	debatEnv := env.NewEnvironment(sujet)
	fmt.Printf("Thème : %s\n\n", sujet)

	// 2. Créer deux débatteurs (un POUR et un CONTRE)
	model := "openai/gpt-oss-20b"
	alice := agt.NewDebaterAgent("Alice", "POUR", model, client, debatEnv)
	bob := agt.NewDebaterAgent("Bob", "CONTRE", model, client, debatEnv)

	// Démarrer les goroutines des agents
	alice.Start()
	bob.Start()

	// Canal simulant le modérateur qui reçoit les arguments
	moderatorChan := make(chan agt.Message)

	// 3. Tour 1 : Donner la parole à Alice (POUR)
	fmt.Println("--> Le modérateur donne la parole à Alice...")
	alice.InChan() <- agt.Message{
		From:  "Moderator",
		Type:  agt.GiveFloor,
		Reply: moderatorChan,
	}

	msgAlice := <-moderatorChan
	fmt.Printf("\n[%s (%s)] :\n%s\n\n", msgAlice.From, alice.Side(), msgAlice.Content)
	debatEnv.AddArgument(msgAlice.From, alice.Side(), msgAlice.Content, 1)

	// 4. Tour 2 : Donner la parole à Bob (CONTRE), qui va percevoir l'argument d'Alice
	fmt.Println("--> Le modérateur donne la parole à Bob...")
	bob.InChan() <- agt.Message{
		From:  "Moderator",
		Type:  agt.GiveFloor,
		Reply: moderatorChan,
	}

	msgBob := <-moderatorChan
	fmt.Printf("\n[%s (%s)] :\n%s\n\n", msgBob.From, bob.Side(), msgBob.Content)
	debatEnv.AddArgument(msgBob.From, bob.Side(), msgBob.Content, 2)

	// 5. Arrêter proprement les agents
	alice.InChan() <- agt.Message{Type: agt.Stop}
	bob.InChan() <- agt.Message{Type: agt.Stop}

	fmt.Println("=== Résumé complet stocké dans l'environnement ===")
	fmt.Println(debatEnv.FormatTranscript())
	fmt.Println("=== Test débatteurs réussi avec succès ! ===")
}
