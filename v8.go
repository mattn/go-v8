package v8

/*
#include <stdlib.h>
extern void* v8_create();
extern void v8_release(void* ctx);
extern char* v8_execute(void* ctx, char* str);
*/
import "C"
import (
	"errors"
	"unsafe"
)
import "runtime"
import "encoding/json"

import "bytes"

type V8Context struct {
	v8context unsafe.Pointer
}

func NewContext() *V8Context {
	v := &V8Context{C.v8_create()}
	runtime.SetFinalizer(v, func(p *V8Context) {
		C.v8_release(p.v8context)
	})
	return v
}

func (v *V8Context) Eval(in string) (res interface{}, err error) {
	ptr := C.CString(in)
	defer C.free(unsafe.Pointer(ptr))
	ret := C.v8_execute(v.v8context, ptr)
	if ret != nil {
		out := C.GoString(ret)
		var buf bytes.Buffer
		buf.Write([]byte(out))
		dec := json.NewDecoder(&buf)
		err = dec.Decode(&res)
		return
	}
	return nil, errors.New("failed to eval")
}
