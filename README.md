# Go Multi-Agent Debater 🗣️🤖

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Groq API](https://img.shields.io/badge/LLM-Groq%20Cloud-F55036?style=flat)](https://groq.com/)
[![Architecture](https://img.shields.io/badge/Architecture-Multi--Agent%20System-blueviolet?style=flat)]()
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A concurrent **Multi-Agent System (MAS)** in Go that simulates structured, argumentative debates between autonomous AI personas. Powered by **Groq Cloud LLMs**, agents communicate asynchronously via **goroutines** and **Go channels**, sharing a thread-safe blackboard environment.

---

## 🌟 Highlights

- **Concurrence & Asynchronisme** : Chaque agent fonctionne comme un processus indépendant via des **goroutines** et communique par passage de messages sur des **canaux Go typés** (`chan Message`).
- **Environnement Partagé Thread-Safe** : Gestion du tableau noir (*blackboard transcript*) protégée par un `sync.RWMutex` pour garantir l'absence de *data races* lors des accès concurrents.
- **Inférence LLM Ultra-Rapide** : Client HTTP natif Go connecté à l'API **Groq** pour une génération de répliques et d'analyses en temps réel.
- **Rôles Dialectiques Spécialisés** : Modération dynamique, argumentation contradictoire contradictoire (thèse / antithèse), synthèse neutre et arbitrage impartial.

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

| Agent | Rôle & Responsabilité | Comportement & Modèle |
| :--- | :--- | :--- |
| **Moderator (Socrate)** | Maître de cérémonie | Coordonne les tours de parole, lance les rounds de débat, sollicite la synthèse et le verdict. |
| **Debater POUR (Alice)** | Défenseur de la thèse | Formule des arguments convaincants, réfute les objections de l'adversaire et s'appuie sur l'historique. |
| **Debater CONTRE (Bob)** | Opposant contradictoire | Contredit les arguments adverses, propose des contre-exemples et soulève les limites pratiques/morales. |
| **Summarizer (Hermès)** | Synthétiseur neutre | Extrait les points de convergence et de divergence majeurs sans prendre parti. |
| **Judge (Rhadamanthe)** | Juge arbitre | Évalue la rigueur logique, la qualité rhétorique et délivre un verdict motivé avec déclaration du vainqueur. |

---

## 📂 Structure du Projet

```text
debater/
├── agt/               # Définitions et logique des agents autonomes
│   ├── types.go       # Types de messages (GiveFloor, Argument, Summarize, Verdict)
│   ├── debater.go     # Agent débatteur (persona POUR/CONTRE)
│   ├── moderator.go   # Agent modérateur
│   ├── summarizer.go  # Agent de synthèse neutre
│   └── judge.go       # Agent juge / arbitre
├── cmd/               # Points d'entrée exécutables
│   ├── simu/          # Simulation complète (Main de production)
│   ├── tst-agt/       # Tests unitaires / fonctionnels inter-agents
│   └── tst-groq/      # Validation de connectivité avec l'API Groq
├── env/               # Environnement partagé thread-safe (RWMutex)
│   └── env.go
├── net/               # Client HTTP REST pour l'API Groq
│   └── client.go
├── .env.example       # Modèle de configuration des variables d'environnement
├── go.mod
└── README.md
```

---

## 🚀 Démarrage Rapide

### 1. Prérequis
- [Go 1.23+](https://go.dev/dl/)
- Une clé d'API [Groq Cloud](https://console.groq.com/keys) (gratuite)

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
Puis renseignez votre clé d'API Groq dans `.env` :
```env
GROQ_API_KEY=gsk_votre_cle_api_groq_ici
```

### 4. Lancer la simulation de débat
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

## 🛠️ Stack Technique

- **Langage** : Go (Golang)
- **Modèles de langage testés** : `openai/gpt-oss-20b`, `llama-3.3-70b-versatile`, `mixtral-8x7b-32768` via Groq
- **Design Patterns** : Actor Pattern (Agent Goroutine Loop), Blackboard Architecture, Producer-Consumer (Channels)
- **Gestion des dépendances** : Go Modules (`godotenv`)

---

## 👤 Auteur

**Jean Van Gysel**
- GitHub : [@JeanVG23](https://github.com/JeanVG23)

---

## 📄 Licence

Ce projet est sous licence MIT. Libre d'utilisation et d'adaptation.
