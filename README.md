# net-cat 🐱

**net-cat** est une implémentation personnalisée de l'outil réseau NetCat, entièrement développée en Go. Ce projet propose un système de chat TCP client-serveur permettant la communication en temps réel entre plusieurs utilisateurs connectés à un serveur commun.

## 📋 Table des matières

- [Présentation du projet](#-présentation-du-projet)
- [Installation](#-installation)
- [Utilisation](#-utilisation)
- [Architecture technique](#-architecture-technique)
- [Stratégies et algorithmes](#-stratégies-et-algorithmes)
- [Structure du code](#-structure-du-code)
- [Fonctionnalités](#-fonctionnalités)
- [Technologies utilisées](#-technologies-utilisées)
- [Contribution](#-contribution)
- [Licence](#-licence)

## 🎯 Présentation du projet

**net-cat** est un outil de communication réseau basé sur le protocole TCP qui permet de :
- Créer un serveur de chat TCP capable d'accueillir jusqu'à 10 utilisateurs simultanément
- Connecter des clients à un serveur distant pour échanger des messages en temps réel
- Conserver un historique des conversations dans un fichier JSON
- Afficher les messages avec des couleurs distinctives et des horodatages
- Gérer les connexions et déconnexions des utilisateurs de manière fluide

Ce projet est une réimplémentation moderne du célèbre outil NetCat de Linux, spécialement conçu pour les applications de messagerie instantanée.

## 🔧 Installation

### Prérequis

- **Go** version 1.21.5 ou supérieure
- Un système d'exploitation compatible (Linux, macOS, Windows)
- Un terminal/console pour l'exécution

### Étapes d'installation

1. **Cloner le dépôt**
```bash
git clone https://github.com/i2sac/net-cat.git
cd net-cat
```

2. **Compiler le projet**
```bash
go build -o net-cat
```

Cette commande génère un exécutable binaire nommé `net-cat` dans le répertoire courant.

3. **Vérifier l'installation**
```bash
ls -lh net-cat
```

L'exécutable devrait avoir une taille d'environ 4 Mo.

## 🚀 Utilisation

### Démarrer un serveur

Pour créer un serveur de chat sur un port spécifique :

```bash
./net-cat <port>
```

**Exemple :**
```bash
./net-cat 8080
```

Si aucun port n'est spécifié, le serveur écoute par défaut sur le port **8989** :
```bash
./net-cat
# Écoute sur localhost:8989
```

### Connecter un client

Pour se connecter à un serveur de chat existant :

```bash
./net-cat <adresse_ip> <port>
```

**Exemple :**
```bash
./net-cat localhost 8080
```

### Flux d'utilisation typique

1. **Démarrage du serveur :**
   ```bash
   Terminal 1 : ./net-cat 8989
   # Output : Listening on the port : 8989
   ```

2. **Connexion des clients :**
   ```bash
   Terminal 2 : ./net-cat localhost 8989
   # Le client voit le message de bienvenue ASCII art
   # Il entre son nom d'utilisateur
   ```

3. **Communication :**
   - Les messages sont précédés du timestamp et du nom d'utilisateur : `[2025-10-25 15:20:00][Alice]:`
   - Les notifications d'entrée/sortie sont affichées en orange
   - Les messages sont affichés en bleu
   - Les erreurs sont affichées en rouge

## 🏗️ Architecture technique

### Modèle Client-Serveur

Le projet suit une architecture client-serveur classique basée sur TCP :

```
┌─────────────┐         TCP          ┌─────────────┐
│   Client 1  │◄─────────────────────►│             │
└─────────────┘                       │             │
                                      │   Serveur   │
┌─────────────┐         TCP          │   net-cat   │
│   Client 2  │◄─────────────────────►│             │
└─────────────┘                       │             │
                                      │             │
┌─────────────┐         TCP          │             │
│   Client N  │◄─────────────────────►│             │
└─────────────┘                       └─────────────┘
                                            │
                                            │
                                      ┌─────▼─────┐
                                      │ msglogs.  │
                                      │   json    │
                                      └───────────┘
```

### Composants principaux

1. **Server** : Gère les connexions entrantes et la distribution des messages
2. **Client** : Établit des connexions et envoie/reçoit des messages
3. **Message Handler** : Encode/décode les messages au format JSON
4. **Validation** : Vérifie les noms d'utilisateur et les entrées

## 🧠 Stratégies et algorithmes

### 1. Gestion des connexions concurrentes

**Stratégie :** Utilisation de goroutines pour gérer plusieurs clients simultanément

```go
// Chaque connexion client utilise 3 goroutines :
go s.ShowLogin(conn)      // Affichage du message de bienvenue
go s.readLoop(conn)       // Lecture des messages du client
go s.printLoop(conn)      // Distribution des messages aux autres clients
```

**Avantages :**
- Non-bloquant : Un client lent n'affecte pas les autres
- Scalabilité : Peut gérer jusqu'à 10 clients sans dégradation
- Parallélisme natif de Go

### 2. Système de diffusion (Broadcasting)

**Algorithme :**
```
Pour chaque message reçu :
  1. Encoder le message en JSON
  2. Ajouter au log des messages
  3. Pour chaque client connecté :
     - Si client ≠ auteur du message :
       - Envoyer le message au client
```

**Implémentation :**
```go
func (s *Server) BroadcastMsg(msg []byte, excluded string) {
    for conn, usr := range s.clients {
        if usr != excluded {
            conn.Write([]byte(msg))
        }
    }
}
```

**Complexité :** O(n) où n = nombre de clients connectés

### 3. Validation des noms d'utilisateur

**Méthode :** Validation alphanumerique stricte

```go
func IsAlphaNumeric(s string) bool {
    for _, r := range s {
        if r == ' ' || !((r >= 'a' && r <= 'z') || 
                         (r >= 'A' && r <= 'Z') || 
                         (r >= '0' && r <= '9')) {
            return false
        }
    }
    return true
}
```

**Contraintes :**
- Caractères autorisés : a-z, A-Z, 0-9
- Pas d'espaces ni de caractères spéciaux
- Unicité garantie par la map `ExistingUsers`

### 4. Gestion de l'historique des messages

**Stratégie :** Persistance JSON pour la récupération des logs

```
Nouveau client se connecte :
  1. Vérification de l'existence de msglogs.json
  2. Si logs existent :
     - Envoyer notification "logs" au client
     - Client lit et affiche l'historique
  3. Client continue à recevoir les nouveaux messages
```

**Structure du message :**
```json
{
  "Type": "msg|notif|error|logs",
  "Author": "nom_utilisateur",
  "Text": "contenu_du_message",
  "Date": "2025-10-25 15:20:00"
}
```

### 5. Validation des entrées réseau

**Méthode pour les adresses IP :**
```go
// Utilisation de regex pour valider les octets IPv4
oct := `([1-9]|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])`
// Accepte aussi "localhost" comme alias
```

**Méthode pour les ports :**
```go
// Validation numérique simple (0-9 uniquement)
func IsPort(s string) bool {
    for _, r := range s {
        if r < '0' || r > '9' {
            return false
        }
    }
    return true
}
```

### 6. Coloration des messages

**Stratégie :** Codes ANSI pour différencier visuellement les types de messages

```go
Orange = ColorAnsiStart(255, 94, 0)   // Notifications
Red = ColorAnsiStart(255, 0, 0)       // Erreurs
Blue = ColorAnsiStart(0, 60, 255)     // Messages normaux
```

**Formule RGB vers ANSI :**
```go
fmt.Sprintf("\033[38;2;%d;%d;%dm", R, G, B)
```

### 7. Gestion de la déconnexion

**Algorithme :**
```
Lors de la détection d'EOF :
  1. Récupérer le nom du client déconnecté
  2. Créer un message de notification
  3. Broadcaster la notification aux autres clients
  4. Supprimer le client de la map des clients
  5. Marquer le nom d'utilisateur comme disponible
  6. Fermer la connexion TCP
```

## 📁 Structure du code

```
net-cat/
├── main.go                    # Point d'entrée de l'application
├── go.mod                     # Fichier de dépendances Go
├── welcome-text.txt           # Message de bienvenue ASCII art
├── msglogs.json              # Historique des messages (généré)
└── handlers/                  # Package contenant la logique métier
    ├── run.go                # Gestion des arguments et démarrage
    ├── server.go             # Logique du serveur TCP
    ├── client.go             # Logique du client TCP
    └── regex.go              # Fonctions de validation
```

### Description des fichiers

#### `main.go`
Point d'entrée minimaliste qui délègue l'exécution au package handlers.

#### `handlers/run.go`
- Analyse les arguments de ligne de commande
- Détermine si l'utilisateur veut créer un serveur ou se connecter
- Valide les paramètres (port, IP)

#### `handlers/server.go`
Contient la logique complète du serveur :
- **`Server` struct** : Représente un serveur avec son listener, channels et clients
- **`Start()`** : Lance le serveur et commence à écouter
- **`acceptLoop()`** : Boucle d'acceptation des nouvelles connexions
- **`readLoop()`** : Lit les messages d'un client spécifique
- **`printLoop()`** : Distribue les messages à tous les clients
- **`BroadcastMsg()`** : Envoie un message à tous les clients sauf l'auteur
- **`AddClient()`** : Ajoute un nouveau client avec validation

#### `handlers/client.go`
Gère le côté client :
- **`Msg` struct** : Structure des messages JSON
- **`ConnectNewUser()`** : Établit la connexion au serveur
- **`SendMsg()`** : Envoie les messages du client au serveur
- **`UserMessages()`** : Reçoit et affiche les messages du serveur
- **`AskUserName()`** : Demande et valide le nom d'utilisateur
- Fonctions utilitaires : `EncodeMsg()`, `DecodeMsg()`, `IsAlphaNumeric()`, etc.

#### `handlers/regex.go`
Fonctions de validation :
- **`IsPort()`** : Vérifie si la chaîne est un port valide
- **`IsIP()`** : Valide les adresses IPv4 et "localhost"

## ✨ Fonctionnalités

### Fonctionnalités du serveur

1. **Gestion multi-clients**
   - Jusqu'à 10 utilisateurs simultanés
   - Refus automatique si le serveur est plein

2. **Broadcasting intelligent**
   - Distribution des messages à tous les clients sauf l'émetteur
   - Messages d'état (connexion/déconnexion) envoyés à tous

3. **Persistance des messages**
   - Sauvegarde automatique dans `msglogs.json`
   - Restauration de l'historique pour les nouveaux arrivants

4. **Messages de bienvenue**
   - Affichage d'ASCII art personnalisé
   - Configurable via `welcome-text.txt`

5. **Gestion des erreurs**
   - Noms d'utilisateur en double détectés
   - Messages d'erreur clairs pour le client

### Fonctionnalités du client

1. **Interface utilisateur interactive**
   - Affichage en temps réel des messages
   - Curseur repositionné automatiquement
   - Timestamp et nom sur chaque message

2. **Validation des entrées**
   - Noms alphanumériques uniquement
   - Messages lisibles uniquement (pas de caractères de contrôle)

3. **Gestion de l'historique**
   - Téléchargement automatique des messages précédents
   - Affichage propre avec préservation du curseur

4. **Colorisation**
   - Messages en bleu
   - Notifications en orange
   - Erreurs en rouge

## 🛠️ Technologies utilisées

### Langage
- **Go (Golang) 1.21.5**
  - Goroutines pour la concurrence
  - Channels pour la communication inter-goroutines
  - Package `net` pour les connexions TCP
  - Package `encoding/json` pour la sérialisation

### Protocole réseau
- **TCP (Transmission Control Protocol)**
  - Garantit la livraison des messages
  - Connexion orientée
  - Communication bidirectionnelle

### Bibliothèques standard Go

```go
import (
    "net"              // Connexions TCP
    "encoding/json"    // Sérialisation des messages
    "bufio"            // Lecture buffered des entrées
    "regexp"           // Validation par expressions régulières
    "time"             // Horodatage des messages
    "os"               // Opérations système
)
```

### Formats de données
- **JSON** : Stockage et transfert des messages structurés
- **ANSI escape codes** : Coloration du terminal

## 💡 Processus de développement

### Méthodologie

1. **Conception initiale**
   - Définition du protocole de communication
   - Choix de la structure JSON pour les messages
   - Architecture client-serveur classique

2. **Développement itératif**
   - Implémentation du serveur de base
   - Ajout du client et connexion
   - Ajout des fonctionnalités (historique, couleurs, validation)

3. **Optimisations**
   - Utilisation de goroutines pour la concurrence
   - Channels pour éviter les race conditions
   - Buffer de 4096 octets pour la lecture efficace

### Défis techniques résolus

1. **Synchronisation des messages**
   - Solution : Channel `msgch` avec buffer de 10 messages
   - Évite les blocages lors de pics de trafic

2. **Affichage concurrent**
   - Solution : Codes ANSI pour repositionner le curseur
   - Préserve la ligne de saisie utilisateur

3. **Gestion de la mémoire**
   - Solution : Fermeture propre des connexions
   - Suppression des clients déconnectés des structures

4. **Unicité des noms**
   - Solution : Map `ExistingUsers` partagée
   - Vérification avant l'ajout d'un nouveau client

## 🤝 Contribution

Les contributions sont les bienvenues ! Si vous souhaitez contribuer à ce projet :

1. **Fork** le dépôt
2. **Créez** une branche pour votre fonctionnalité (`git checkout -b feature/AmazingFeature`)
3. **Committez** vos changements (`git commit -m 'Add some AmazingFeature'`)
4. **Push** vers la branche (`git push origin feature/AmazingFeature`)
5. **Ouvrez** une Pull Request

### Suggestions d'améliorations

- Authentification par mot de passe
- Salons de discussion multiples
- Messages privés entre utilisateurs
- Interface graphique (GUI)
- Chiffrement des communications
- Support IPv6
- Compression des messages

## 📄 Licence

Ce projet est distribué sous licence MIT. Voir le fichier `LICENSE` pour plus d'informations.

---

**Développé avec ❤️ en Go**