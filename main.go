package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type HangManData struct {
	Word             string
	ToFind           string
	Attempts         int
	HangmanPositions [10]string
	GuessedLetters   map[string]bool
}

// Choisir un mot aléatoire de fichier words.txt
func getRandomWord(words string) (string, error) {
	file, err := os.Open(words)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rand.Seed(time.Now().UnixNano())

	var chosenWord string
	count := 0

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		count++
		if rand.Intn(count) == 0 {
			chosenWord = word
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if count == 0 {
		return "", fmt.Errorf("Aucun mot trouvé dans le fichier")
	}

	return chosenWord, nil
}

// Construction du pendu
func parseHangmanPositions(hangman string) ([10]string, error) {
	var positions [10]string
	data, err := os.ReadFile(hangman)
	if err != nil {
		return positions, err
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) < 80 {
		return positions, fmt.Errorf("Le fichier hangman.txt doit contenir au moins 80 lignes")
	}
	for i := 0; i < 10; i++ {
		positions[i] = strings.Join(lines[i*8:(i+1)*8], "\n")
	}
	return positions, nil
}

// Révéler des lettres aléatoires du mot choisi
func revealLetters(word string) string {
	revealed := make([]rune, len(word))
	for i := range revealed {
		if word[i] == ' ' {
			revealed[i] = ' '
		} else {
			revealed[i] = '_'
		}
	}

	n := len(word)/2 - 1
	randIndices := rand.Perm(len(word))[:n]
	for _, i := range randIndices {
		revealed[i] = rune(word[i])
	}

	return string(revealed)
}

// Logique du jeu
func playHangman(hangman *HangManData) {
	reader := bufio.NewReader(os.Stdin)

	for hangman.Attempts > 0 {
		fmt.Println("Mot actuel:", hangman.Word)
		fmt.Printf("Tentatives restantes: %d\n", hangman.Attempts)
		fmt.Println(hangman.HangmanPositions[10-hangman.Attempts])

		fmt.Print("Choisissez une lettre, un mot, ou tapez 'STOP' pour sauvegarder ou 'QUIT' pour quitter : ")
		input, _ := reader.ReadString('\n')
		guess := strings.TrimSpace(input)

		if len(guess) == 0 {
			fmt.Println("Entrée invalide. Veuillez essayer de nouveau.")
			continue
		}

		if strings.EqualFold(guess, "STOP") {
			saveGame(hangman)
			fmt.Println("Jeu sauvegardé dans save.txt.")
			return
		}

		if strings.EqualFold(guess, "QUIT") {
			fmt.Println("Vous avez quitté le jeu.")
			return
		}

		handleGuess(hangman, guess)

		if hangman.Word == hangman.ToFind {
			fmt.Println("Félicitations ! Vous avez trouvé le mot :", hangman.ToFind)
			return
		}
	}

	fmt.Println("Fin du jeu ! Le mot était :", hangman.ToFind)
}

// Gérer la proposition du joueur
func handleGuess(hangman *HangManData, guess string) {
	if hangman.GuessedLetters[guess] {
		fmt.Println("Vous avez déjà proposé cela, essayez autre chose.")
		return
	}

	hangman.GuessedLetters[guess] = true

	// Vérifier le mot choisi
	if len(guess) > 1 {
		if strings.EqualFold(guess, hangman.ToFind) {
			hangman.Word = hangman.ToFind
			fmt.Println("Félicitations ! Vous avez trouvé le mot :", hangman.ToFind)
			hangman.Attempts = 0
		} else {
			fmt.Println("Mot incorrect ! Vous perdez 2 tentatives.")
			hangman.Attempts -= 2
		}
		return
	}

	// Vérifier la lettre choisie
	guessRune := rune(guess[0])
	if strings.ContainsRune(hangman.ToFind, guessRune) {
		var newWord strings.Builder
		for i, c := range hangman.ToFind {
			if c == ' ' || c == guessRune || hangman.Word[i] != '_' {
				newWord.WriteRune(c)
			} else {
				newWord.WriteRune('_')
			}
		}
		hangman.Word = newWord.String()
	} else {
		fmt.Println("Lettre incorrecte ! Vous perdez 1 tentative.")
		hangman.Attempts--
	}
}

// Sauvegarder le jeu dans un fichier JSON
func saveGame(hangman *HangManData) {
	data, err := json.Marshal(hangman)
	if err != nil {
		fmt.Println("Erreur lors de la sauvegarde :", err)
		return
	}

	err = os.WriteFile("save.txt", data, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier :", err)
	}
}

// Charger le jeu depuis un fichier JSON
func loadGame(filename string) (*HangManData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var hangman HangManData
	err = json.Unmarshal(data, &hangman)
	if err != nil {
		return nil, err
	}

	return &hangman, nil
}

func main() {
	startWith := flag.String("startWith", "", "Charger un jeu sauvegardé à partir d'un fichier JSON.")
	flag.Parse()

	var hangman HangManData

	if *startWith != "" {
		// Charger depuis le fichier save.txt
		savedGame, err := loadGame(*startWith)
		if err != nil {
			fmt.Println("Erreur lors du chargement du jeu :", err)
			return
		}
		hangman = *savedGame
		fmt.Printf("Bienvenue de retour ! Vous avez %d tentatives restantes.\n", hangman.Attempts)
	} else {
		// Initialiser un nouveau jeu
		word, err := getRandomWord("words.txt")
		if err != nil {
			fmt.Println("Erreur :", err)
			return
		}

		positions, err := parseHangmanPositions("hangman.txt")
		if err != nil {
			fmt.Println("Erreur :", err)
			return
		}

		hangman = HangManData{
			Word:             revealLetters(word),
			ToFind:           word,
			Attempts:         10,
			HangmanPositions: positions,
			GuessedLetters:   make(map[string]bool),
		}
		fmt.Println("Bonne chance, vous avez 10 tentatives.")
	}

	// Démarrer le jeu
	playHangman(&hangman)
}
