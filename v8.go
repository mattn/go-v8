package v8

/*
#include <stdlib.h>
#include "v8wrap.h"

extern char* _go_v8_callback(unsigned int id, char* n, char* a);

static char*
_c_v8_callback(unsigned int id, char* n, char* a) {
	return _go_v8_callback(id, n, a);
}

static void
v8_callback_init() {
	v8_init((void*) _c_v8_callback);
}
*/
// #cgo LDFLAGS: -L. -lv8wrap -lstdc++
import "C"
import (
	"bytes"
	"encoding/json"
	"errors"
	"runtime"
	"text/template"
	"unsafe"
)

var contexts = make(map[uint32]*V8Context)

var tmpl = template.Must(template.New("go-v8").Parse(`
function {{.name}}() {
  return _go_call({{.id}}, "{{.name}}", JSON.stringify([].slice.call(arguments)));
}`))

//export _go_v8_callback
func _go_v8_callback(id uint32, n, a *C.char) *C.char {
	c := contexts[id]
	f := c.funcs[C.GoString(n)]
	if f != nil {
		var argv []interface{}
		json.Unmarshal([]byte(C.GoString(a)), &argv)
		ret := f(argv...)
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

func (v *V8Context) Eval(in string) (res interface{}, err error) {
	ptr := C.CString(in)
	defer C.free(unsafe.Pointer(ptr))
	C.v8_callback_init()
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
