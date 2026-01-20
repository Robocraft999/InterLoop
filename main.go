package main

import "fmt"

func main() {
	/*var reader = bufio.NewReader(os.Stdin)
	var bytes = make([]byte, 1024)
	var read, _ = reader.Read(bytes)
	bytes = bytes[:read]
	var res = string(bytes)*/
	var res = "x := x + 20000 LOOP x DO y := y + 2 END\n y := x + 5 	x := y - 10"
	var tokens, identsCount, idents, numbers = Lex(res)
	fmt.Println(tokens)
	var interpreter = NewInterpreter(tokens, identsCount, idents, numbers)
	interpreter.Interpret()
}
