package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
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
}

// choisir le mot aleatoire de fichier words.txt
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
		word := scanner.Text()
		count++
		if rand.Intn(count) == 0 {
			chosenWord = word
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if count == 0 {
		return "", fmt.Errorf("no words found in the file")
	}
	return chosenWord, nil
}

// avancer la construction de pendu
func parseHangmanPositions(hangman string) ([10]string, error) {
	var positions [10]string
	data, err := ioutil.ReadFile(hangman)
	if err != nil {
		return positions, err
	}
	lines := strings.Split(string(data), "\n")
	for i := 0; i < 10; i++ {
		positions[i] = strings.Join(lines[i*8:(i+1)*8], "\n")
	}
	return positions, nil
}

// afficher des lettres aleatoires
func revealLetters(words string) string {
	revealed := make([]rune, len(words))
	for i := range revealed {
		revealed[i] = '_'
	}
	n := len(words)/2 - 1
	randIndices := rand.Perm(len(words))[:n]
	for _, i := range randIndices {
		revealed[i] = rune(words[i])
	}
	return string(revealed)
}

// logique de jeu
func playHangman(hangman HangManData) {
	reader := bufio.NewReader(os.Stdin)

	for hangman.Attempts > 0 {
		fmt.Println("Current word:", hangman.Word)
		fmt.Printf("Attempts left: %d\n", hangman.Attempts)
		fmt.Println(hangman.HangmanPositions[10-hangman.Attempts])

		fmt.Print("Choose a letter: ")
		input, _ := reader.ReadString('\n')
		guess := strings.TrimSpace(input)

		if len(guess) != 1 {
			fmt.Println("Please enter a single letter.")
			continue
		}

		guessRune := rune(guess[0])
		if strings.ContainsRune(hangman.ToFind, guessRune) {
			for i, c := range hangman.ToFind {
				if c == guessRune {
					hangman.Word = hangman.Word[:i] + string(c) + hangman.Word[i+1:]
				}
			}
			if hangman.Word == hangman.ToFind {
				fmt.Println("Congrats! You've found the word:", hangman.ToFind)
				return
			}
		} else {
			hangman.Attempts--
			fmt.Println("Wrong guess!")
		}
	}
	fmt.Println("Game Over! The word was:", hangman.ToFind)
	fmt.Println(hangman.HangmanPositions[0])
}

// condition de jeux fonction main
func main() {
	word, err := getRandomWord("words.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Random word:", word)

	positions, err := parseHangmanPositions("hangman.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	hangman := HangManData{
		Word:             revealLetters(word),
		ToFind:           word,
		Attempts:         10,
		HangmanPositions: positions,
	}

	fmt.Println("Good Luck, you have 10 attempts.")
	playHangman(hangman)
}
