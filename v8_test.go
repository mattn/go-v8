package v8

import (
	"fmt"
	"sync"
	"testing"
)

var (
	JS_SIMPLE string = `a = 1; a++;`
	JS_FIB    string = `function fib(n) {
			var f = 0,
				n1 = 1,
				n2 = 0;
			if (n <= 1) {
				return n
			}

			for (var i = 1; i < n; i++) {
				f = n1 + n2;
				n2 = n1;
				n1 = f;

			}
			return f;
		}
		fib(%d);`
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
		num := args[0].(float64)
		argVal := int(num)
		if argVal != 10 {
			t.Fatal("Unexpected value for arg 0, expected 10, received:", argVal)
		}

		// Second argument
		argVal2 := args[1].(string)
		if argVal2 != `Test string` {
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

func TestEvalScriptMap(t *testing.T) {
	ctx := NewContext()

	res, err := ctx.Eval(`var a = {"foo": "bar"}; a`)
	if err != nil {
		t.Fatal("Unexpected error on eval,", err)
	}
	if res == nil {
		t.Fatal("Expected result from eval, received nil")
	}

	v, ok := res.(map[string]interface{})
	if !ok {
		t.Fatal("Expected map, but not:", res)
	}

	r, ok := v["foo"].(string)
	if !ok {
		t.Fatal("Expected string, but %q", res)
	}

	if r != "bar" {
		t.Fatalf("Expected %q, but %q", "bar", r)
	}
}

func TestEvalScriptArray(t *testing.T) {
	ctx := NewContext()

	res, err := ctx.Eval(`var a = ["foo", true]; a`)
	if err != nil {
		t.Fatal("Unexpected error on eval,", err)
	}
	if res == nil {
		t.Fatal("Expected result from eval, received nil")
	}

	v, ok := res.([]interface{})
	if !ok {
		t.Fatal("Expected array, but not:", res)
	}

	if len(v) != 2 {
		t.Fatalf("Expected two elements, but %d elements", len(v))
	}

	e1, ok := v[0].(string)
	if !ok {
		t.Fatal("Expected a[0] is string, but not")
	}
	e2, ok := v[1].(bool)
	if !ok {
		t.Fatal("Expected a[1] is boolean, but not")
	}

	if e1 != "foo" || e2 != true {
		t.Fatal(`Expected a[0]=="foo", a[1]==true but not`)
	}
}

func TestAddFuncReturnObject(t *testing.T) {
	ctx := NewContext()
	err := ctx.AddFunc("testFunc", func(args ...interface{}) interface{} {
		return map[string]interface{}{
			"arg0": args[0].(float64),
			"arg1": args[1].(string),
			"arg2": args[2].(bool),
		}
	})
	if err != nil {
		t.Fatal("Expected to be able to add function, received error ", err)
	}

	res, err := ctx.Eval(`testFunc(10, "something", true).arg0`)
	if err != nil {
		t.Fatal("Unexpected error on testFunc eval ", err)
	}
	if res == nil {
		t.Fatal("Expected result from testFunc eval, received nil")
	}
	if int(res.(float64)) != 10 {
		t.Fatal("Expected result to be 10, got", res)
	}

	res, err = ctx.Eval(`testFunc(10, "something", false)`)
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
	arg2 := resMap["arg2"].(bool)
	if arg2 != false {
		t.Fatal("Expected arg2 value to be false got ", arg1)
	}
}

func v8EvalRoutine(i int, wg *sync.WaitGroup, t *testing.T) {
	defer wg.Done()

	ctx := NewContext()
	res, err := ctx.Eval(fmt.Sprintf(JS_FIB, i))
	if err != nil {
		t.Fatal("Failed to evaluate test fib function for index,", i, "error:", err)
	}
	if res == nil {
		t.Fatal("Unexpected nil for result of test fib function for index", i)
	}
	r := uint64(res.(float64))
	if !((i == 80 && r == 23416728348467684) ||
		(i == 50 && r == 12586269025) ||
		(i == 20 && r == 6765) ||
		(i == 60 && r == 1548008755920)) {
		t.Fatal("Failed to calculate correct fib for index", i, "received value,", r)
	}
}

func TestMultiRoutines(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go v8EvalRoutine(80, &wg, t)

	wg.Add(1)
	go v8EvalRoutine(50, &wg, t)

	wg.Add(1)
	go v8EvalRoutine(20, &wg, t)

	wg.Add(1)
	go v8EvalRoutine(60, &wg, t)

	wg.Wait()
}

func BenchmarkCreateContext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewContext()
	}
}

func BenchmarkEvalSimple(b *testing.B) {
	b.StopTimer()
	ctx := NewContext()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		ctx.Eval(JS_SIMPLE)
	}
}
