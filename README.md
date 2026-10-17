# Go Multi-Agent Debater

Projet réalisé dans le cadre de l'UV **IA04 (Systèmes distribués et multi-agents)** à l'UTC. Il modélise un système multi-agents en **Go** : chaque agent est une goroutine qui reçoit des messages typés sur un canal.

Le programme organise un débat contradictoire entre deux débatteurs (POUR et CONTRE), arbitré par un modérateur, puis résumé par un synthétiseur et tranché par un juge. Chaque rôle est un appel à un modèle de langage via l'API **Groq**, avec son propre prompt.

---

## Ce que le projet met en œuvre

- **Agents en goroutines** : chaque agent tourne dans sa propre goroutine et lit ses messages sur un canal (`chan Message`), à la manière du modèle acteur.
- **Messages typés** : cinq types (`GiveFloor`, `Argument`, `Summarize`, `Verdict`, `Stop`). Chaque message porte un canal de réponse.
- **Environnement partagé (blackboard)** : le transcript du débat est protégé par un `sync.RWMutex`. Seul le modérateur écrit ; les autres agents lisent.
- **Fan-out en fin de débat** : `main` envoie la demande de synthèse et la demande de verdict sans attendre la première réponse, puis attend les deux (`sync.WaitGroup`). Les deux appels au modèle s'exécutent en parallèle.
- **Cycle perception, délibération, action** : un débatteur lit le transcript, appelle le modèle, puis renvoie son argument au modérateur.
- **Client HTTP Go** : appels à l'API de Groq avec la bibliothèque standard, sans SDK.

---

## Architecture

```text
                     +--------------------------+
                     |    Moderator (Socrate)   |
                     +--------------------------+
                              |        |
             Donne la parole  |        |  Donne la parole
             (GiveFloor)      v        v  (GiveFloor)
                 +--------------+    +---------------+
                 | Debater POUR |    | Debater CONTRE|
                 |   (Alice)    |    |    (Bob)      |
                 +--------------+    +---------------+
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

### Rôles

| Agent                    | Rôle                  | Comportement                                                                              |
| :----------------------- | :-------------------- | :---------------------------------------------------------------------------------------- |
| **Moderator (Socrate)**  | Animateur             | Donne la parole à tour de rôle, enregistre chaque réplique dans l'environnement.          |
| **Debater POUR (Alice)** | Défend la thèse       | Lit le transcript, argumente et réfute son adversaire.                                    |
| **Debater CONTRE (Bob)** | S'oppose à la thèse   | Lit le transcript, contredit les arguments adverses.                                      |
| **Summarizer (Hermès)**  | Synthèse              | Résume les points d'accord et de désaccord.                                               |
| **Judge (Rhadamanthe)**  | Arbitre               | Évalue la solidité des arguments et désigne un gagnant.                                   |

---

## Structure du dépôt

```text
debater/
├── agt/               # Les agents
│   ├── types.go       # Messages typés (GiveFloor, Argument, Summarize, Verdict, Stop)
│   ├── debater.go     # Débatteur
│   ├── moderator.go   # Modérateur
│   ├── summarizer.go  # Synthétiseur
│   └── judge.go       # Juge
├── cmd/               # Points d'entrée
│   ├── simu/          # Simulation complète du débat
│   ├── tst-agt/       # Essai manuel des débatteurs
│   └── tst-groq/      # Essai manuel de la connexion à Groq
├── env/               # Environnement partagé (blackboard avec RWMutex)
│   └── env.go
├── net/               # Client HTTP pour l'API Groq
│   └── client.go
├── .env.example       # Modèle de variables d'environnement
├── go.mod
└── README.md
```

---

## Objectifs du cours

1. Prendre en main la concurrence en Go : goroutines, canaux, verrous.
2. Définir un protocole d'interaction entre agents de rôles différents (attribution de la parole, recueil des arguments).
3. Concevoir un tableau noir partagé et le protéger par un verrou lecture/écriture.
4. Cadrer des modèles de langage par des rôles et des consignes précises.

---

## Limites

- **Débat séquentiel** : pendant le débat, un seul agent travaille à la fois, car le modérateur attend la réponse d'un débatteur avant d'appeler le suivant. Le seul parallélisme est celui du synthétiseur et du juge en fin de débat. Le `RWMutex` n'est donc jamais en contention entre un lecteur et un écrivain.
- **Aucun test, par choix** : ce projet est entièrement écrit à la main, y compris ses défauts. Je n'ai fait écrire ni code ni tests à un assistant.
- **Erreurs d'appel masquées** : si un appel au modèle échoue, le débatteur insère un argument de repli dans le transcript et le juge renvoie un message d'erreur à la place du verdict. L'échec n'interrompt pas le programme et se confond avec une sortie normale.
- **Pas de `context.Context`** : aucun contrôle fin des délais d'appel, ni d'arrêt propre de l'ensemble des goroutines.
- **Un seul modèle** : tous les agents utilisent `openai/gpt-oss-20b`, et la qualité des arguments n'est pas évaluée.

Pistes d'évolution :

- un protocole asynchrone, où un débatteur peut demander la parole ;
- un jury de plusieurs juges, en prolongeant le fan-out actuel ;
- des tests Go (`go test -race ./...`) avec un faux serveur HTTP (`httptest.Server`).

---

## Démarrage

### Prérequis

- [Go 1.23+](https://go.dev/dl/)
- Une clé d'API [Groq](https://console.groq.com/keys)

### Installation

```bash
git clone git@github.com:JeanVG23/go-multiagent-debater.git
cd go-multiagent-debater
```

### Configuration

Créez le fichier `.env` à la racine du projet :

```bash
cp .env.example .env
```

Renseignez votre clé d'API Groq dans `.env` :

```env
GROQ_API_KEY=gsk_votre_cle_api_groq_ici
```

### Lancement

```bash
go run cmd/simu/main.go
```

---

## Exemple d'exécution

Sortie réelle (stdout) du 17 octobre 2026, modèle `openai/gpt-oss-20b`, tronquée avec `[...]`. Les journaux de coordination vont sur stderr.

```text
==================================================
   SIMULATION MULTI-AGENTS : DÉBAT & VERDICT
