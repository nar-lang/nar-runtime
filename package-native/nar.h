#ifndef NAR_H
#define NAR_H

#include <stdint.h>
#include <wchar.h>

#if __DBL_DIG__ != 15 || __DBL_MANT_DIG__ != 53 || __DBL_MIN_EXP__ != -1021 || __DBL_MAX_EXP__ != 1024
#error "IEEE 754 floating point support is required"
#endif

typedef uint32_t nar_runtime_t;

typedef uint64_t nar_size_t;

typedef uint8_t nar_bool_t;
#define nar_false 0
#define nar_true 1
typedef int64_t nar_int_t;
typedef uint64_t nar_uint_t;
typedef wchar_t nar_char_t;
typedef double nar_float_t;
typedef wchar_t *nar_string_t;
typedef void *nar_ptr_t;
typedef nar_int_t (*nar_cmp_native_fn_t)(nar_ptr_t a, nar_ptr_t b);
typedef struct {
    nar_ptr_t ptr;
    nar_cmp_native_fn_t cmp;
} nar_native_t;

typedef uint64_t nar_object_t;
typedef struct {
    nar_size_t size;
    nar_string_t *keys;
    nar_object_t *values;
} nar_record_t;
typedef struct {
    nar_size_t size;
    nar_object_t *items;
} nar_list_t;
typedef nar_list_t nar_tuple_t;
typedef struct {
    nar_string_t name;
    nar_size_t size;
    nar_object_t *values;
} nar_option_t;
typedef struct {
    nar_ptr_t ptr;
    nar_size_t arity;
} nar_func_t;

typedef enum {
    NAR_UNKNOWN,
    NAR_UNIT,
    NAR_INT,
    NAR_FLOAT,
    NAR_STRING,
    NAR_CHAR,
    NAR_RECORD,
    NAR_TUPLE,
    NAR_LIST,
    NAR_OPTION,
    NAR_FUNCTION,
    NAR_CLOSURE,
    NAR_NATIVE,
} nar_object_kind_t;

typedef struct init_data_t {
    nar_runtime_t runtime;

    void (*register_def)(nar_runtime_t runtime, nar_string_t module_name, nar_string_t def_name, nar_object_t def);
    nar_object_kind_t (*nar_object_kind)(nar_runtime_t runtime, nar_object_t obj);
    void (*nar_fail)(nar_runtime_t runtime, nar_string_t message);
    nar_object_t (*nar_apply)(nar_runtime_t runtime, nar_object_t fn, nar_size_t num_args, nar_object_t *args);
    void *(*nar_alloc)(nar_runtime_t runtime, nar_size_t size);
    void (*nar_print)(nar_runtime_t runtime, nar_string_t message);

    nar_object_t (*nar_unit)(nar_runtime_t runtime);
    nar_object_t (*nar_char)(nar_runtime_t runtime, nar_char_t value);
    nar_object_t (*nar_int)(nar_runtime_t runtime, nar_int_t value);
    nar_object_t (*nar_float)(nar_runtime_t runtime, nar_float_t value);
    nar_object_t (*nar_string)(nar_runtime_t runtime, nar_string_t value);
    nar_object_t (*nar_record)(nar_runtime_t runtime, nar_size_t size, nar_string_t *keys, nar_object_t *values);
    nar_object_t (*nar_list)(nar_runtime_t runtime, nar_size_t size, nar_object_t *items);
    nar_object_t (*nar_list_cons)(nar_runtime_t runtime, nar_object_t head, nar_object_t tail);
    nar_object_t (*nar_tuple)(nar_runtime_t runtime, nar_size_t size, nar_object_t *items);
    nar_object_t (*nar_bool)(nar_runtime_t runtime, nar_bool_t value);
    nar_object_t (*nar_option)(nar_runtime_t runtime, nar_string_t name, nar_size_t size, nar_object_t *items);
    nar_object_t (*nar_func)(nar_runtime_t runtime, nar_ptr_t fn, nar_size_t arity);
    nar_object_t (*nar_native)(nar_runtime_t runtime, nar_ptr_t ptr, nar_cmp_native_fn_t size);

    void (*nar_to_unit)(nar_runtime_t runtime, nar_object_t obj);
    nar_char_t (*nar_to_char)(nar_runtime_t runtime, nar_object_t obj);
    nar_int_t (*nar_to_int)(nar_runtime_t runtime, nar_object_t obj);
    nar_float_t (*nar_to_float)(nar_runtime_t runtime, nar_object_t obj);
    nar_string_t (*nar_to_string)(nar_runtime_t runtime, nar_object_t obj);
    nar_record_t (*nar_to_record)(nar_runtime_t runtime, nar_object_t obj);
    nar_list_t (*nar_to_list)(nar_runtime_t runtime, nar_object_t obj);
    nar_tuple_t (*nar_to_tuple)(nar_runtime_t runtime, nar_object_t obj);
    nar_bool_t (*nar_to_bool)(nar_runtime_t runtime, nar_object_t obj);
    nar_option_t (*nar_to_option)(nar_runtime_t runtime, nar_object_t obj);
    nar_func_t (*nar_to_func)(nar_runtime_t runtime, nar_object_t obj);
    nar_native_t (*nar_to_native)(nar_runtime_t runtime, nar_object_t obj);
} init_data_t;
#endif
