#include "native.h"
#include <narc.c>

nar_object_t hello_world() {
    return nar_string(L"Hello, World!");
}

void register_world(void) {
    nar_string_t module_name = L"My.Hello";
    nar_register_def(module_name, L"world", nar_func(&hello_world, 0));
}
