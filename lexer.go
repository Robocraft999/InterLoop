package main

import "strconv"

type Token = int

const (
	EOF Token = iota
	PLUS
	MINUS
	ASSIGN
	SEMICOLON
	LOOP
	DO
	END
	IDENT
	NUM
)

func Lex(input string) ([]Token, []string, []int) {
	var tokens = make([]Token, 0)
	var currentRaw = ""
	var idents = make([]string, 0)
	var numbers = make([]int, 0)
	for _, char := range input {
		if char != ' ' && char != '\t' && char != '\n' {
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
			case ";":
				//tokens = append(tokens, SEMICOLON)
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
					idents = append(idents, currentRaw)
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
	return tokens, idents, numbers
}
