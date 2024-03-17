#include <stdlib.h>
#include "_cgo_export.h"
#include "../lib/nar.h"
#include "../lib/nar-package.h"


nar_int_t init_wrapper(void* init_fn, nar_runtime_t runtime) {
	nar_t *nar = (nar_t*)malloc(sizeof(nar_t));

  nar->register_def = &xnar_register_def;
  nar->fail = &xnar_fail;
  nar->apply_func = &xnar_apply_func;
  nar->alloc = &xnar_alloc;
  nar->print = &xnar_print;
  nar->get_object_kind = &xnar_get_object_kind;
  nar->new_unit = &xnar_new_unit;
  nar->new_char = &xnar_new_char;
  nar->new_int = &xnar_new_int;
  nar->new_float = &xnar_new_float;
  nar->new_string = &xnar_new_string;
  nar->new_record = &xnar_new_record;
  nar->new_list = &xnar_new_list;
  nar->new_list_cons = &xnar_new_list_cons;
  nar->new_tuple = &xnar_new_tuple;
  nar->new_bool = &xnar_new_bool;
  nar->new_option = &xnar_new_option;
  nar->new_func = &xnar_new_func;
  nar->new_native = &xnar_new_native;
  nar->to_unit = &xnar_to_unit;
  nar->to_char = &xnar_to_char;
  nar->to_int = &xnar_to_int;
  nar->to_float = &xnar_to_float;
  nar->to_string = &xnar_to_string;
  nar->to_record = &xnar_to_record;
  nar->to_list = &xnar_to_list;
  nar->to_tuple = &xnar_to_tuple;
  nar->to_bool = &xnar_to_bool;
  nar->to_option = &xnar_to_option;
  nar->to_func = &xnar_to_func;
  nar->to_native = &xnar_to_native;

	return ((init_fn_t)init_fn)(nar, runtime);
}
