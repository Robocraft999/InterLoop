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
	for c := i.current(); c != EOF && c != END; {
		i.interpretStatement()
	}
}

func (i *Interpreter) interpretStatement() {
	var current = i.current()
	i.index++
	if current == LOOP {
		if SYNTAX_CHECK_ENABLED && i.current() != IDENT {
			panic("Expected IDENT in LOOP head")
		}
		i.index++

		var loopAmountIndex = i.identCount
		i.identCount++

		//fmt.Println("LOOP START WITH ", i.idents[loopAmountIndex], " = ", i.vars[i.idents[loopAmountIndex]])
		if SYNTAX_CHECK_ENABLED && i.current() != DO {
			panic("Expected DO in LOOP head")
		}
		i.index++

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

		if SYNTAX_CHECK_ENABLED && i.current() != END {
			panic("Expected END in LOOP tail")
		}
		i.index++
		//fmt.Println("LOOP END\n")
		return
	}
	if current == IDENT {
		var currentIdentIndex = i.identCount
		i.identCount++

		if SYNTAX_CHECK_ENABLED && i.current() != ASSIGN {
			panic("Expected ASSIGN in statement")
		}
		i.index++

		if SYNTAX_CHECK_ENABLED && i.current() != IDENT {
			panic("Expected IDENT in statement")
		}
		i.index++
		var otherIdentIndex = i.identCount
		i.identCount++

		var operationToken = i.current()
		i.index++

		if SYNTAX_CHECK_ENABLED && i.current() != NUM {
			panic("Expected NUM in statement")
		}
		i.index++
		var numberIndex = i.numbersCount
		i.numbersCount++
		var number = i.numbers[numberIndex]

		if operationToken == PLUS {
			//fmt.Println(i.idents[currentIdentIndex], "=", i.vars[i.idents[otherIdentIndex]]+number)
			i.vars[i.idents[currentIdentIndex]] = i.vars[i.idents[otherIdentIndex]] + number
		}
		if operationToken == MINUS {
			//fmt.Println(i.idents[currentIdentIndex], "=", i.vars[i.idents[otherIdentIndex]]-number)
			i.vars[i.idents[currentIdentIndex]] = i.vars[i.idents[otherIdentIndex]] - number
		}
		if SYNTAX_CHECK_ENABLED && operationToken != PLUS && operationToken != MINUS {
			panic("Expected PLUS or MINUS in statement")
		}
		return
	}
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
		}
		if tok == END {
			count--
		}
		if tok == IDENT {
			i.identCount++
		}
		if tok == NUM {
			i.numbersCount++
		}
	}
	i.index = j
}
