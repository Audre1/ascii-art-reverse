
package main
import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
)

func main() {
	// Analyser les arguments de la ligne de commande
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		// Effectuer l'opération inverse s'il n'y a pas d'arguments supplémentaires
		reverse(args)
	} else {
		fmt.Println("Trop d'arguments")
	}
}

// Définir un drapeau pour spécifier le fichier d'entrée pour l'opération inverse
var readFlags = flag.String("reverse", "example.txt", "lire le fichier à partir du drapeau")

// La fonction reverse gère les vérifications d'entrée et effectue l'opération inverse.
func reverse(args []string) {
	// Vérifier les erreurs d'utilisation potentielles
	checkForAudit()

	// Définir le chemin vers le fichier de polices ASCII
	fonts := "standard.txt"

	// Définir les instructions d'utilisation
	const usage = "Usage : go run . [OPTION]\n\nEX : go run . --reverse=<fileName>"

	// Vérifier si le drapeau --reverse est fourni avec un seul argument
	if !strings.Contains(*readFlags, "--reverse=") && len(args) == 1 {
		fmt.Println(args[0])
		os.Exit(0)
	}

	// Vérifier la présence d'arguments supplémentaires
	if len(args) > 0 {
		fmt.Println(usage)
		return
	}

	// Lire le contenu du fichier d'entrée
	input, err := os.ReadFile(*readFlags)
	if err != nil {
		fmt.Printf("Impossible de lire le contenu du fichier en raison de %v", err)
	}

	// Diviser l'entrée en lignes
	matrix := strings.Split(string(input), "\n")

	// Supprimer les signes de dollar en fin de chaque ligne
	matrix2 := delDollarSigns(matrix)

	// Trouver les colonnes vides (espaces) dans l'entrée utilisateur
	spaces := findSpace(matrix2)

	// Diviser l'entrée utilisateur en fonction des colonnes vides
	userInput := splitUserInput(matrix2, spaces)

	// Mapper la saisie utilisateur divisée
	userInputMap := userInputMapping(userInput)

	// Obtenir les polices graphiques ASCII
	asciiGraphic := getASCIIgraphicFont(fonts)

	// Associer la saisie utilisateur avec les polices graphiques ASCII et générer la sortie
	output := mapUserInputWithASCIIgraphicFont(userInputMap, asciiGraphic)

	// Afficher la sortie générée
	fmt.Println(output)
}

// delDollarSigns supprime les signes de dollar en fin de chaque ligne dans la matrice.
func delDollarSigns(matrix []string) []string {
	var matrix2 []string
	for _, v := range matrix {
		lenv := len(v)
		if lenv <= 1 {
			matrix2 = append(matrix2, "")
		} else {
			matrix2 = append(matrix2, v[:lenv-1])
		}
	}
	return matrix2
}

// findSpace identifie les colonnes vides dans la matrice.
func findSpace(matrix []string) []int {
	var emptyColumns []int
	count := 0

	for column := 0; column < len(matrix[0]); column++ {
		for row := 0; row < len(matrix)-1; row++ {
			if matrix[row][column] == ' ' {
				count++
			} else {
				count = 0
				break
			}

			if count == len(matrix)-1 {
				emptyColumns = append(emptyColumns, column)
				count = 0
			}
		}
	}

	// Vérifier les espaces supplémentaires et les convertir en conséquence
	count = 5
	var indexToRem []int
	for i := range emptyColumns {
		if count == 0 {
			count = 5
			continue
		}
		if i > 0 {
			if emptyColumns[i] == (emptyColumns[i-1])+1 {
				indexToRem = append(indexToRem, i)
				count -= 1
			}
		}
	}

	// Supprimer les espaces supplémentaires
	for i := len(indexToRem) - 1; i >= 0; i-- {
		emptyColumns = removeIndex(emptyColumns, indexToRem[i])
	}

	return emptyColumns
}

// removeIndex supprime un élément à un index donné d'une tranche.
func removeIndex(s []int, index int) []int {
	if index < 0 || index >= len(s) {
		return s
	}
	return append(s[:index], s[index+1:]...)
}

// checkForAudit vérifie les erreurs d'utilisation potentielles.
func checkForAudit() {
	if strings.Contains(os.Args[1], "--") && !strings.Contains(os.Args[1], "=") {
		fmt.Println("Usage : go run . [OPTION]\n\nEX : go run . --reverse=<fileName>")
		os.Exit(0)
	}
}

// splitUserInput divise l'entrée utilisateur en fonction des colonnes vides.
func splitUserInput(matrix []string, emptyColumns []int) string {
	var result string
	result = "\n"
	start := 0
	end := 0

	for _, column := range emptyColumns {
		if end < len(matrix[0]) {
			end = column

			for _, characters := range matrix {
				if len(characters) > 0 {
					columns := characters[start:end]
					result = result + columns + " "
				}
				result = result + "\n"
			}

			start = end + 1
		}
	}

	return result
}

// userInputMapping mappe l'entrée utilisateur pour la recherche.
func userInputMapping(result string) map[int][]string {
	strSlice := strings.Split(result, "\n")
	graphicInput := make(map[int][]string)
	j := 0

	for _, ch := range strSlice {
		if ch == "" {
			j++
		} else {
			graphicInput[j] = append(graphicInput[j], ch)
		}
	}

	return graphicInput
}

// getASCIIgraphicFont lit les polices graphiques ASCII à partir d'un fichier.
func getASCIIgraphicFont(fonts string) map[int][]string {
	readFile, err := os.ReadFile(fonts)
	if err != nil {
		fmt.Printf("Impossible de lire le contenu du fichier en raison de %v", err)
	}

	slice := strings.Split(string(readFile), "\n")
	ascii := make(map[int][]string)
	i := 31

	for _, ch := range slice {
		if ch == "" {
			i++
		} else {
			ascii[i] = append(ascii[i], ch)
		}
	}

	return ascii
}

// mapUserInputWithASCIIgraphicFont associe l'entrée utilisateur avec les polices graphiques ASCII et renvoie la sortie.
func mapUserInputWithASCIIgraphicFont(graphicInput, ascii map[int][]string) string {
	var keys []int
	for k := range graphicInput {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	var output string
	var sliceOfBytes []byte

	for _, value := range keys {
		graphicValue := graphicInput[value]
		for asciiKey, asciiValue := range ascii {
			if reflect.DeepEqual(asciiValue, graphicValue) {
				sliceOfBytes = append(sliceOfBytes, byte(asciiKey))
			}
		}
		output = string(sliceOfBytes)
	}

	return output
}