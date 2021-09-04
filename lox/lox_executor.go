package lox

import "fmt"

var lexer = NewScanner()

func Eval_lisp(expr string) {
	lexer.Eval(expr)
	fmt.Printf("%v+", lexer.Tokens)
}

//https://gigamonkeys.com/book/lather-rinse-repeat-a-tour-of-the-repl.html
//1. Number
//10

//2. Expression
//(+ 2 3)

//3. String
//"hello, world"

//4. Call function
//(format t "hello, world")

//5. define function
//(defun hello-world () (format t "hello, world"))

//6. call user defined function
//(hello-world)

//7. load a file
//(load "hello.lisp")

//8. call function from loaded file
// (hello-world)

//9. compile and load file
// (load (compile-file "hello.lisp"))
