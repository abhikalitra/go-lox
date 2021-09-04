package main

import (
	"bufio"
	"fmt"
	"lisp"
	"os"
	"os/exec"
	"strings"
	"time"
)

func printRepl() {
	fmt.Print("\nlisp> ")
}

func recoverExp(text string) {
	if r := recover(); r != nil {
		fmt.Println("lisp> unknown command ", text)
	}
}

func handle(text string) {
	// We might have a panic here we so need DEFER + RECOVER
	defer recoverExp(text)
	// \n Will be ignored
	t := strings.TrimSuffix(text, "\n")
	if t != "" {
		lisp.Eval_lisp(t)
	}
}

func get(r *bufio.Reader) string {
	t, _ := r.ReadString('\n')
	return strings.TrimSpace(t)
}

func shouldContinue(text string) bool {
	if strings.EqualFold("exit", text) {
		return false
	}
	return true
}

func help() {
	fmt.Println("lisp> Welcome to LISP! ")
	fmt.Println("lisp> time   - Prints current date / time ")
	fmt.Println("lisp> ")
}

func cls() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func now() {
	fmt.Println("lisp> ", time.Now().Format(time.RFC850))
}

func main() {
	commands := map[string]interface{}{
		"help": help,
		"cls":  cls,
		"time": now,
	}
	reader := bufio.NewReader(os.Stdin)
	help()
	printRepl()
	text := get(reader)
	for ; shouldContinue(text); text = get(reader) {
		if value, exists := commands[text]; exists {
			value.(func())()
		} else {
			handle(text)
		}
		printRepl()
	}
	fmt.Println("Bye!")

}
