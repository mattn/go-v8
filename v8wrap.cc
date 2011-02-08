#include <v8.h>
#include <string.h>

extern "C" {

class V8Context {
public:
    V8Context() {
		v8::HandleScope scope;

		v8::Handle<v8::ObjectTemplate> global = v8::ObjectTemplate::New();
		v8::Handle<v8::Context> context = v8::Context::New(NULL, global);

        context_ = v8::Persistent<v8::Context>::New(context);
    };

    virtual ~V8Context() {
        context_.Dispose();
    };

	v8::Handle<v8::Context> context() { return context_; };

private:
	v8::Persistent<v8::Context> context_;
};

void*
v8_create() {
	return (void*) new V8Context; 
}

void
v8_release(void* ctx) {
	delete  static_cast<V8Context *>(ctx);
}

char*
v8_execute(V8Context *ctx, char* source) {
    V8Context *context = static_cast<V8Context *>(ctx);
	v8::HandleScope scope;
	v8::TryCatch try_catch;

	v8::Context::Scope context_scope(context->context());

	v8::Handle<v8::Script> script
        = v8::Script::Compile(v8::String::New(source), v8::Undefined());
    if (script.IsEmpty()) {
		v8::ThrowException(try_catch.Exception());
        return NULL;
    }
    else {
		v8::Handle<v8::Value> result = script->Run();
        if (result.IsEmpty()) {
			v8::ThrowException(try_catch.Exception());
			return NULL;
        }
        else {
			v8::Handle<v8::Object> json = v8::Handle<v8::Object>::Cast(
					context->context()->Global()->Get(
						v8::String::New("JSON")));
			v8::Handle<v8::Function> func = v8::Handle<v8::Function>::Cast(
					json->GetRealNamedProperty(
							v8::String::New("stringify")));
			v8::Handle<v8::Value> args[1];
			args[0] = result;
			v8::String::Utf8Value ret(
					func->Call(context->context()->Global(), 1, args)
						->ToString());
			return strdup(*ret);
        }
    }
}

}
