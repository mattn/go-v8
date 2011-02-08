package main

import "v8"
import "reflect"

func main() {
	v8ctx := v8.NewContext()
	ret, err := v8ctx.Eval(`
var a = 1;
a += 2;
a;
`)
	if err != nil {
		println(err.String())
	} else {
		println(ret.(float64))
	}

	ret, err = v8ctx.Eval(`
a+'b'
`)
	println(reflect.NewValue(ret).Type().Name())
	if err != nil {
		println(err.String())
	} else {
		println(ret.(string))
	}
}
