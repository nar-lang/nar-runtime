#ifndef NAR_BASE_NARC_H
#define NAR_BASE_NARC_H

#include "nar.h"

/**
 * Called by runtime to initialize the native library.
 * It should register all the native functions.
 * Redirect init_data to nar_init_lib(...) to initialize helper functions.
 */
int init(init_data_t *init_data);

void nar_init_lib(init_data_t *init_data);
void nar_register_def(nar_string_t module_name, nar_string_t def_name, nar_object_t def);
nar_object_kind_t nar_object_kind(nar_object_t obj);
nar_object_t nar_fail(nar_string_t message);
nar_object_t nar_apply(nar_object_t fn, nar_size_t num_args, nar_object_t *args);
void *nar_alloc(size_t size);
void nar_print(nar_string_t message);

nar_object_t nar_unit(void);
nar_object_t nar_char(nar_char_t value);
nar_object_t nar_int(nar_int_t value);
nar_object_t nar_float(nar_float_t value);
nar_object_t nar_string(nar_string_t value);
nar_object_t nar_record(nar_size_t size, nar_string_t *keys, nar_object_t *values);
nar_object_t nar_list(nar_size_t size, nar_object_t *items);
nar_object_t nar_list_cons(nar_object_t head, nar_object_t tail);
nar_object_t nar_tuple(nar_size_t size, nar_object_t *items);
nar_object_t nar_bool(nar_bool_t value);
nar_object_t nar_option(nar_string_t name, nar_size_t size, nar_object_t *values);
nar_object_t nar_func(nar_ptr_t fn, nar_size_t arity);
nar_object_t nar_native(nar_ptr_t ptr, nar_cmp_native_fn_t cmp);

void nar_to_unit(nar_object_t obj);
nar_char_t nar_to_char(nar_object_t obj);
nar_int_t nar_to_int(nar_object_t obj);
nar_float_t nar_to_float(nar_object_t obj);
nar_string_t nar_to_string(nar_object_t obj);
nar_record_t nar_to_record(nar_object_t obj);
nar_list_t nar_to_list(nar_object_t obj);
nar_tuple_t nar_to_tuple(nar_object_t obj);
nar_bool_t nar_to_bool(nar_object_t obj);
nar_option_t nar_to_option(nar_object_t obj);
nar_func_t nar_to_func(nar_object_t obj);
nar_native_t nar_to_native(nar_object_t obj);

#endif //NAR_BASE_NARC_H
