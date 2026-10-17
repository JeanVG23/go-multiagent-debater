# Go Multi-Agent Debater 🗣️🤖

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Groq API](https://img.shields.io/badge/LLM-Groq%20Cloud-F55036?style=flat)](https://groq.com/)
[![Course](https://img.shields.io/badge/Academic-IA04%20Multi--Agent%20Systems-purple?style=flat)]()
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Projet académique développé dans le cadre de l'UV **IA04 (Systèmes Multi-Agents)** explorant la modélisation de systèmes multi-agents et l'initiation aux paradigmes de concurrence en **Go** (goroutines, canaux typés, synchronisation mémoire).

Le projet simule une arène de débat contradictoire structurée entre des personas IA autonomes propulsées par l'API **Groq Cloud**. Chaque agent s'exécute dans sa propre boucle concurrente (goroutine worker) et communique par passage de messages typés, orchestré selon un protocole délibératif au tour par tour.

---

## 🌟 Points Clés & Paradigmes Explorés

- **Apprentissage des Primitives Concurrentes** : Modélisation d'agents indépendants s'exécutant dans des **goroutines** dédiées et communiquant par passage de messages sur des **canaux Go typés** (`chan Message`).
- **Environnement Partagé Thread-Safe (Blackboard)** : Gestion d'un tableau noir partagé protégé par un `sync.RWMutex` pour garantir l'intégrité des accès concurrents (lectures concurrentes par les agents observateurs, écriture contrôlée par le modérateur).
- **Architecture Multi-Agents Formelle** : Implémentation du cycle canonique **Perception - Délibération - Action** pour chaque entité du système.
- **Rôles Dialectiques Spécialisés** : Modération, argumentation contradictoire contrainte (thèse / antithèse), synthèse neutre et arbitrage impartial.
- **Inférence LLM Ultra-Rapide** : Client HTTP natif Go connecté à **Groq** pour une génération temps réel avec personas et contraintes de prompt strictes.

---

## 🏛️ Architecture Multi-Agents

```text
                     +--------------------------+
                     |    Moderator (Socrate)   |
                     +--------------------------+
                              |        |
             Donne la parole  |        |  Donne la parole
             (GiveFloor)      v        v  (GiveFloor)
                 +--------------+    +--------------+
                 | Debater POUR |    | Debater CONTRE|
                 |   (Alice)    |    |    (Bob)     |
                 +--------------+    +--------------+
                        \                 /
          Publie réplique\               / Publie réplique
                          v             v
                    +-------------------------+
                    |  Shared Environment     |  <-- sync.RWMutex
                    |  (Transcript & Topic)   |
                    +-------------------------+
                              |        |
           Lecture transcript |        | Lecture transcript
                              v        v
                 +--------------+    +--------------+
                 |  Summarizer  |    |    Judge     |
                 |   (Hermès)   |    | (Rhadamanthe)|
                 +--------------+    +--------------+
```

### 🎭 Rôles des Agents

| Agent | Rôle & Responsabilité | Comportement & Cycle |
| :--- | :--- | :--- |
| **Moderator (Socrate)** | Maître de cérémonie | Coordonne les tours de parole, applique le protocole de débat, sollicite la synthèse et le verdict. |
| **Debater POUR (Alice)** | Défenseur de la thèse | Perçoit l'état du débat, formule des arguments convaincants et réfute les objections adverses. |
| **Debater CONTRE (Bob)** | Opposant contradictoire | Perçoit l'état du débat, contredit les arguments adverses et soulève les limites pratiques/morales. |
| **Summarizer (Hermès)** | Synthétiseur neutre | Extrait les points de convergence et de divergence majeurs sans parti pris. |
| **Judge (Rhadamanthe)** | Juge arbitre | Évalue la rigueur logique et rhétorique, délivre un verdict motivé et désigne le vainqueur. |

---

## 📂 Structure du Projet

```text
debater/
├── agt/               # Définitions et logique des agents autonomes
│   ├── types.go       # Messages typés (GiveFloor, Argument, Summarize, Verdict, Stop)
│   ├── debater.go     # Agent débatteur (cycle perception-délibération-action)
│   ├── moderator.go   # Agent modérateur orchestrateur
│   ├── summarizer.go  # Agent de synthèse neutre
│   └── judge.go       # Agent juge / arbitre
├── cmd/               # Points d'entrée exécutables
│   ├── simu/          # Point d'entrée principal de la simulation de débat
│   ├── tst-agt/       # Script de validation rapide des échanges inter-agents
│   └── tst-groq/      # Script de validation de la connectivité API Groq
├── env/               # Environnement partagé thread-safe (Blackboard avec RWMutex)
│   └── env.go
├── net/               # Client HTTP REST natif pour l'API Groq
│   └── client.go
├── .env.example       # Modèle de variables d'environnement
├── go.mod
└── README.md
```

---

## 🎯 Contexte & Objectifs Pédagogiques

Ce projet a été conçu dans le cadre du cours de **Systèmes Multi-Agents (IA04)** à l'Université de Technologie de Compiègne (UTC).

Les objectifs d'apprentissage principaux étaient :
1. **L'exploration des paradigmes concurrents en Go** : mise en place de structures d'agents autonomes basées sur des boucles d'écoute concurrentes et l'échange de messages typés sur canaux.
2. **La coordination d'agents hétérogènes** : définition d'un protocole d'interaction délibératif (attribution de parole, recueil des arguments, passage de témoin).
3. **La gestion de mémoire partagée concurrente** : conception d'un tableau noir (*Blackboard*) sécurisé par verrous R/W pour prévenir les accès concurrents non synchronisés.
4. **L'intégration pragmatique de LLMs** : cadrage de modèles de langage par des rôles dialectiques et des contraintes d'évaluation rigides.

---

## 🔭 Limites Identifiées & Perspectives d'Évolution

Dans une démarche d'ingénierie logicielle et de recul critique sur l'architecture :

- **Orchestration Tour par Tour vs Plein Asynchronisme** : Bien que chaque agent dispose de sa propre goroutine et de son canal d'entrée, l'orchestration par le modérateur demeure pour l'instant synchrone au niveau logique (attente de la réponse avant d'attribuer le tour suivant). Une perspective d'évolution est d'implémenter un protocole asynchrone avec interruptions et demandes spontanées de droit de réponse.
- **Parallélisation de l'Évaluation (Fan-out / Fan-in)** : Permettre à un jury composé de multiples agents d'évaluer le débat en parallèle à chaque tour à l'aide de `sync.WaitGroup`.
- **Gestion de Cycle de Vie avec `context.Context`** : Intégrer la propagation de contextes pour la gestion fine des timeouts d'inférence HTTP et l'arrêt contrôlé (*graceful shutdown*) de l'ensemble des goroutines.
- **Suite de Tests Découplée** : Remplacer les scripts de validation manuels par des tests Go idiomatiques (`go test -race ./...`) s'appuyant sur des interfaces et des mocks de transport HTTP (`httptest.Server`).

---

## 🚀 Démarrage Rapide

### 1. Prérequis
- [Go 1.23+](https://go.dev/dl/)
- Une clé d'API [Groq Cloud](https://console.groq.com/keys)

### 2. Installation
```bash
git clone git@github.com:JeanVG23/go-multiagent-debater.git
cd go-multiagent-debater
```

### 3. Configuration
Créez votre fichier d'environnement `.env` à la racine du projet :
```bash
cp .env.example .env
```
Renseignez votre clé d'API Groq dans `.env` :
```env
GROQ_API_KEY=gsk_votre_cle_api_groq_ici
```

### 4. Lancer la simulation
```bash
go run cmd/simu/main.go
```

---

## 🧪 Exemple d'Exécution

```text
==================================================
   SIMULATION MULTI-AGENTS : DÉBAT & VERDICT      
==================================================

Sujet du débat : « Le vote obligatoire devrait-il être instauré dans toutes les démocraties ? »

[Modérateur-Socrate] Début du débat en 2 rounds.
--- ROUND 1 ---
[Alice (POUR)] Le vote obligatoire renforce la légitimité démocratique et incite les citoyens à s'informer...
[Bob (CONTRE)] Forcer le vote porte atteinte à la liberté individuelle fondamentale de ne pas exprimer d'opinion...

--- ROUND 2 ---
[Alice (POUR)] En réponse à Bob, le vote blanc demeure une option de contestation tout en assurant une représentativité maximale...
[Bob (CONTRE)] La représentativité forcée ne résout pas la défiance politique; elle ne fait que la masquer...

--------------------------------------------------
   SYNTHÈSE DU DÉBAT (AGENT SYNTHÉTISEUR)         
--------------------------------------------------
Synthèse neutre d'Hermès :
- Points d'accord : Importance de la légitimité des institutions.
- Points de friction : Liberté individuelle versus devoir civique collectif.

--------------------------------------------------
   VERDICT DU JUGE (AGENT JUGE)                   
--------------------------------------------------
Verdict de Rhadamanthe :
- Analyse rhétorique et solidité argumentative...
- Vainqueur désigné : Alice (POUR) pour sa prise en compte constructive du vote blanc face à la critique libertarienne.
```

---

## 🛠️ Stack & Concepts

- **Langage** : Go (Golang)
- **Modèles LLM** : `openai/gpt-oss-20b`, `llama-3.3-70b-versatile` via Groq Cloud
- **Paradigmes** : Agent Worker Loops (Goroutines), Communication par canaux (`chan Message`), Blackboard Pattern (`sync.RWMutex`)
- **Dépendances** : `github.com/joho/godotenv`

---

## 👤 Auteur

**Jean Van Gysel**
- GitHub : [@JeanVG23](https://github.com/JeanVG23)

---

## 📄 Licence

Ce projet est sous licence MIT.
