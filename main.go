package main

import (
	"fmt"
	"strings"
	// "log"
	// "encoding/json"
)

type production struct {
	left  string
	right []string
}

func getNonTerminals(productions []production) []string {
	nonTerminals := make(map[string]bool)
	for _, p := range productions {
		nonTerminals[p.left] = true
	}
	result := make([]string, 0, len(nonTerminals))
	for nt := range nonTerminals {
		result = append(result, nt)
	}
	return result
}
func getTerminals(productions []production) []string {
	terminals := make(map[string]bool)
	for _, p := range productions {
		for _, s := range p.right {
			if !isNonTerminal(s) {
				terminals[s] = true
			}
		}
	}
	result := make([]string, 0, len(terminals))
	for t := range terminals {
		result = append(result, t)
	}
	return result
}

func isNonTerminal(symbol string) bool {
	return symbol[0] >= 'A' && symbol[0] <= 'Z' && (symbol != "lambda")
}
func getProductionsForNonTerminal(left string, productions []production) []production {
	result := make([]production, 0)
	for _, p := range productions {
		if p.left == left {
			result = append(result, p)
		}
	}
	return result
}
func getGeneratedSymbols(left string, productions []production) []string {
	productionsForNonTerminal := getProductionsForNonTerminal(left, productions)
	result := make([]string, 0)
	for _, p := range productionsForNonTerminal {
		result = append(result, p.right...)
	}
	return result
}

// func getPredictions(p production, firsts map[string][]string) []string {
//     var predictions []string
//     for _, symbol := range p.right {
//         if isTerminal(symbol) {
//             predictions = append(predictions, symbol)
//             break
//         } else {
//             x2 := symbol
//             x2Firsts := getFirst(x2, firsts)
//             predictions = append(predictions, x2Firsts...)
//             if !contains(x2Firsts, "lambda") {
//                 break
//             }
//         }
//     }
//     return predictions
// }

// getFirstSet devuelve el conjunto First para el símbolo dado en la gramática dada.
// func getFirstSet(symbol string, productions []production) []string {
//     firstSet := make([]string, 0)
//     if !isNonTerminal(symbol) {
//         firstSet = append(firstSet, symbol)
//         return firstSet
//     }
//     for _, p := range productions {
//         if p.left == symbol {
//             firstSymbol := p.right[0]
//             if firstSymbol != symbol {
//                 firstSet = append(firstSet, getFirstSet(firstSymbol, productions)...)
// 				fmt.Println("Primero:", symbol,  getFirstSet(firstSymbol, productions))
//             }
//             if contains(p.right, "lambda") && len(p.right) > 1 {
//                 rest := p.right[1:]
//                 // for symbolRight := 0; symbolRight < len(p.right); symbolRight++ {
// 					// if (p.right[symbolRight] != "lambda" && !isNonTerminal(p.right[symbolRight]) ){
// 						restFirstSet := getFirstSet(strings.Join(rest, ""), productions)
// 						firstSet = append(firstSet, restFirstSet...)
// 						fmt.Println("Primero:", symbol,  getFirstSet(firstSymbol, productions))
// 						// break
// 					// }
// 				// }
//             }
//         }
//     }
//     // return removeDuplicates(firstSet)
//     return firstSet
// }

func getFirst(left string, right []string, firstSets map[string][]string, productions []production) []string {
	var first []string
	// Check if the first symbol is a terminal
	if (len(right) == 1) && isNonTerminal(right[0]) {
		for _, p := range productions {
			if left == p.left && right[0] == left {
				first = getFirst(right[0], p.right, firstSets, productions)
				// fmt.Println("p.left, p.reght", p.left, first)
			}
		}
	}
	if isTerminal(right[0]) && right[0] != "lambda" {
		first = append(first, right[0])
	} else if !isTerminal(right[0]) && right[0] != "lambda" && len(right) >= 1 {
		// Find the first non-terminal symbol
		for _, symbol := range right {
			if !isTerminal(symbol) {
				// Get the first set of the non-terminal symbol

				var x2First []string
				// recorrer producciones    --- envio el array de productions de la produccion
				for _, p := range productions {
					if left == p.left && right[0] == left {
						x2First := getFirst(right[0], p.right, firstSets, productions)
						getFirst(right[0], p.right, firstSets, productions)
						fmt.Println("Entrp")
						first = append(first, x2First...)
						// fmt.Println("p.left, p.reght", p.left, first, x2First)
					} else if left == p.left{
						x2First := getFirst(right[0], right[0:1], firstSets, productions)
						first = append(first, x2First...)
						break
					}
				}


				// Add the first set of the non-terminal symbol to the first set of the current symbol
				// first = append(first, x2First...)

				// If the first set of the non-terminal symbol does not contain lambda, stop looking for more symbols
				if !contains(x2First, "lambda") {
					break
				}
			} else {
				// If the symbol is a terminal, add it to the first set of the current symbol and stop looking for more symbols
				first = append(first, symbol)
				break
			}
		}
	}

	// If all symbols in the right side of the production are lambda, or the right side is empty, add lambda to the first set
	if len(right) == 1 && right[0] == "lambda" {
		first = append(right, "lambda")
	} else if len(right) >= 1 && right[0] == "lambda" {
		for symbolRight := 0; symbolRight < len(right); symbolRight++ {
			if right[symbolRight] != "lambda" {
				if isTerminal(right[symbolRight]) {
					first = append(first, right[symbolRight])
					break
				} else if !isTerminal(right[symbolRight]) && right[symbolRight] == left {
					fmt.Println("dos", right[symbolRight], right[symbolRight:symbolRight+1])
					x2First := getFirst(right[symbolRight], right[symbolRight:symbolRight+1], firstSets, productions)
					// getFirst(right[symbolRight], right[symbolRight:symbolRight+1], firstSets, productions)
					first = append(first, x2First...)
					break
				} else if !isTerminal(right[symbolRight]){
					x2First := getFirst(right[symbolRight], right[symbolRight:symbolRight+1], firstSets, productions)
					first = append(first, x2First...)
					break
				}
			}
		}
	}

	// Add the first set of the current symbol to the first sets map
	firstSets[left] = append(firstSets[left], first...)

	return first
}

