package main

import "strconv"

type Token int

const (
	EOF Token = iota
	PLUS
	MINUS
	ASSIGN
	LOOP
	DO
	END
	IDENT
	NUM
)

func Lex(input string) ([]Token, int, []int, []int) {
	input += "\n"
	var tokens = make([]Token, 0)
	var currentRaw = ""
	var uniqueIdents = make(map[string]int)
	var idents = make([]int, 0)
	var numbers = make([]int, 0)
	for _, char := range input {
		if char != ' ' && char != '\t' && char != '\n' && char != ';' {
			currentRaw += string(char)

			switch currentRaw {
			case "":
				tokens = append(tokens, EOF)
				currentRaw = ""
			case "+":
				tokens = append(tokens, PLUS)
				currentRaw = ""
			case "-":
				tokens = append(tokens, MINUS)
				currentRaw = ""
			case ":=":
				tokens = append(tokens, ASSIGN)
				currentRaw = ""
			case "LOOP":
				tokens = append(tokens, LOOP)
				currentRaw = ""
			case "DO":
				tokens = append(tokens, DO)
				currentRaw = ""
			case "END":
				tokens = append(tokens, END)
				currentRaw = ""
			}
		} else {
			if currentRaw != "" {
				var num, err = strconv.Atoi(currentRaw)
				if err != nil {
					if _, ok := uniqueIdents[currentRaw]; !ok {
						uniqueIdents[currentRaw] = len(uniqueIdents)
					}
					idents = append(idents, uniqueIdents[currentRaw])
					tokens = append(tokens, IDENT)
				} else {
					numbers = append(numbers, num)
					tokens = append(tokens, NUM)
				}
			}
			currentRaw = ""
		}
	}
	tokens = append(tokens, EOF)
	return tokens, len(uniqueIdents), idents, numbers
}

func (t Token) String() string {
	return [...]string{"EOF", "PLUS", "MINUS", "ASSIGN", "LOOP", "DO", "END", "IDENT", "NUM"}[t]
}
