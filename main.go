package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
)

func main() {
	env, err := cel.NewEnv(cel.Variable("name", cel.StringType))
	// Check err for environment setup errors.
	if err != nil {
		log.Fatalf("failed initializing CEL: %s", err)
	}
	ast, iss := env.Compile(`"Hello world! I'm " + name + "."`)
	// Check iss for compilation errors.
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	prg, err := env.Program(ast)
	out, _, err := prg.Eval(map[string]interface{}{
		"name": "CEL",
	})
	fmt.Println(out)

}
