package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
)

const csvFilename = "./user_profiles.csv"

func main() {
	db, err := readCSV(csvFilename)
	if err != nil {
		log.Fatalf("failed reading CSV: %s", err)
	}

	env, err := cel.NewEnv(cel.Variable("user", cel.MapType(cel.StringType, cel.AnyType)))
	if err != nil {
		log.Fatalf("failed initializing CEL: %s", err)
	}

	if len(os.Args) != 2 {
		fmt.Printf(`Usage:

	%s <statement>

Statements must be valid CEL statements, which result in booleans, evaluated on
the contents of %s. Statements can utilize the following keys:

	user.id        (string)
	user.gender    (string)
	user.age       (integer)
	user.balance   (integer)

`, os.Args[0], csvFilename)

		os.Exit(2)
	}
	statement := os.Args[1]

	for _, line := range db {
		fmt.Println(line)
		out, err := eval(env, statement, line)
		if err != nil {
			log.Fatalf("failed evaluating example statement: %s", err)
		}
		fmt.Println(out)
	}
}

func readCSV(filepath string) ([]map[string]interface{}, error) {
	csvfile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer csvfile.Close()
	r := csv.NewReader(csvfile)
	data, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var ret []map[string]interface{}
	headers := data[0]
	for _, line := range data[1:] {
		entry := make(map[string]interface{}, 0)
		for i := 0; i < len(line); i++ {
			field := headers[i]
			cell := line[i]

			// First, we need to rename some fields; the CSV lists one thing but the example statements list a different one
			if field == "accountBalance" {
				field = "balance"
			}
			if field == "userId" {
				field = "id"
			}

			// Next, if we have age or balance fields, we need to convert to an int. Anything else is a string.
			if field == "age" || field == "balance" {
				convCell, err := strconv.Atoi(cell)
				if err != nil {
					return nil, err
				}
				entry[field] = convCell
			} else {
				entry[field] = cell
			}
		}
		ret = append(ret, entry)
	}
	return ret, nil
}

func eval(env *cel.Env, statement string, user map[string]interface{}) (bool, error) {
	ast, iss := env.Compile(statement)
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
