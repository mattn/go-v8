package v8

/*
#cgo LDFLAGS: -L. -lv8wrap -lstdc++

#include <stdlib.h>
#include "v8wrap.h"

extern char* _go_v8_callback(unsigned int, char*, v8data*, int);

static char*
_c_v8_callback(unsigned int id, char* n, v8data* d, int a) {
  return _go_v8_callback(id, n, d, a);
}

static void
v8_callback_init() {
  v8_init((void*) _c_v8_callback);
}
*/
import "C"
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"text/template"
	"unsafe"

	"github.com/cwc/jsregexp"
)

var contexts = make(map[uint32]*V8Context)

var tmpl = template.Must(template.New("go-v8").Parse(`
function {{.name}}() {
  return _go_call({{.id}}, "{{.name}}", arguments);
}`))

type V8Function struct {
	ctx  *V8Context
	repr string
}

type V8Object struct {
	Name string
}

func (f V8Function) Call(args ...interface{}) (interface{}, error) {
	var arguments bytes.Buffer
	for i, arg := range args {
		fn, ok := arg.(func(...interface{}) interface{})
		if ok {
			// arg is a Go func
			name := fmt.Sprintf("anonymous%v", fn)
			f.ctx.funcs[name] = fn
			buf := bytes.NewBufferString("")
			tmpl.Execute(buf, map[string]interface{}{
				"id":   f.ctx.id,
				"name": name,
			})
			arguments.WriteString("(" + buf.String() + ")")
		} else {
			obj, ok := arg.(V8Object)
			if ok {
				// arg is a JavaScript object
				arguments.WriteString(obj.Name)
			} else {
				// arg is a Go object; marshal to JSON
				b, err := json.Marshal(arg)
				if err != nil {
					return nil, err
				}
				arguments.WriteString(string(b))
			}
		}
		if i != len(args)-1 {
			arguments.WriteString(",")
		}
	}

	return f.ctx.Eval("(" + f.repr + ")(" + arguments.String() + ")")
}

func (f V8Function) String() string {
	return f.repr
}

//export _go_v8_callback
func _go_v8_callback(contextId uint32, functionName *C.char, v8Objects *C.v8data, count C.int) *C.char {
	ctx := contexts[contextId]
	fn := ctx.funcs[C.GoString(functionName)]

	if fn != nil {
		var argv []interface{}

		// Parse objects
		i := C.int(0)
		for ; i < count; i++ {
			obj := C.v8_get_array_item(v8Objects, i)

			switch obj.obj_type {
			case C.v8regexp:
				argv = append(argv, regexp.MustCompile(jsregexp.Translate(C.GoString(obj.repr))))
				break
			case C.v8function:
				argv = append(argv, V8Function{ctx, C.GoString(obj.repr)})
				break
			default:
				// Should be a JSON string, so pass it as-is
				argv = append(argv, C.GoString(obj.repr))
			}
		}

		// Call function
		ret := fn(argv...)
		if ret != nil {
			b, _ := json.Marshal(ret)
			return C.CString(string(b))
		}
		return nil
	}
	return C.CString("undefined")
}

func init() {
  C.v8_callback_init()
}

type V8Context struct {
	id        uint32
	v8context unsafe.Pointer
	funcs     map[string]func(...interface{}) interface{}
}

func NewContext() *V8Context {
	v := &V8Context{
		uint32(len(contexts)),
		C.v8_create(),
		make(map[string]func(...interface{}) interface{}),
	}
	contexts[v.id] = v
	runtime.SetFinalizer(v, func(p *V8Context) {
		C.v8_release(p.v8context)
	})
	return v
}

func (v *V8Context) MustEval(in string) (res interface{}) {
	res, err := v.Eval(in)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (v *V8Context) Eval(in string) (res interface{}, err error) {
	ptr := C.CString(in)
	defer C.free(unsafe.Pointer(ptr))
	ret := C.v8_execute(v.v8context, ptr)
	if ret != nil {
		out := C.GoString(ret)
		if out != "" {
			C.free(unsafe.Pointer(ret))
			var buf bytes.Buffer
			buf.Write([]byte(out))
			dec := json.NewDecoder(&buf)
			err = dec.Decode(&res)
			return
		}
		return nil, nil
	}
	ret = C.v8_error(v.v8context)
	out := C.GoString(ret)
	C.free(unsafe.Pointer(ret))
	return nil, errors.New(out)
}

func (v *V8Context) AddFunc(name string, f func(...interface{}) interface{}) error {
	v.funcs[name] = f
	b := bytes.NewBufferString("")
	tmpl.Execute(b, map[string]interface{}{
		"id":   v.id,
		"name": name,
	})
	_, err := v.Eval(b.String())
	return err
}
