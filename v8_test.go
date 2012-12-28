package v8

import (
	"testing"
)

func TestEvalScript(t *testing.T) {
	ctx := NewContext()

	res, err := ctx.Eval(`var a = 10; a`)
	if err != nil {
		t.Fatal("Unexpected error on eval,", err)
	}
	if res == nil {
		t.Fatal("Expected result from eval, received nil")
	}

	switch res.(type) {
	case float64:
	default:
		t.Fatal("Expected float64 type")
	}
	if 10 != int(res.(float64)) {
		t.Fatal("Expected result to be 10, received:", res)
	}
}

func TestAddFunc(t *testing.T) {
	ctx := NewContext()

	err := ctx.AddFunc("_gov8_testFunc", func(args ...interface{}) interface{} {
		if len(args) != 2 {
			t.Fatal("Unexpected number of _gov8_testFunc's arguments.", len(args))
		}
		// First argument
		arg := args[0]
		switch arg.(type) {
		case float64:
		default:
			t.Fatal("Unexpected arg 0 type, expecting float64")
		}
		argVal := int(arg.(float64))
		if argVal != 10 {
			t.Fatal("Unexpected value for arg 0, expected 10, received:", argVal)
		}

		// Second argument
		arg = args[1]
		switch arg.(type) {
		case string:
		default:
			t.Fatal("Unexpected arg 1 type, expected string")
		}
		argVal2 := arg.(string)
		if argVal2 != "Test string" {
			t.Fatal("Unexpected value for arg 1, expected Test string, received:", argVal2)
		}

		return "testFunc return value"
	})
	if err != nil {
		t.Fatal("Expected to be able to add function, received error ", err)
	}

	res, err := ctx.Eval(`_gov8_testFunc(10, "Test string");`)
	if err != nil {
		t.Fatal("Unexpected error on testFunc eval,", err)
	}
	if res == nil {
		t.Fatal("Expected result from testFunc eval, received nil")
	}
	if res.(string) != "testFunc return value" {
		t.Fatal("Unexpected result from eval,", res)
	}
}

func TestAddFuncReturnObject(t *testing.T) {
	ctx := NewContext()
	err := ctx.AddFunc("testFunc", func(args ...interface{}) interface{} {
		return map[string]interface{}{
			"arg0": int(args[0].(float64)),
			"arg1": args[1].(string),
		}
	})
	if err != nil {
		t.Fatal("Expected to be able to add function, received error ", err)
	}

	res, err := ctx.Eval(`testFunc(10, "something").arg0`)
	if err != nil {
		t.Fatal("Unexpected error on testFunc eval ", err)
	}
	if res == nil {
		t.Fatal("Expected result from testFunc eval, received nil")
	}
	if int(res.(float64)) != 10 {
		t.Fatal("Expected result to be 10, got", res)
	}

	res, err = ctx.Eval(`testFunc(10, "something")`)
	if err != nil {
		t.Fatal("Unexpected error on testFunc eval ", err)
	}
	if res == nil {
		t.Fatal("Expected result from testFunc eval, received nil")
	}
	resMap := res.(map[string]interface{})
	arg0 := int(resMap["arg0"].(float64))
	if arg0 != 10 {
		t.Fatal("Expected arg0 value to be 10 got ", arg0)
	}
	arg1 := resMap["arg1"].(string)
	if arg1 != "something" {
		t.Fatal("Expected arg1 value to be something got ", arg1)
	}
}
