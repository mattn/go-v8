#include <v8.h>
#include <iostream>
#include <sstream>
#include <cstring>
#include <cstdlib>
#include "v8wrap.h"

extern "C" {

char*
__strdup(const char* ptr) {
  int l = strlen(ptr);
  char* p = new char[l + 1];
  strcpy(p, ptr);
  return p;
}

static volatile v8wrap_callback ___go_v8_callback = NULL;

static std::string
to_json(v8::Handle<v8::Value> value) {
  v8::HandleScope scope(v8::Isolate::GetCurrent());
  v8::TryCatch try_catch;
  v8::Handle<v8::Object> json = v8::Handle<v8::Object>::Cast(
    v8::Isolate::GetCurrent()->GetCurrentContext()->Global()->Get(v8::String::NewFromUtf8(v8::Isolate::GetCurrent(), "JSON")));
  v8::Handle<v8::Function> func = v8::Handle<v8::Function>::Cast(
    json->GetRealNamedProperty(v8::String::NewFromUtf8(v8::Isolate::GetCurrent(), "stringify")));
  v8::Handle<v8::Value> args[1];
  args[0] = value;
  v8::String::Utf8Value ret(
    func->Call(v8::Isolate::GetCurrent()->GetCurrentContext()->Global(), 1, args)->ToString());
  return (char*) *ret;
}

v8::Handle<v8::Value>
from_json(std::string str) {
  v8::HandleScope scope(v8::Isolate::GetCurrent());
  v8::TryCatch try_catch;

  v8::Handle<v8::Object> json = v8::Handle<v8::Object>::Cast(
    v8::Isolate::GetCurrent()->GetCurrentContext()->Global()->Get(v8::String::NewFromUtf8(v8::Isolate::GetCurrent(), "JSON")));
  v8::Handle<v8::Function> func = v8::Handle<v8::Function>::Cast(
    json->GetRealNamedProperty(v8::String::NewFromUtf8(v8::Isolate::GetCurrent(), "parse")));
  v8::Handle<v8::Value> args[1];
  args[0] = v8::String::NewFromUtf8(v8::Isolate::GetCurrent(), str.c_str());
  return func->Call(v8::Isolate::GetCurrent()->GetCurrentContext()->Global(), 1, args);
}

v8data
v8_get_array_item(v8data* array, int index) {
  return array[index];
}

void
_go_call(const v8::FunctionCallbackInfo<v8::Value>& args) {
  v8::Locker locker(v8::Isolate::GetCurrent());
  uint32_t id = args[0]->ToUint32()->Value();
  v8::String::Utf8Value name(args[1]);

  // Parse arguments
  v8::Array* realArgs = v8::Array::Cast(*args[2]);
  v8data* data = (v8data*) malloc(sizeof(v8data) * realArgs->Length());

  for (int i = 0; i < realArgs->Length(); i++) {
    v8::Local<v8::Value> arg = realArgs->Get(i);

    v8::String::Utf8Value argString(arg);
    if (arg->IsRegExp()) {
      data[i].obj_type = v8regexp;
      data[i].repr = __strdup(*argString);
    } else if (arg->IsFunction()) {
      data[i].obj_type = v8function;
      data[i].repr = __strdup(*argString);
    } else if (arg->IsNumber()) {
      data[i].obj_type = v8number;
      data[i].repr = __strdup(*argString);
    } else if (arg->IsBoolean()) {
      data[i].obj_type = v8boolean;
      data[i].repr = __strdup(*argString);
    } else if (arg->IsString()) {
      data[i].obj_type = v8string;
      data[i].repr = __strdup(*argString);
    } else if (arg->IsArray()) {
      data[i].obj_type = v8array;
      data[i].repr = __strdup(*argString);
    } else if (arg->IsObject()) {
      data[i].obj_type = v8object;
      data[i].repr = __strdup(*argString);
    } else {
      data[i].obj_type = v8string;
      data[i].repr = __strdup(to_json(arg).c_str());
    }
  }

  v8::TryCatch try_catch;
  char* retv;
  retv = ___go_v8_callback(id, *name, data, realArgs->Length());

  // Free args memory
  for (int i = 0; i < realArgs->Length(); i++) {
      free(data[i].repr);
  }

  free(data);

  if (retv != NULL) {
    v8::Handle<v8::Value> ret = from_json(retv);
    free(retv);
    args.GetReturnValue().Set(ret);
  }
  args.GetReturnValue().Set(v8::Undefined(v8::Isolate::GetCurrent()));
}

class V8Context {
public:
  V8Context() : err_("") {
    v8::Locker locker(v8::Isolate::GetCurrent());
    v8::Isolate* isolate = v8::Isolate::GetCurrent();
    v8::HandleScope scope(isolate);
    v8::Handle<v8::ObjectTemplate> global = v8::ObjectTemplate::New();
    global->Set(v8::String::NewFromUtf8(v8::Isolate::GetCurrent(), "_go_call"),
      v8::FunctionTemplate::New(isolate, _go_call));

    v8::Local<v8::Context> context = v8::Context::New(isolate, NULL, global);
    context_.Reset(isolate, context);

  };

  virtual ~V8Context() {
    context_.Reset();
  };
  v8::Local<v8::Context> context() { return v8::Local<v8::Context>::New(v8::Isolate::GetCurrent(), context_); };
  const char* err() const { return err_.c_str(); };
  void err(const char* e) { this->err_ = std::string(e); }

private:
  v8::Persistent<v8::Context> context_;
  std::string err_;
};

void
v8_init(void *p) {
  ___go_v8_callback = (v8wrap_callback) p;
}

void*
v8_create() {
  return (void*) new V8Context();
}

void
v8_release(void* ctx) {
  delete static_cast<V8Context *>(ctx);
}

char*
v8_error(void* ctx) {
  V8Context *context = static_cast<V8Context *>(ctx);
  return __strdup(context->err());
}

static std::string
report_exception(v8::TryCatch& try_catch) {
  v8::Handle<v8::Message> message = try_catch.Message();
  v8::String::Utf8Value exception(try_catch.Exception());
  std::stringstream ss;
  if (message.IsEmpty()) {
    ss << *exception << std::endl;
  } else {
    v8::String::Utf8Value filename(message->GetScriptResourceName());
    const char* filename_string = *filename;
    int linenum = message->GetLineNumber();
    ss
      << filename_string
      << ":" << linenum
      << ": " << *exception << std::endl;
    v8::String::Utf8Value sourceline(message->GetSourceLine());
    ss << *sourceline << std::endl;
    int start = message->GetStartColumn();
    for (int n = 0; n < start; n++) {
      ss << " ";
    }
    int end = message->GetEndColumn();
    for (int n = start; n < end; n++) {
      ss << "^";
    }
    ss << std::endl;
    v8::String::Utf8Value stack_trace(try_catch.StackTrace());
    if (stack_trace.length() > 0) {
      const char* stack_trace_string = *stack_trace;
      ss << stack_trace_string << std::endl;
    }
  }
  return ss.str();
}

char*
v8_execute(void *ctx, char* source) {
  v8::Locker locker(v8::Isolate::GetCurrent());
  V8Context *context = static_cast<V8Context *>(ctx);
  v8::HandleScope scope(v8::Isolate::GetCurrent());
  v8::TryCatch try_catch;

  v8::Context::Scope context_scope(context->context());

  context->err("");
  v8::ScriptCompiler::Source src(v8::String::NewFromUtf8(v8::Isolate::GetCurrent(), source));
  v8::Handle<v8::Script> script
    = v8::ScriptCompiler::Compile(v8::Isolate::GetCurrent(), &src);
  if (script.IsEmpty()) {
    //v8::Isolate::ThrowException(try_catch.Exception());
    context->err(report_exception(try_catch).c_str());
    return NULL;
  } else {
    v8::Handle<v8::Value> result = script->Run();
    if (result.IsEmpty()) {
      //v8::ThrowException(try_catch.Exception());
      context->err(report_exception(try_catch).c_str());
      return NULL;
    } else if (result->IsUndefined()) {
      return __strdup("");
    } else if (result->IsFunction()) {
      v8::Handle<v8::Function> func = v8::Handle<v8::Function>::Cast(result);
      v8::String::Utf8Value ret(func->ToString());
      return __strdup(*ret);
    } else if (result->IsRegExp()) {
      v8::Handle<v8::RegExp> re = v8::Handle<v8::RegExp>::Cast(result);
      v8::String::Utf8Value ret(re->ToString());
      return __strdup(*ret);
    } else {
      return __strdup(to_json(result).c_str());
    }
  }
}

}

// vim:set et sw=2 ts=2 ai:
