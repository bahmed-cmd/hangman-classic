Jeu du Pendu en Go

Bienvenue dans le jeu du Pendu implémenté en Go ! Cette version vous permet de deviner des mots tout en sauvegardant votre progression et vous offre la flexibilité de quitter ou de sauvegarder à tout moment.

Fonctionnalités

Jouez au jeu du Pendu avec 10 tentatives.

Sauvegardez en utilisant la commande STOP.

Reprenez la partie  en utilisant l'état sauvegardé.

Quittez la partie avec la commande QUIT.

Prérequis

Go (version 1.15 ou plus récente)

Installation

Clonez ce dépôt :

git clone https://github.com/bahmed-cmd/hangman-classic.git
cd hangman-game

Assurez-vous d'avoir les fichiers suivants :

words.txt : Contient la liste des mots .

hangman.txt : Contient  la progression du pendu.

Utilisation

Démarrer une Nouvelle Partie

Exécutez la commande suivante pour démarrer une nouvelle partie :

go run main.go words.txt

Vous aurez 10 tentatives pour deviner le mot.

Sauvegarder et Reprendre la Partie

Pendant la partie, tapez STOP pour sauvegarder votre progression. Le jeu sera sauvegardé dans save.txt.

Pour reprendre une partie sauvegardée, exécutez :

go run main.go --startWith save.txt

Résumé des Commandes

STOP : Sauvegarder votre progression et quitter la partie.

QUIT : Quitter la partie sans sauvegarder.

Exemple

$ go run main.go words.txt
Bonne chance, vous avez 10 tentatives.
Mot actuel: _ _ _ _ _
Tentatives restantes: 10
Choisissez une lettre, un mot, ou tapez 'STOP' pour sauvegarder ou 'QUIT' pour quitter : A
Lettre incorrecte ! Vous perdez 1 tentative.

