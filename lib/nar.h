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
typedef uint8_t nar_bool_t;
#define nar_false 0
#define nar_true 1
typedef uint8_t nar_byte_t;
typedef wchar_t nar_char_t;
typedef int64_t nar_int_t;
typedef uint64_t nar_uint_t;
typedef double nar_float_t;
typedef wchar_t *nar_string_t;
typedef void *nar_ptr_t;
typedef nar_int_t (*nar_cmp_native_fn_t)(nar_runtime_t rt, nar_ptr_t a, nar_ptr_t b);
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

typedef void *nar_bytecode_t;

#endif // NAR_H
