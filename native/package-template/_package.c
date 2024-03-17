#include "_package.h"

int init(nar_t *n, nar_runtime_t rt) {
    nar = n;
    register_hello(rt);
    return 0;
}
