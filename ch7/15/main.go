// Exercise 7.15: Write a program that reads a single expression from the standard input,
// prompts the user to provide values for any variables, then evaluates the expression in the
// resulting environment. Handle all errors gracefully.

package main

import (
	"fmt"
	"gopl-exercises/ch7/eval"
	"log"
	"os"
	"regexp"
	"strconv"
)

const (
	usage = `./main "<expression>" "argument1" "argument2" ...
	argument format: <single letter name>=<single numeric value>`
)

var r = regexp.MustCompile(`^([a-zA-Z])[=]([\d]+)$`)

// !+parseExprAndCheck
func parseExprAndCheck(s string, env *eval.Env) (eval.Expr, error) {
	if s == "" {
		return nil, fmt.Errorf("empty expression")
	}
	expr, err := eval.Parse(s)
	if err != nil {
		return nil, err
	}
	vars := make(map[eval.Var]bool)
	if err := expr.Check(vars); err != nil {
		return nil, err
	}

	for v := range vars {
		if _, ok := (*env)[v]; !ok {
			return nil, fmt.Errorf("undefined variable: %s", v)
		}
	}

	return expr, nil
}

// !+parseEnv
func parseEnv(input []string) (eval.Env, error) {
	env := make(eval.Env)

	for _, s := range input {
		k, v, err := parseVariable(s)
		if err != nil {
			fmt.Printf("invalid variable %s\n", s)
			continue
		}
		env[eval.Var(k)] = v
	}

	return env, nil
}

func parseVariable(input string) (variable string, value float64, err error) {
	allSubmatches := r.FindStringSubmatch(input)
	if len(allSubmatches) == 0 {
		err = fmt.Errorf("no variable matched in %s", input)
		return
	}

	variable = allSubmatches[1]
	value, err = strconv.ParseFloat(allSubmatches[2], 64)
	return
}

// !+evaluate
func evaluate(input []string) (result float64, err error) {
	env, err := parseEnv(input[2:])
	if err != nil {
		err = fmt.Errorf("bad variables provided: %s", err)
		return
	}

	e := input[1]
	expr, err := parseExprAndCheck(e, &env)
	if err != nil {
		err = fmt.Errorf("bad expression: %s", err)
		return
	}

	result = expr.Eval(env)
	return
}

func main() {
	if len(os.Args) < 2 {
		log.Println("No expression provided. Usage " + usage)
		os.Exit(0)
	}

	result, err := evaluate(os.Args)
	if err != nil {
		fmt.Println(err.Error() + "\n" + "Usage: " + usage)
		os.Exit(1)
	}

	fmt.Printf("Result: %f\n", result)
}
