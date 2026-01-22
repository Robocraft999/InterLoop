package main

import (
	"bufio"
	"fmt"
	"os"
)

const SYNTAX_CHECK_ENABLED bool = false

func main() {
	var reader = bufio.NewReader(os.Stdin)
	var bytes = make([]byte, 1024)
	var read, _ = reader.Read(bytes)
	bytes = bytes[:read]
	var input = string(bytes)
	//var input = "x := x + 500000 y := y + 500000 LOOP x DO\n  z := z + 1\nEND\nLOOP y DO\n  z := z + 1\nEND"
	//var init = "n := zero + 29 "
	//var input = init + "a := zero + 0\nb := zero + 1\n\nLOOP n DO\n  sum := zero + 0\n  LOOP a DO\n    sum := sum + 1\n  END\n  LOOP b DO\n    sum := sum + 1\n  END\n  \n  a := zero + 0\n  LOOP b DO\n    a := a + 1\n  END\n  \n  b := zero + 0\n  LOOP sum DO\n    b := b + 1\n  END\nEND"
	var tokens, valIndices, vars = Lex(input)
	fmt.Println(tokens)
	var interpreter = NewInterpreter(tokens, valIndices, vars)
	interpreter.Interpret()
}
