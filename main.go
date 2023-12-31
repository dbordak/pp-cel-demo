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

var errNonBool error = errors.New("statements must evaluate to a boolean")

func main() {
	db, err := readCSV(csvFilename)
	if err != nil {
		log.Fatalf("failed reading CSV: %s", err)
	}

	if len(os.Args) != 2 {
		fmt.Printf(`Usage:

	%[1]s <statement>

Statements must be valid CEL statements, which must result in booleans, evaluated
linewise on the contents of %[2]s.

Example Usage:
		%[1]s 'user.gender == "male"'
		%[1]s 'user.balance >= 500.0 && user.gender == "female" && user.age <= 30'

Statements can utilize the following keys:

	user.id        (string)
	user.gender    (string)
	user.age       (integer)
	user.balance   (integer)

`, os.Args[0], csvFilename)

		os.Exit(2)
	}
	statement := os.Args[1]

	prg, err := initCEL(statement)
	if err != nil {
		log.Fatalf("failed initializing CEL: %s", err)
	}
	out, err := average(prg, db)
	if err != nil {
		log.Fatalf("failed evaluating average: %s", err)
	}
	fmt.Printf("Average Balance: %f\n", out)
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

func initCEL(statement string) (cel.Program, error) {
	env, err := cel.NewEnv(cel.Variable("user", cel.MapType(cel.StringType, cel.AnyType)))
	if err != nil {
		return nil, err
	}
	ast, iss := env.Compile(statement)
	if iss.Err() != nil {
		return nil, err
	}
	return env.Program(ast)
}

func evalToBool(prg cel.Program, user map[string]interface{}) (bool, error) {
	out, _, err := prg.Eval(map[string]interface{}{
		"user": user,
	})
	if err != nil {
		return false, err
	}

	// *technically* the requirements only state that the filter is to include
	// anything that evaluates to true -- meaning you could construct a
	// statement that evaluates to a boolean for some lines and an object for
	// others, where the latter should be treated as false. However, I'm going
	// to assume this was not a gotcha and that erroring out on seemingly
	// invalid results is the more correct thing to do.
	if !types.IsBool(out) {
		return false, errNonBool
	}
	refType := reflect.TypeOf(true)
	val, err := out.ConvertToNative(refType)
	return val.(bool), err
}

func average(prg cel.Program, db []map[string]interface{}) (float64, error) {
	sum := 0
	count := 0
	for _, line := range db {
		out, err := evalToBool(prg, line)
		if err != nil {
			return 0, err
		}
		if out {
			count++
			sum += line["balance"].(int)
		}
	}
	return float64(sum) / float64(count), nil
}
