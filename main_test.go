package main

import (
	"math"
	"testing"
)

var db []map[string]interface{}

func TestMain(t *testing.T) {
	db, _ = readCSV("./user_profiles.csv")
}

func TestEvalNonBool(t *testing.T) {
	// statement 'user' means the result of each expression will be an object
	// (i.e. the entire row from the CSV). Since this is a filter operation, we
	// can't work with that, so it should error out.
	failCase, _ := initCEL("user")
	for _, line := range db {
		_, err := evalToBool(failCase, line)
		if err != errNonBool {
			t.Fatalf("Should have failed with non bool error, instead got %v", err)
		}
	}
}

func TestEvalTrue(t *testing.T) {
	alwaysTrue, _ := initCEL(`true`)
	for _, line := range db {
		out, err := evalToBool(alwaysTrue, line)
		if err != nil {
			t.Fatalf("Err when none expected %v", err)
		}
		if !out {
			t.Fatal("False output on always-true test")
		}
	}
}

func TestEvalExample(t *testing.T) {
	// examples from the project description
	case1, _ := initCEL(`user.balance >= 500.0 && user.gender == "female" && user.age <= 30 && user.id == "yf2"`)
	expected1 := []bool{false, false, false, false, false, false, false, true, false, false, false, false}
	for i, line := range db {
		out, err := evalToBool(case1, line)
		if err != nil {
			t.Fatalf("Err when none expected %v", err)
		}
		if out != expected1[i] {
			t.Fatalf("Mismatched output on case 1 row %d: Expected %t, Actual %t", i, expected1[i], out)
		}
	}

	case2, _ := initCEL(`user.gender == "male"`)
	expected2 := []bool{true, true, true, true, true, true, false, false, false, false, false, false}
	for i, line := range db {
		out, err := evalToBool(case2, line)
		if err != nil {
			t.Fatalf("Err when none expected %v", err)
		}
		if out != expected2[i] {
			t.Fatalf("Mismatched output on case 2 row %d: Expected %t, Actual %t", i, expected2[i], out)
		}
	}

	case3, _ := initCEL(`user.balance >= 500.0 && user.gender == "female" && user.age <= 30`)
	expected3 := []bool{false, false, false, false, false, false, true, true, false, false, false, false}
	for i, line := range db {
		out, err := evalToBool(case3, line)
		if err != nil {
			t.Fatalf("Err when none expected %v", err)
		}
		if out != expected3[i] {
			t.Fatalf("Mismatched output on case 3 row %d: Expected %t, Actual %t", i, expected3[i], out)
		}
	}
}

func TestAverageExample(t *testing.T) {
	// examples from the project description
	case1, _ := initCEL(`user.balance >= 500.0 && user.gender == "female" && user.age <= 30 && user.id == "yf2"`)
	expected1 := 1000.0
	out, err := average(case1, db)
	if err != nil {
		t.Fatalf("Err when none expected %v", err)
	}
	if math.Abs(out-expected1) > 0.1 {
		t.Fatalf("Mismatched output on case 1: Expected %f, Actual %f", expected1, out)
	}

	case2, _ := initCEL(`user.gender == "male"`)
	expected2 := 383.333333
	out, err = average(case2, db)
	if err != nil {
		t.Fatalf("Err when none expected %v", err)
	}
	if math.Abs(out-expected2) > 0.1 {
		t.Fatalf("Mismatched output on case 2: Expected %f, Actual %f", expected2, out)
	}

	case3, _ := initCEL(`user.balance >= 500.0 && user.gender == "female" && user.age <= 30`)
	expected3 := 950.0
	out, err = average(case3, db)
	if err != nil {
		t.Fatalf("Err when none expected %v", err)
	}
	if math.Abs(out-expected3) > 0.1 {
		t.Fatalf("Mismatched output on case 3: Expected %f, Actual %f", expected3, out)
	}
}
