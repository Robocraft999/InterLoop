package main

import "fmt"

type Interpreter struct {
	tokens       []Token
	idents       []string
	identCount   int
	numbers      []int
	numbersCount int
	vars         map[string]int
	index        int
}

func NewInterpreter(tokens []Token, indents []string, numbers []int) *Interpreter {
	return &Interpreter{
		tokens:       tokens,
		idents:       indents,
		identCount:   0,
		numbers:      numbers,
		numbersCount: 0,
		vars:         make(map[string]int),
		index:        0,
	}
}

func (i *Interpreter) Interpret() {
	for i.current() != EOF {
		i.interpretStatement()
	}
}

func (i *Interpreter) interpretStatement() {
	var current = i.advance()
	fmt.Println(i.index, current)
	switch current {
	case LOOP:
		{
			var identToken = i.advance()
			if identToken != IDENT {
				panic("Expected IDENT in LOOP head")
			}
			var doToken = i.advance()
			if doToken != DO {
				panic("Expected DO in LOOP head")
			}

			var loopAmountVar = i.idents[i.identCount]
			var loopAmount = i.vars[loopAmountVar]
			i.identCount++

			if loopAmount == 0 {
				i.jumpToEnd()
				return
			}

			if loopAmount > 1 {
				var currentPc = i.index
				var currentIdentCount = i.identCount
				var currentNumbersCount = i.numbersCount

				for range loopAmount - 1 {
					i.interpretStatement()
					i.index = currentPc
					i.identCount = currentIdentCount
					i.numbersCount = currentNumbersCount
				}
			}
			i.interpretStatement()
			var endToken = i.advance()
			if endToken != END {
				panic("Expected END in LOOP tail")
			}
		}
	case IDENT:
		{
			var currentIdentIndex = i.identCount
			i.identCount++
			var assignToken = i.advance()
			if assignToken != ASSIGN {
				panic("Expected ASSIGN in statement")
			}
			var otherVarToken = i.advance()
			if otherVarToken != IDENT {
				panic("Expected IDENT in statement")
			}
			var otherIdentIndex = i.identCount
			i.identCount++

			var operationToken = i.advance()

			var numberToken = i.advance()
			if numberToken != NUM {
				panic("Expected NUM in statement")
			}
			var numberIndex = i.numbersCount
			i.numbersCount++
			var number = i.numbers[numberIndex]

			if operationToken == PLUS {
				i.vars[i.idents[currentIdentIndex]] = i.vars[i.idents[otherIdentIndex]] + number
			} else if operationToken == MINUS {
				i.vars[i.idents[currentIdentIndex]] = i.vars[i.idents[otherIdentIndex]] - number
			} else {
				panic("Expected PLUS or MINUS in statement")
			}

		}
	default:
		{
			panic("unreachable")
		}
	}
}

func (i *Interpreter) advance() Token {
	if i.index >= len(i.tokens) {
		return EOF
	}
	var current = i.tokens[i.index]
	i.index++
	return current
}

func (i *Interpreter) current() Token {
	return i.tokens[i.index]
}

func (i *Interpreter) peek() Token {
	if i.index+1 >= len(i.tokens) {
		return EOF
	}
	return i.tokens[i.index+1]
}

func (i *Interpreter) jumpToEnd() {
	var currentIndex = i.index
	var count = 1
	j := currentIndex
	for ; count > 0; j++ {
		var tok = i.tokens[j]
		if tok == LOOP {
			count++
		} else if tok == END {
			count--
		}
	}
	i.index = j + 1
}
