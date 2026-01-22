package main

import "fmt"

type Interpreter struct {
	tokens     []Token
	valIndices []int
	valIndex   int
	vars       []int
	index      int
}

func NewInterpreter(tokens []Token, valIndices []int, vars []int) *Interpreter {
	return &Interpreter{
		tokens:     tokens,
		valIndices: valIndices,
		valIndex:   0,
		vars:       vars,
		index:      0,
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
	var current = i.current()
	i.index++
	if current == LOOP {
		if SYNTAX_CHECK_ENABLED && i.current() != IDENT {
			panic("Expected IDENT in LOOP head")
		}
		i.index++

		var loopAmountIndex = i.valIndex
		i.valIndex++

		//fmt.Println("LOOP START WITH ", i.valIndices[loopAmountIndex], " = ", i.vars[i.valIndices[loopAmountIndex]])
		if SYNTAX_CHECK_ENABLED && i.current() != DO {
			panic("Expected DO in LOOP head")
		}
		i.index++

		var loopAmount = i.vars[i.valIndices[loopAmountIndex]]

		if loopAmount == 0 {
			i.jumpToEnd()
			//fmt.Println("LOOP SKIPPED\n")
			return
		}

		if loopAmount > 1 {
			var currentPc = i.index
			var currentValIndex = i.valIndex

			for range loopAmount - 1 {
				i.interpretStatements()
				i.index = currentPc
				i.valIndex = currentValIndex
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
		var currentIdentIndex = i.valIndex
		i.valIndex++

		if SYNTAX_CHECK_ENABLED && i.current() != ASSIGN {
			panic("Expected ASSIGN in statement")
		}
		i.index++

		if SYNTAX_CHECK_ENABLED && i.current() != IDENT {
			panic("Expected IDENT in statement")
		}
		i.index++
		var otherIdentIndex = i.valIndex
		i.valIndex++

		var operationToken = i.current()
		i.index++

		if SYNTAX_CHECK_ENABLED && i.current() != NUM {
			panic("Expected NUM in statement")
		}
		i.index++
		var numberIndex = i.valIndex
		i.valIndex++
		var number = i.vars[i.valIndices[numberIndex]]

		if operationToken == PLUS {
			//fmt.Println(i.valIndices[currentIdentIndex], "=", i.vars[i.valIndices[otherIdentIndex]]+number)
			i.vars[i.valIndices[currentIdentIndex]] = i.vars[i.valIndices[otherIdentIndex]] + number
		}
		if operationToken == MINUS {
			//fmt.Println(i.valIndices[currentIdentIndex], "=", i.vars[i.valIndices[otherIdentIndex]]-number)
			i.vars[i.valIndices[currentIdentIndex]] = i.vars[i.valIndices[otherIdentIndex]] - number
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
		if tok == IDENT || tok == NUM {
			i.valIndex++
		}
	}
	i.index = j
}
