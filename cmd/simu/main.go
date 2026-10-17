package main

import (
	"fmt"
	"log"
	"sync"

	"debater/agt"
	"debater/env"
	net "debater/net"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("   SIMULATION MULTI-AGENTS : DÉBAT & VERDICT      ")
	fmt.Println("==================================================")

	// 1. Initialisation du client Groq
	client, err := net.NewClientGroq()
	if err != nil {
		log.Fatalf("Erreur client Groq : %v", err)
	}

	// 2. Création de l'environnement partagé
	sujet := "Le vote obligatoire devrait-il être instauré dans toutes les démocraties ?"
	debatEnv := env.NewEnvironment(sujet)
	fmt.Printf("\nSujet du débat : « %s »\n\n", sujet)

	// Modèles utilisés (possibilité d'associer des modèles différents)
	modelDebaters := "openai/gpt-oss-20b"
	modelJudge := "openai/gpt-oss-20b"
	modelSummarizer := "openai/gpt-oss-20b"

	// 3. Instanciation des agents
	alice := agt.NewDebaterAgent("Alice", "POUR", modelDebaters, client, debatEnv)
	bob := agt.NewDebaterAgent("Bob", "CONTRE", modelDebaters, client, debatEnv)
	judge := agt.NewJudgeAgent("Juge-Rhadamanthe", modelJudge, client, debatEnv)
	summarizer := agt.NewSummarizerAgent("Synthétiseur-Hermès", modelSummarizer, client, debatEnv)

	// 4. Lancement des goroutines des agents
	alice.Start()
	bob.Start()
	judge.Start()
	summarizer.Start()

	// 5. Création et lancement du modérateur (débat en 2 rounds)
	debaters := []*agt.DebaterAgent{alice, bob}
	moderator := agt.NewModeratorAgent("Modérateur-Socrate", debatEnv, debaters, 2, client, modelDebaters)

	moderator.RunDebate()

	// 6. Summarizer et Juge parallèle

	var wg sync.WaitGroup
	wg.Add(2)

	var summaryContent, verdictContent string

	summaryChan := make(chan agt.Message, 1)
	go func() {
		defer wg.Done()
		msg := <-summaryChan
		summaryContent = msg.Content
	}()

	judgeChan := make(chan agt.Message, 1)
	go func() {
		defer wg.Done()
		msg := <-judgeChan
		verdictContent = msg.Content
	}()

	summarizer.InChan() <- agt.Message{From: "Main", Type: agt.Summarize, Reply: summaryChan}
	judge.InChan() <- agt.Message{From: "Main", Type: agt.Verdict, Reply: judgeChan}

	wg.Wait()

	fmt.Println("\n--------------------------------------------------")
	fmt.Println("   SYNTHÈSE DU DÉBAT (AGENT SYNTHÉTISEUR)         ")
	fmt.Println("--------------------------------------------------")
	fmt.Println(summaryContent)

	fmt.Println("\n--------------------------------------------------")
	fmt.Println("   VERDICT FINAL (AGENT JUGE)                     ")
	fmt.Println("--------------------------------------------------")
	fmt.Println(verdictContent)

	// 7. Arrêt propre des agents restants
	summarizer.InChan() <- agt.Message{Type: agt.Stop}
	judge.InChan() <- agt.Message{Type: agt.Stop}

	fmt.Println("\n==================================================")
	fmt.Println("   FIN DE LA SIMULATION DU DÉBAT                  ")
	fmt.Println("==================================================")
}
