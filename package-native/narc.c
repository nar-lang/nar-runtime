#include "nar.h"
#include "narc.h"

init_data_t *env = NULL;

void nar_init_lib(init_data_t *init_data) {
    env = init_data;
}

void nar_register_def(nar_string_t module_name, nar_string_t def_name, nar_object_t def) {
    (*env->register_def)(env->runtime, module_name, def_name, def);
}

nar_object_t nar_fail(nar_string_t message) {
    env->nar_fail(env->runtime, message);
    return nar_unit();
}

nar_object_t nar_apply(nar_object_t fn, nar_size_t num_args, nar_object_t *args) {
    return env->nar_apply(env->runtime, fn, num_args, args);
}

nar_object_kind_t nar_object_kind(nar_object_t obj) {
    return env->nar_object_kind(env->runtime, obj);
}

nar_object_t nar_unit(void) {
    return env->nar_unit(env->runtime);
}

nar_object_t nar_char(nar_char_t value) {
    return env->nar_char(env->runtime, value);
}

nar_object_t nar_int(nar_int_t value) {
    return env->nar_int(env->runtime, value);
}

nar_object_t nar_float(nar_float_t value) {
    return env->nar_float(env->runtime, value);
}

nar_object_t nar_string(nar_string_t value) {
    return env->nar_string(env->runtime, value);
}

nar_object_t nar_record(nar_size_t size, nar_string_t *keys, nar_object_t *values) {
    return env->nar_record(env->runtime, size, keys, values);
}

nar_object_t nar_list(nar_size_t size, nar_object_t *items) {
    return env->nar_list(env->runtime, size, items);
}

nar_object_t nar_list_cons(nar_object_t head, nar_object_t tail) {
    return env->nar_list_cons(env->runtime, head, tail);
}

nar_object_t nar_tuple(nar_size_t size, nar_object_t *items) {
    return env->nar_tuple(env->runtime, size, items);
}

nar_object_t nar_bool(nar_bool_t value) {
    return env->nar_bool(env->runtime, value);
}

nar_object_t nar_option(nar_string_t name, nar_size_t size, nar_object_t *items) {
    return env->nar_option(env->runtime, name, size, items);
}

nar_object_t nar_func(nar_ptr_t fn, nar_size_t arity) {
    return env->nar_func(env->runtime, fn, arity);
}

nar_object_t nar_native(nar_ptr_t ptr, nar_cmp_native_fn_t cmp) {
    return env->nar_native(env->runtime, ptr, cmp);
}

void nar_to_unit(nar_object_t obj) {
    env->nar_to_unit(env->runtime, obj);
}

nar_char_t nar_to_char(nar_object_t obj) {
    return env->nar_to_char(env->runtime, obj);
}

nar_int_t nar_to_int(nar_object_t obj) {
    return env->nar_to_int(env->runtime, obj);
}

nar_float_t nar_to_float(nar_object_t obj) {
    return env->nar_to_float(env->runtime, obj);
}

nar_string_t nar_to_string(nar_object_t obj) {
    return env->nar_to_string(env->runtime, obj);
}

nar_record_t nar_to_record(nar_object_t obj) {
    return env->nar_to_record(env->runtime, obj);
}

nar_list_t nar_to_list(nar_object_t obj) {
    return env->nar_to_list(env->runtime, obj);
}

nar_tuple_t nar_to_tuple(nar_object_t obj) {
    return env->nar_to_tuple(env->runtime, obj);
}

nar_bool_t nar_to_bool(nar_object_t obj) {
    return env->nar_to_bool(env->runtime, obj);
}

nar_option_t nar_to_option(nar_object_t obj) {
    return env->nar_to_option(env->runtime, obj);
}

nar_func_t nar_to_func(nar_object_t obj) {
    return env->nar_to_func(env->runtime, obj);
}

nar_native_t nar_to_native(nar_object_t obj) {
    return env->nar_to_native(env->runtime, obj);
}

void *nar_alloc(size_t size) {
    return env->nar_alloc(env->runtime, size);
}

void nar_print(nar_string_t message) {
    return env->nar_print(env->runtime, message);
}
