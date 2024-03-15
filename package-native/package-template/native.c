#include "native.h"

int init(void *init_data) {
    nar_init_lib(init_data);
    register_hello();
    return 0;
}
