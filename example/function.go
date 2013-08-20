package main

import (
	"fmt"
	"github.com/mattn/go-v8"
)

func main() {
	ctx := v8.NewContext()

	ctx.MustEval(`
	this.console = { "log": function(args) { _console_log.apply(null, arguments) }}`)
	ctx.AddFunc("_console_log", func(args ...interface{}) (interface{}, error) {
		for i := 0; i < len(args); i++ {
			fmt.Printf("%v ", args[i])
		}
		fmt.Println()
		return "", nil
	})

	ctx.MustEval(`
var b = function() { console.log('hi', arguments[0]); }

function a() {
    b(arguments[0]);
}
	`)

	ctx.MustEval("a").(v8.Function).Call("Cameron Currie")
}
