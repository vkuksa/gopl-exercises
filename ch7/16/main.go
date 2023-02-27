// Exercise 7.16: Write a web-based calculator program.

package main

import (
	"fmt"
	"gopl-exercises/ch7/eval"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

const (
	usage = `./main "<expression>" "argument1" "argument2" ...
	argument format: <single letter name>=<single numeric value>`
)

var r = regexp.MustCompile(`^[\d]+$`)

// !+parseAndCheckExpr
func parseAndCheckExpr(s string, env *eval.Env) (eval.Expr, error) {
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

func parseVariable(input string) (value float64, err error) {
	allSubmatches := r.FindStringSubmatch(input)
	if len(allSubmatches) == 0 {
		err = fmt.Errorf("no variable matched in %s", input)
		return
	}

	value, err = strconv.ParseFloat(allSubmatches[0], 64)
	return
}

func extractRawExprAndEnv(values url.Values) (expr string, env eval.Env, err error) {
	if !values.Has("expr") {
		err = fmt.Errorf("no expr provided")
		return
	}
	expr = values.Get("expr")

	env = make(eval.Env)
	for k, v := range values {
		if k != "expr" {
			val, err := parseVariable(v[0])
			if err != nil {
				fmt.Printf("invalid variable %s\n", v[0])
				continue
			}
			env[eval.Var(k)] = val
		}
	}

	return
}

func calculator(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	e, env, err := extractRawExprAndEnv(r.Form)
	if err != nil {
		http.Error(w, "bad query provided: "+err.Error(), http.StatusBadRequest)
		return
	}

	expr, err := parseAndCheckExpr(e, &env)
	if err != nil {
		http.Error(w, "bad expr: "+err.Error(), http.StatusBadRequest)
		return
	}

	result := expr.Eval(env)
	fmt.Fprintln(w, result)
	return
}

func main() {
	http.HandleFunc("/", calculator)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