func isTerminal(symbol string) bool {
	// Check if the symbol is lowercase or an empty string
	return symbol == strings.ToLower(symbol) || symbol == "lambda"
}

// contains devuelve true si la lista dada contiene el elemento dado, false en caso contrario.
func contains(list []string, element string) bool {
	for _, e := range list {
		if e == element {
			return true
		}
	}
	return false
}

// removeDuplicates devuelve una copia de la lista dada sin elementos duplicados.
func removeDuplicates(list []string) []string {
	uniqueList := make([]string, 0)
	seen := make(map[string]bool)
	for _, element := range list {
		if !seen[element] {
			uniqueList = append(uniqueList, element)
			seen[element] = true
		}
	}
	return uniqueList
}
func main() {
	productions := []production{
		// {"E", []string{"lambda","T", "EP"}},
		// {"EP", []string{"+", "T", "EP"}},
		// {"EP", []string{"lambda"}},
		// {"T", []string{"lambda","lambda", "F", "TP"}},
		// // {"T", []string{"M", "TP"}},
		// {"TP", []string{"lambda","*", "F", "TP"}},
		// {"TP", []string{"lambda"}},
		// {"F", []string{"id"}},
		// {"F", []string{"(", "E", ")"}},
		// // {"M", []string{"(", "TP"}},

		{"AL", []string{"id", ":=", "P"}},
		{"P", []string{"D", "PP"}},
		{"PP", []string{"or","D", "PP"}},
		{"PP", []string{"lambda"}},
		{"D", []string{"C", "DP"}},
		{"DP", []string{"and","C","DP"}},
		{"DP", []string{"lambda"}},
		{"C", []string{"S"}},
		{"C", []string{"not", "(", "P", ")"}},
		{"S", []string{"(", "P", ")"}},
		{"S", []string{"OP", "REL", "OP"}},
		{"S", []string{"true"}},
		{"S", []string{"false"}},
		{"REL", []string{"=",}},
		{"REL", []string{"<","RP"}},
		{"REL", []string{">", "EP"}},
		{"RP", []string{"="}},
		{"RP", []string{">"}},
		{"RP", []string{"lambda"}},
		{"EP", []string{"="}},
		{"EP", []string{"lambda"}},
		{"OP", []string{"id"}},
		{"OP", []string{"num"}},

	}

	fmt.Println("Productions:")
	for _, p := range productions {
		fmt.Printf("%s -> %v\n", p.left, p.right)
	}

	firstSets := make(map[string][]string)
	fmt.Println("Primeros:")
	for _, p := range productions {
		// fmt.Println(p.left, p.right)
		// first := getFirst(p.left, p.right, firstSets, productions)
		getFirst(p.left, p.right, firstSets, productions)
		// fmt.Println("p.left", p.left, "p.reght", p.right, first)
	}

	// Imprime los primeros de cada no terminal
	for key, value := range firstSets {
		fmt.Printf("FIRST(%s): %v\n", key, value)
	}

	// predictions := getPredictions(productions)
	// fmt.Println("Predictions:")
	// for nt, preds := range predictions {
	//     fmt.Printf("%s -> %v\n", nt, preds)
	// }

	// if isLL1(productions) {
	// 	fmt.Println("La gramática es una gramática LL(1).")
	// } else {
	// 	fmt.Println("La gramática no es una gramática LL(1).")
	// }

}

// isLL1 devuelve true si la gramática dada es una gramática LL(1), false en caso contrario.
// func isLL1(productions []production) bool {
//     predictions := getPredictions(productions)
//     for nt, preds := range predictions {
//         // Si hay una producción que comienza con lambda y otro símbolo, la gramática no es LL(1).
//         if contains(preds, "lambda") && len(preds) > 1 {
//             return false
//         }
//         // Si hay dos producciones que comienzan con el mismo símbolo de entrada, la gramática no es LL(1).
//         firstSets := make(map[string][]string)
//         for _, p := range productions {
//             if p.left == nt {
//                 firstSets[strings.Join(p.right, "")] = getFirstSet(strings.Join(p.right, ""), productions)
//             }
//         }
//         for i := 0; i < len(preds); i++ {
//             for j := i + 1; j < len(preds); j++ {
//                 if contains(firstSets[preds[i]], preds[j]) || contains(firstSets[preds[j]], preds[i]) {
//                     return false
//                 }
//             }
//         }
//     }
//     return true
// }
