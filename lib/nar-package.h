#ifndef NATIVE_PACKAGE_H
#define NATIVE_PACKAGE_H

#include "nar.h"

typedef struct {
    void              (*register_def)   (nar_runtime_t runtime, nar_string_t module_name, nar_string_t def_name, nar_object_t def);
    void              (*fail)           (nar_runtime_t runtime, nar_string_t message);
    nar_object_t      (*apply_func)     (nar_runtime_t runtime, nar_object_t fn, nar_size_t num_args, nar_object_t *args);
    void *            (*alloc)          (nar_runtime_t runtime, nar_size_t size);
    void              (*print)          (nar_runtime_t runtime, nar_string_t message);
    nar_object_kind_t (*get_object_kind)(nar_runtime_t runtime, nar_object_t obj);

    nar_object_t (*new_unit)      (nar_runtime_t runtime);
    nar_object_t (*new_char)      (nar_runtime_t runtime, nar_char_t value);
    nar_object_t (*new_int)       (nar_runtime_t runtime, nar_int_t value);
    nar_object_t (*new_float)     (nar_runtime_t runtime, nar_float_t value);
    nar_object_t (*new_string)    (nar_runtime_t runtime, nar_string_t value);
    nar_object_t (*new_record)    (nar_runtime_t runtime, nar_size_t size, nar_string_t *keys, nar_object_t *values);
    nar_object_t (*new_list)      (nar_runtime_t runtime, nar_size_t size, nar_object_t *items);
    nar_object_t (*new_list_cons) (nar_runtime_t runtime, nar_object_t head, nar_object_t tail);
    nar_object_t (*new_tuple)     (nar_runtime_t runtime, nar_size_t size, nar_object_t *items);
    nar_object_t (*new_bool)      (nar_runtime_t runtime, nar_bool_t value);
    nar_object_t (*new_option)    (nar_runtime_t runtime, nar_string_t name, nar_size_t size, nar_object_t *items);
    nar_object_t (*new_func)      (nar_runtime_t runtime, nar_ptr_t fn, nar_size_t arity);
    nar_object_t (*new_native)    (nar_runtime_t runtime, nar_ptr_t ptr, nar_cmp_native_fn_t size);

    void          (*to_unit)    (nar_runtime_t runtime, nar_object_t obj);
    nar_char_t    (*to_char)    (nar_runtime_t runtime, nar_object_t obj);
    nar_int_t     (*to_int)     (nar_runtime_t runtime, nar_object_t obj);
    nar_float_t   (*to_float)   (nar_runtime_t runtime, nar_object_t obj);
    nar_string_t  (*to_string)  (nar_runtime_t runtime, nar_object_t obj);
    nar_record_t  (*to_record)  (nar_runtime_t runtime, nar_object_t obj);
    nar_list_t    (*to_list)    (nar_runtime_t runtime, nar_object_t obj);
    nar_tuple_t   (*to_tuple)   (nar_runtime_t runtime, nar_object_t obj);
    nar_bool_t    (*to_bool)    (nar_runtime_t runtime, nar_object_t obj);
    nar_option_t  (*to_option)  (nar_runtime_t runtime, nar_object_t obj);
    nar_func_t    (*to_func)    (nar_runtime_t runtime, nar_object_t obj);
    nar_native_t  (*to_native)  (nar_runtime_t runtime, nar_object_t obj);
} nar_t;

typedef nar_int_t (*init_fn_t)(nar_t *, nar_runtime_t);

#endif //NATIVE_PACKAGE_H
