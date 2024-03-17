#include "_package.h"

nar_object_t hello_world(nar_runtime_t rt) {
    return nar->string(rt, L"Hello, World!");
}

void register_world(nar_runtime_t rt) {
    nar_string_t module_name = L"My.Hello";
    nar->register_def(rt, module_name, L"world", nar_func(&hello_world, 0));
}
