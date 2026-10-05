package main

import (
	"fmt"
	"log"

	net "debater/net"
)

func main() {
	fmt.Println("=== Test du package debater/net ===")

	client, err := net.NewClientGroq()
	if err != nil {
		log.Fatalf("Erreur initialisation client Groq: %v", err)
	}

	model := "openai/gpt-oss-20b"
	fmt.Printf("Modèle sélectionné: %s\n\n", model)

	// Test 1: Simple requête avec Prompt Système
	fmt.Println("--- Test 1 : Query simple ---")
	systemPrompt := "Tu es un assistant concis et formel."
	userPrompt := "Explique en une phrase ce qu'est un système multi-agents."

	reponse, err := client.Query(model, systemPrompt, userPrompt)
	if err != nil {
		log.Fatalf("Erreur lors de la requête: %v", err)
	}
	fmt.Printf("Réponse : %s\n\n", reponse)

	// Test 2: Conversation avec historique
	fmt.Println("--- Test 2 : Conversation avec historique ---")
	historique := []net.ChatMessage{
		{Role: "user", Content: "Bonjour, je m'appelle Alice et je défends la thèse du libre arbitre."},
		{Role: "assistant", Content: "Bonjour Alice, ravi de l'apprendre. Quel est votre premier argument ?"},
		{Role: "user", Content: "Rappelle-moi mon prénom et ce que je défends brièvement."},
	}

	reponseHistorique, err := client.Chat(model, historique)
	if err != nil {
		log.Fatalf("Erreur lors de la conversation: %v", err)
	}
	fmt.Printf("Réponse : %s\n", reponseHistorique)
	fmt.Println("\n=== Test com validé avec succès ! ===")
}
