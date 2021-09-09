package lisp

import (
	"fmt"
	"testing"
)

func TestBasicExpr(t *testing.T) {

	eval("100")
	eval("45.65")
	//eval("hello") failing
	eval("nil")
}

func TestList(t *testing.T) {
	eval("1 2 3 4")
	eval("(1 2 3 4)")
	//eval("george kate james joyce") failing
	//eval("(a (b c) (d (e f)))") failing
	eval("()")
}

func eval(str string) {
	lexer := NewScanner()
	lexer.Eval(str)
	parser := NewParser(lexer.Tokens)
	ast := parser.Parse()
	for _, expr := range ast {
		expr.PrintExpr()
	}
	fmt.Println()
}

//https://www.cs.unm.edu/~luger/ai-final2/LISP/CH%2011_S-expressions,%20The%20Syntax%20of%20Lisp.pdf
//symbolic expressions -> atom or list

//list => (atom | list)
//(1 2 3 4)
//(george kate james joyce)
//(a (b c) (d (e f)))
//() -> this is nil, its both atom and a list

//(+ 14 5)
//(+ 1 2 3 4)
//(– (+ 3 4) 7)
//(* (+ 2 5) (– 7 (/ 21 7)))
//(= (+ 2 3) 5)
//(> (* 5 6) (+ 4 5))
//(a b c)
