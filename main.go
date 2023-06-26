package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
)

func main() {
	env, err := cel.NewEnv(cel.Variable("user", cel.MapType(cel.StringType, cel.AnyType)))
	if err != nil {
		log.Fatalf("failed initializing CEL: %s", err)
	}

	example := `user.balance >= 500.0 && user.gender == "female" && user.age <= 30 && user.id == "yf2"`
	ast, iss := env.Compile(example)
	// Check iss for compilation errors.
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	prg, err := env.Program(ast)
	out, _, err := prg.Eval(map[string]interface{}{
		"user": map[string]interface{}{
			"userId":         "yn1",
			"gender":         "male",
			"age":            10,
			"accountBalance": 500,
		},
	})
	fmt.Println(out)

}
