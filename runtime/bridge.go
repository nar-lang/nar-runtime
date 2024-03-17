package runtime

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>
#include "../lib/nar.h"
#include "../lib/nar-package.h"

nar_int_t init_wrapper(void* init_fn, nar_runtime_t runtime);
*/
import "C"
import (
	"fmt"
	"os"
	"path/filepath"
	goruntime "runtime"
	"unsafe"
)

func (rt *Runtime) registerPackage(packageName string, packageVersion int, libsPath string) error {
	var initPtr unsafe.Pointer
	if goruntime.GOOS == "darwin" {
		libPath := filepath.Join(libsPath, fmt.Sprintf("lib%s.%d.dylib", packageName, packageVersion))
		if _, err := os.Stat(libPath); err != nil {
			return nil
		}

		handle := C.dlopen(C.CString(libPath), C.RTLD_LAZY)
		if handle == nil {
			return fmt.Errorf("dlopen failed: %s", (string)(C.GoString(C.dlerror())))
		}
		initPtr = C.dlsym(handle, C.CString("init"))
		if initPtr == nil {
			return fmt.Errorf("dlsym failed: %s", (string)(C.GoString(C.dlerror())))
		}
	}

	result := C.init_wrapper(initPtr, C.nar_runtime_t(rt.Id()))
	if result != 0 {
		return fmt.Errorf("failed to init package %s with code %d", packageName, int(result))
	}
	return nil
}
