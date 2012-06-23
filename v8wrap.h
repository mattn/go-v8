#ifndef _V8WRAP_H_
#ifdef __cplusplus
extern "C" {
#endif
extern void v8_init(void*);
extern void* v8_create();
extern void v8_release(void* ctx);
extern char* v8_execute(void* ctx, char* str);
extern char* v8_error(void* ctx);

typedef char* (*v8wrap_callback)(unsigned int, char*, char*);
#ifdef __cplusplus
}
#endif
#endif /* _V8WRAP_H_ */
