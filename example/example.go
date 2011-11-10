package main

import (
	"fmt"
	"reflect"
	"github.com/mattn/go-v8/v8"
)

func main() {
	v8ctx := v8.NewContext()
	ret, err := v8ctx.Eval(`
var a = 1;
a += 2;
a;
`)
	if err != nil {
		fmt.Println(err)
	} else {
		println(ret.(float64))
	}

	ret, err = v8ctx.Eval(`
a+'b'
`)
	println(reflect.ValueOf(ret).Type().Name())
	if err != nil {
		fmt.Println(err)
	} else {
		println(ret.(string))
	}
}
