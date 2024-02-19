# net-cat

"net-cat" est une version personnalisée du net-cat de Linux, écrite en Go. Il s'agit d'un outil de réseau qui permet aux utilisateurs de créer des serveurs TCP et de se connecter à des serveurs TCP distants.

## Installation

Pour installer "net-cat", vous devez avoir Go installé sur votre système. Vous pouvez alors cloner ce dépôt et construire l'exécutable avec la commande suivante :
```
bash go build -o net-cat
```


## Utilisation

### Créer un serveur

Pour créer un serveur, exécutez "net-cat" avec le numéro de port sur lequel vous souhaitez que le serveur écoute. Par défaut, le serveur écoute sur le port  8989 si aucun port n'est spécifié. Le programme fonctionne avec le localhost.
```
./net-cat 8989
```


## Contribution

Les contributions sont les bienvenues. Si vous souhaitez contribuer à ce projet, veuillez soumettre une demande de tirage (pull request) avec vos modifications.

## Licence

"net-cat" est distribué sous la licence MIT. Voir le fichier LICENSE pour plus d'informations.
