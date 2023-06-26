package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
)

func main() {
	env, err := cel.NewEnv(cel.Variable("user", cel.MapType(cel.StringType, cel.AnyType)))
	if err != nil {
		log.Fatalf("failed initializing CEL: %s", err)
	}
	user := map[string]interface{}{
		"userId":         "yn1",
		"gender":         "male",
		"age":            10,
		"accountBalance": 500,
	}

	out, err := example(env, user)
	if err != nil {
		log.Fatalf("failed evaluating example statement: %s", err)
	}
	fmt.Println(out)
}

func example(env *cel.Env, user map[string]interface{}) (bool, error) {
	example := `user.balance >= 500.0 && user.gender == "female" && user.age <= 30 && user.id == "yf2"`
	ast, iss := env.Compile(example)

	if iss.Err() != nil {
		return false, iss.Err()
	}

	prg, err := env.Program(ast)
	if err != nil {
		return false, err
	}
	out, _, err := prg.Eval(map[string]interface{}{
		"user": user,
	})
	if err != nil {
		return false, err
	}

	if !types.IsBool(out) {
		return false, errors.New("statements must evaluate to a boolean")
	}
	refType := reflect.TypeOf(true)
	val, err := out.ConvertToNative(refType)
	return val.(bool), err
}