==================================================

Sujet du débat : « Le vote obligatoire devrait-il être instauré dans toutes les démocraties ? »

[Alice (POUR)] :
Le vote obligatoire garantit la représentativité réelle de chaque citoyen, renforçant la légitimité des institutions démocratiques. [...]

[Bob (CONTRE)] :
La participation obligatoire viole la liberté de choix, pénalisant ceux qui refusent de voter et sapant le respect de l'autonomie individuelle. [...]

[Alice (POUR)] : [...]

[Bob (CONTRE)] : [...]

--------------------------------------------------
   SYNTHÈSE DU DÉBAT (AGENT SYNTHÉTISEUR)
--------------------------------------------------
**Camp « Pour » (Alice)**
- Vote obligatoire renforce la représentativité et la légitimité des institutions.
[...]
**Camp « Contre » (Bob)**
- La contrainte viole la liberté de choix et l'autonomie individuelle.
[...]

--------------------------------------------------
   VERDICT FINAL (AGENT JUGE)
--------------------------------------------------
[...]
**Gagnant : Bob (CONTRE)**
[...]
```

Le texte généré n'est pas vérifié : Bob invoque des « enquêtes » qu'il ne source pas.

---

## Stack

- **Langage** : Go
- **Modèle** : `openai/gpt-oss-20b` via l'API Groq
- **Dépendance** : `github.com/joho/godotenv`

---

## Auteur

**Jean Van Gysel** ([@JeanVG23](https://github.com/JeanVG23))

Le code est de ma main ; ce README a été rédigé avec l'aide d'une IA.

Projet sous licence MIT.
