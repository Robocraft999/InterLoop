package main

import "fmt"

type Interpreter struct {
	tokens       []Token
	idents       []int
	identCount   int
	numbers      []int
	numbersCount int
	vars         []int
	index        int
}

func NewInterpreter(tokens []Token, identsCount int, indents []int, numbers []int) *Interpreter {
	var vars = make([]int, identsCount)
	return &Interpreter{
		tokens:       tokens,
		idents:       indents,
		identCount:   0,
		numbers:      numbers,
		numbersCount: 0,
		vars:         vars,
		index:        0,
	}
}

func (i *Interpreter) Interpret() {
	i.interpretStatements()
	for x, v := range i.vars {
		fmt.Println(x, ": ", v)
	}
}

func (i *Interpreter) interpretStatements() {
	for i.current() != EOF && i.current() != END {
		i.interpretStatement()
	}
}

func (i *Interpreter) interpretStatement() {
	var current = i.advance()
	switch current {
	case LOOP:
		{

			var identToken = i.advance()
			if identToken != IDENT {
				panic("Expected IDENT in LOOP head")
			}
			var loopAmountIndex = i.identCount
			i.identCount++

			//fmt.Println("LOOP START WITH ", i.idents[loopAmountIndex], " = ", i.vars[i.idents[loopAmountIndex]])

			var doToken = i.advance()
			if doToken != DO {
				panic("Expected DO in LOOP head")
			}

			var loopAmount = i.vars[i.idents[loopAmountIndex]]

			if loopAmount == 0 {
				i.jumpToEnd()
				//fmt.Println("LOOP SKIPPED\n")
				return
			}

			if loopAmount > 1 {
				var currentPc = i.index
				var currentIdentCount = i.identCount
				var currentNumbersCount = i.numbersCount

				for range loopAmount - 1 {
					i.interpretStatements()
					i.index = currentPc
					i.identCount = currentIdentCount
					i.numbersCount = currentNumbersCount
				}
			}
			i.interpretStatements()
			var endToken = i.advance()
			if endToken != END {
				panic("Expected END in LOOP tail")
			}
			//fmt.Println("LOOP END\n")
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
				//fmt.Println(i.idents[currentIdentIndex], "=", i.vars[i.idents[otherIdentIndex]]+number)
				i.vars[i.idents[currentIdentIndex]] = i.vars[i.idents[otherIdentIndex]] + number
			} else if operationToken == MINUS {
				//fmt.Println(i.idents[currentIdentIndex], "=", i.vars[i.idents[otherIdentIndex]]-number)
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
		} else if tok == IDENT {
			i.identCount++
		} else if tok == NUM {
			i.numbersCount++
		}
	}
	i.index = j
}
