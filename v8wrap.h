#ifndef _V8WRAP_H_
#define _V8WRAP_H_

#ifdef __cplusplus
extern "C" {
#endif

extern void v8_init(void*);
extern void* v8_create();
extern void v8_release(void* ctx);
extern char* v8_execute(void* ctx, char* str);
extern char* v8_error(void* ctx);

typedef enum {
  v8regexp, v8string, v8function
} v8type;

typedef struct {
  v8type obj_type;
  char* repr; // String representation of object
} v8data;

extern char* _go_v8_callback(unsigned int contextId, char* functionName, v8data* v8Objects, int count);

extern v8data v8_get_array_item(v8data* array, int index);

#ifdef __cplusplus
}
#endif

#endif /* _V8WRAP_H_ */
