package main

import (
	"fmt"
	"github.com/mattn/go-v8"
)

func main() {
	v8ctx := v8.NewContext()

	// setup console.log()
	v8ctx.AddFunc("_console_log", func(args ...interface{}) (interface{}, error) {
		fmt.Printf("Go console log: ")
		for i := 0; i < len(args); i++ {
			fmt.Printf("%v ", args[i])
		}
		fmt.Println()
		return "", nil
	})

	v8ctx.Eval(`
	this.console = { "log": function(args) { _console_log.apply(null, arguments) }}`)
	ret := v8ctx.MustEval(`
	var a = 1;
	var b = 'B'
	a += 2;
	a;
	`)
	fmt.Println("Eval result:", int(ret.(float64))) // 3

	v8ctx.Eval(`console.log(a + '年' + b + '組 金八先生！', 'something else')`) // 3b
	v8ctx.Eval(`console.log("Hello World, こんにちわ世界")`)                    // john
	v8ctx.Eval(`console.log({"hoge": "fuga"})`)                          // john
	v8ctx.AddFunc("func_call", func(args ...interface{}) (interface{}, error) {
		f := func(args ...interface{}) (interface{}, error) {
			return "V8", nil
		}
		ret, _ := args[0].(v8.Function).Call("Go", 2, 1, f)
		return ret, nil
	})

	fmt.Println(v8ctx.MustEval(`
		func_call(function() {
			return "Hello " + arguments[0] + (arguments[1] - arguments[2])
				+ ", Hello " + arguments[3]();
		})
		`).(string))
}
