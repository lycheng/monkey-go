package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/lycheng/monkey-go/evaluator"
	"github.com/lycheng/monkey-go/lexer"
	"github.com/lycheng/monkey-go/object"
	"github.com/lycheng/monkey-go/parser"
)

const (
	prompt = ">> "
)

// Start to read input from in and print the parsed result to out
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		switch rv := evaluated.(type) {
		case *object.ReturnValue:
		case *object.Error:
			io.WriteString(out, rv.Inspect())
			io.WriteString(out, "\n")
		case *object.Null:
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
