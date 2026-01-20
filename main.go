package main

import "fmt"

func main() {
	/*var reader = bufio.NewReader(os.Stdin)
	var bytes = make([]byte, 1024)
	var read, _ = reader.Read(bytes)
	bytes = bytes[:read]
	var res = string(bytes)*/
	var res = "x := x + 2 LOOP x DO y := y + 2 END"
	var tokens, idents, numbers = Lex(res)
	fmt.Println(tokens)
	var interpreter = NewInterpreter(tokens, idents, numbers)
	interpreter.Interpret()
}
