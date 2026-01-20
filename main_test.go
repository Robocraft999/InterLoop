package main

import (
	"testing"
)

func BenchmarkLarge(b *testing.B) {
	var res = "x := x + 2000000000 LOOP x DO y := y + 2 END\n y := x + 5 x := y - 10"
	var tokens, identsCount, idents, numbers = Lex(res)
	var interpreter = NewInterpreter(tokens, identsCount, idents, numbers)
	interpreter.Interpret()
}

func BenchmarkAddition(b *testing.B) {
	var res = "x := x + 500000 y := y + 500000 LOOP x DO\n  z := z + 1\nEND\nLOOP y DO\n  z := z + 1\nEND"
	var tokens, identsCount, idents, numbers = Lex(res)
	var interpreter = NewInterpreter(tokens, identsCount, idents, numbers)
	interpreter.Interpret()
}
