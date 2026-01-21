package main

import (
	"testing"
)

func run(input string) {
	var tokens, identsCount, idents, numbers = Lex(input)
	var interpreter = NewInterpreter(tokens, identsCount, idents, numbers)
	interpreter.Interpret()
}

func BenchmarkLarge(b *testing.B) {
	var input = "x := x + 2000000000 LOOP x DO y := y + 2 END\n y := x + 5 x := y - 10"
	run(input)
}

func BenchmarkAddition(b *testing.B) {
	var init = "x := x + 500000 y := y + 500000 "
	var input = init + "LOOP x DO\n  z := z + 1\nEND\nLOOP y DO\n  z := z + 1\nEND"
	run(input)
}

func BenchmarkFactorial(b *testing.B) {
	var init = "n := n + 10 "
	var input = init + "factorial := zero + 1\ni := zero + 0\n\nLOOP n DO\n  i := i + 1\n  \n  oldFactorial := zero + 0\n  LOOP factorial DO\n    oldFactorial := oldFactorial + 1\n  END\n  \n  factorial := zero + 0\n  LOOP i DO\n    LOOP oldFactorial DO\n      factorial := factorial + 1\n    END\n  END\nEND"
	run(input)
}

func BenchmarkFibonacci(b *testing.B) {
	var init = "n := n + 29 "
	var input = init + "a := zero + 0\nb := zero + 1\n\nLOOP n DO\n  sum := zero + 0\n  LOOP a DO\n    sum := sum + 1\n  END\n  LOOP b DO\n    sum := sum + 1\n  END\n  \n  a := zero + 0\n  LOOP b DO\n    a := a + 1\n  END\n  \n  b := zero + 0\n  LOOP sum DO\n    b := b + 1\n  END\nEND"
	run(input)
}
