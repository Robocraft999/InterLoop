package main

import "fmt"

func main() {
	/*var reader = bufio.NewReader(os.Stdin)
	var bytes = make([]byte, 1024)
	var read, _ = reader.Read(bytes)
	bytes = bytes[:read]
	var res = string(bytes)*/
	var res = "x := x + 500000 y := y + 500000 LOOP x DO\n  z := z + 1\nEND\nLOOP y DO\n  z := z + 1\nEND"
	var tokens, identsCount, idents, numbers = Lex(res)
	fmt.Println(tokens)
	var interpreter = NewInterpreter(tokens, identsCount, idents, numbers)
	interpreter.Interpret()
}
