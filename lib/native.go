package main

/*
#include <stdlib.h>
#include "nar.h"

typedef nar_object_t (*func0)(nar_runtime_t);
typedef nar_object_t (*func1)(nar_runtime_t, nar_object_t);
typedef nar_object_t (*func2)(nar_runtime_t, nar_object_t, nar_object_t);
typedef nar_object_t (*func3)(nar_runtime_t, nar_object_t, nar_object_t, nar_object_t);
typedef nar_object_t (*func4)(nar_runtime_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t);
typedef nar_object_t (*func5)(nar_runtime_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t);
typedef nar_object_t (*func6)(nar_runtime_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t);
typedef nar_object_t (*func7)(nar_runtime_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t);
typedef nar_object_t (*func8)(nar_runtime_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t, nar_object_t);

static nar_object_t call_func0(func0 fn, nar_runtime_t rt) {
	return fn(rt);
}
static nar_object_t call_func1(func1 fn, nar_runtime_t rt, nar_object_t a) {
	return fn(rt, a);
}
static nar_object_t call_func2(func2 fn, nar_runtime_t rt, nar_object_t a, nar_object_t b) {
	return fn(rt, a, b);
}
static nar_object_t call_func3(func3 fn, nar_runtime_t rt, nar_object_t a, nar_object_t b, nar_object_t c) {
	return fn(rt, a, b, c);
}
static nar_object_t call_func4(func4 fn, nar_runtime_t rt, nar_object_t a, nar_object_t b, nar_object_t c, nar_object_t d) {
	return fn(rt, a, b, c, d);
}
static nar_object_t call_func5(func5 fn, nar_runtime_t rt, nar_object_t a, nar_object_t b, nar_object_t c, nar_object_t d, nar_object_t e) {
	return fn(rt, a, b, c, d, e);
}
static nar_object_t call_func6(func6 fn, nar_runtime_t rt, nar_object_t a, nar_object_t b, nar_object_t c, nar_object_t d, nar_object_t e, nar_object_t f) {
	return fn(rt, a, b, c, d, e, f);
}
static nar_object_t call_func7(func7 fn, nar_runtime_t rt, nar_object_t a, nar_object_t b, nar_object_t c, nar_object_t d, nar_object_t e, nar_object_t f, nar_object_t g) {
	return fn(rt, a, b, c, d, e, f, g);
}
static nar_object_t call_func8(func8 fn, nar_runtime_t rt, nar_object_t a, nar_object_t b, nar_object_t c, nar_object_t d, nar_object_t e, nar_object_t f, nar_object_t g, nar_object_t h) {
	return fn(rt, a, b, c, d, e, f, g, h);
}

*/
import "C"
import (
	"bytes"
	"fmt"
	"github.com/nar-lang/nar-common/ast"
	"github.com/nar-lang/nar-common/bytecode"
	"github.com/nar-lang/nar-runtime/runtime"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	goruntime "runtime"
	"unsafe"
)

func main() {}

//export nar_fail
func nar_fail(rtid C.nar_runtime_t, message C.nar_string_t) {
	msg := stringToGo(message)
	setLastError(runtime.GetById(runtime.Id(rtid)), msg)
	if runtime.KDebug {
		panic(msg)
	}
}

//export nar_get_last_error
func nar_get_last_error(rtid C.nar_runtime_t) C.nar_string_t {
	rt := runtime.GetById(runtime.Id(rtid))
	return stringToC(getLastError(rt))
}

//export nar_load_bytecode
func nar_load_bytecode(size C.nar_size_t, data *C.nar_byte_t) C.nar_bytecode_t {
	buf := (*[1 << 30]byte)(unsafe.Pointer(data))
	btc, err := bytecode.Read(bytes.NewReader(buf[:size]))
	if err != nil {
		setLastError(nil, "failed to load bytecode")
		return C.nar_bytecode_t(nil)
	}
	return C.nar_bytecode_t(btc)
}

//export nar_create_runtime
func nar_create_runtime(btc C.nar_bytecode_t, libsPath C.nar_string_t) C.nar_runtime_t {
	rt, err := runtime.NewRuntime((*bytecode.Binary)(btc), string(stringToGo(libsPath)))
	if err != nil {
		setLastError(nil, runtime.TString(err.Error()))
		return 0
	}
	return C.nar_runtime_t(rt.Id())
}

//export nar_destroy_runtime
func nar_destroy_runtime(rtid C.nar_runtime_t) {
	nar_free_frame_memory(rtid)
	rt := runtime.GetById(runtime.Id(rtid))
	rt.Destroy()
}

//export nar_register_def
func nar_register_def(rtid C.nar_runtime_t, module_name C.nar_string_t, def_name C.nar_string_t, def C.nar_object_t) {
	rt := runtime.GetById(runtime.Id(rtid))
	rt.RegisterDef(bytecode.QualifiedIdentifier(stringToGo(module_name)), ast.Identifier(stringToGo(def_name)), runtime.Object(def))
}

//export nar_apply
func nar_apply(rtid C.nar_runtime_t, name C.nar_string_t, num_args C.nar_size_t, args *C.nar_object_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	args_ := make([]runtime.Object, 0, num_args)
	argsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(args))
	for i := C.nar_size_t(0); i < num_args; i++ {
		args_ = append(args_, runtime.Object(argsSlice[i]))
	}
	res, err := rt.Apply(bytecode.FullIdentifier(stringToGo(name)), args_...)
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	return C.nar_object_t(res)
}

//export nar_apply_func
func nar_apply_func(rtid C.nar_runtime_t, fn C.nar_object_t, num_args C.nar_size_t, args *C.nar_object_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	args_ := make([]runtime.Object, 0, num_args)
	argsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(args))
	for i := C.nar_size_t(0); i < num_args; i++ {
		args_ = append(args_, runtime.Object(argsSlice[i]))
	}
	res, err := rt.ApplyFunc(runtime.Object(fn), args_...)
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	return C.nar_object_t(res)
}

//export nar_get_object_kind
func nar_get_object_kind(_ C.nar_runtime_t, obj C.nar_object_t) C.nar_object_kind_t {
	return C.nar_object_kind_t(runtime.Object(obj).Kind())
}

//export nar_new_unit
func nar_new_unit(rtid C.nar_runtime_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	return C.nar_object_t(rt.NewUnit())
}

//export nar_new_char
func nar_new_char(rtid C.nar_runtime_t, value C.nar_char_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	return C.nar_object_t(rt.NewChar(runtime.TChar(value)))
}

//export nar_new_int
func nar_new_int(rtid C.nar_runtime_t, value C.nar_int_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	return C.nar_object_t(rt.NewInt(runtime.TInt(value)))
}

//export nar_new_float
func nar_new_float(rtid C.nar_runtime_t, value C.nar_float_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	return C.nar_object_t(rt.NewFloat(runtime.TFloat(value)))
}

//export nar_new_string
func nar_new_string(rtid C.nar_runtime_t, value C.nar_string_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	return C.nar_object_t(rt.NewString(stringToGo(value)))
}

//export nar_new_record
func nar_new_record(rtid C.nar_runtime_t, size C.nar_size_t, keys *C.nar_string_t, values *C.nar_object_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	fields_ := make([]runtime.TString, 0, size)
	keys_ := make([]runtime.Object, 0, size)
	keysSlice := (*[1 << 30]C.nar_string_t)(unsafe.Pointer(keys))
	valuesSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(values))
	for i := C.nar_size_t(0); i < size; i++ {
		fields_ = append(fields_, stringToGo(keysSlice[i]))
		keys_ = append(keys_, runtime.Object(valuesSlice[i]))
	}

	return C.nar_object_t(rt.NewRecord(fields_, keys_))
}

//export nar_new_list
func nar_new_list(rtid C.nar_runtime_t, size C.nar_size_t, items *C.nar_object_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	if size <= 0 {
		return C.nar_object_t(rt.NewList())
	}
	items_ := make([]runtime.Object, 0, size)
	itemsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(items))
	for i := C.nar_size_t(0); i < size; i++ {
		items_ = append(items_, runtime.Object(itemsSlice[i]))
	}
	return C.nar_object_t(rt.NewList(items_...))
}

//export nar_new_list_cons
func nar_new_list_cons(rtid C.nar_runtime_t, head C.nar_object_t, tail C.nar_object_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	return C.nar_object_t(rt.NewListItem(runtime.Object(head), runtime.Object(tail)))
}

//export nar_new_tuple
func nar_new_tuple(rtid C.nar_runtime_t, size C.nar_size_t, items *C.nar_object_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	items_ := make([]runtime.Object, 0, size)
	itemsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(items))
	for i := C.nar_size_t(0); i < size; i++ {
		items_ = append(items_, runtime.Object(itemsSlice[i]))
	}
	return C.nar_object_t(rt.NewTuple(items_...))
}

//export nar_new_bool
func nar_new_bool(rtid C.nar_runtime_t, value C.nar_bool_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	return C.nar_object_t(rt.NewBool(value != 0))
}

//export nar_new_option
func nar_new_option(rtid C.nar_runtime_t, name C.nar_string_t, size C.nar_size_t, items *C.nar_object_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	items_ := make([]runtime.Object, 0, size)
	itemsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(items))
	for i := C.nar_size_t(0); i < size; i++ {
		items_ = append(items_, runtime.Object(itemsSlice[i]))
	}
	optName := stringToGo(name)
	return C.nar_object_t(rt.NewOption(optName, items_...))
}

var funcWrappers = map[C.nar_ptr_t]unsafe.Pointer{}

//export nar_new_func
func nar_new_func(rtid C.nar_runtime_t, fn C.nar_ptr_t, arity C.nar_size_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	wrapper, ok := funcWrappers[fn]
	if !ok {
		switch int(arity) {
		case 0:
			{
				wfn := func() runtime.Object {
					return runtime.Object(C.call_func0(C.func0(fn), rtid))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 1:
			{
				wfn := func(a runtime.Object) runtime.Object {
					return runtime.Object(C.call_func1(C.func1(fn), rtid, C.nar_object_t(a)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 2:
			{
				wfn := func(a runtime.Object, b runtime.Object) runtime.Object {
					return runtime.Object(C.call_func2(C.func2(fn), rtid, C.nar_object_t(a), C.nar_object_t(b)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 3:
			{
				wfn := func(a runtime.Object, b runtime.Object, c runtime.Object) runtime.Object {
					return runtime.Object(C.call_func3(C.func3(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 4:
			{
				wfn := func(a runtime.Object, b runtime.Object, c runtime.Object, d runtime.Object) runtime.Object {
					return runtime.Object(C.call_func4(C.func4(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c), C.nar_object_t(d)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 5:
			{
				wfn := func(a runtime.Object, b runtime.Object, c runtime.Object, d runtime.Object, e runtime.Object) runtime.Object {
					return runtime.Object(C.call_func5(C.func5(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c), C.nar_object_t(d), C.nar_object_t(e)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 6:
			{
				wfn := func(a runtime.Object, b runtime.Object, c runtime.Object, d runtime.Object, e runtime.Object, f runtime.Object) runtime.Object {
					return runtime.Object(C.call_func6(C.func6(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c), C.nar_object_t(d), C.nar_object_t(e), C.nar_object_t(f)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 7:
			{
				wfn := func(a runtime.Object, b runtime.Object, c runtime.Object, d runtime.Object, e runtime.Object, f runtime.Object, g runtime.Object) runtime.Object {
					return runtime.Object(C.call_func7(C.func7(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c), C.nar_object_t(d), C.nar_object_t(e), C.nar_object_t(f), C.nar_object_t(g)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 8:
			{
				wfn := func(a runtime.Object, b runtime.Object, c runtime.Object, d runtime.Object, e runtime.Object, f runtime.Object, g runtime.Object, h runtime.Object) runtime.Object {
					return runtime.Object(C.call_func8(C.func8(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c), C.nar_object_t(d), C.nar_object_t(e), C.nar_object_t(f), C.nar_object_t(g), C.nar_object_t(h)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		}
		funcWrappers[fn] = wrapper
	}
	return C.nar_object_t(rt.NewFunc(wrapper, uint8(arity)))
}

//export nar_new_native
func nar_new_native(rtid C.nar_runtime_t, ptr C.nar_ptr_t, cmp C.nar_cmp_native_fn_t) C.nar_object_t {
	rt := runtime.GetById(runtime.Id(rtid))
	return C.nar_object_t(rt.NewNative(unsafe.Pointer(ptr), unsafe.Pointer(cmp)))
}

//export nar_to_unit
func nar_to_unit(rtid C.nar_runtime_t, obj C.nar_object_t) {
	rt := runtime.GetById(runtime.Id(rtid))
	_, err := runtime.ToUnit(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
}

//export nar_to_char
func nar_to_char(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_char_t {
	rt := runtime.GetById(runtime.Id(rtid))
	v, err := runtime.ToChar(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	return C.nar_char_t(v)
}

//export nar_to_int
func nar_to_int(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_int_t {
	rt := runtime.GetById(runtime.Id(rtid))
	v, err := runtime.ToInt(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	return C.nar_int_t(v)
}

//export nar_to_float
func nar_to_float(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_float_t {
	rt := runtime.GetById(runtime.Id(rtid))
	v, err := runtime.ToFloat(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	return C.nar_float_t(v)
}

//export nar_to_string
func nar_to_string(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_string_t {
	rt := runtime.GetById(runtime.Id(rtid))
	v, err := runtime.ToString(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	return stringToC(v)
}

//export nar_to_record
func nar_to_record(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_record_t {
	rt := runtime.GetById(runtime.Id(rtid))
	fields, err := runtime.ToRecordFields(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	rec := C.nar_record_t{
		size:   C.nar_size_t(len(fields)),
		keys:   (*C.nar_string_t)(nar_alloc(rtid, C.nar_size_t(len(fields)*C.sizeof_nar_string_t))),
		values: (*C.nar_object_t)(nar_alloc(rtid, C.nar_size_t(len(fields)*C.sizeof_nar_object_t))),
	}
	keysSlice := (*[1 << 30]C.nar_string_t)(unsafe.Pointer(rec.keys))
	valuesSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(rec.values))
	for i, f := range fields {
		keysSlice[i] = stringToC(f.K)
		valuesSlice[i] = C.nar_object_t(f.V)
	}
	return rec
}

//export nar_to_list
func nar_to_list(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_list_t {
	rt := runtime.GetById(runtime.Id(rtid))
	objects, err := runtime.ToList(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	list := C.nar_list_t{
		size:  C.nar_size_t(uint64(len(objects))),
		items: (*C.nar_object_t)(nar_alloc(rtid, C.nar_size_t(len(objects)*C.sizeof_nar_object_t))),
	}
	itemsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(list.items))
	for i, item := range objects {
		itemsSlice[i] = C.nar_object_t(item)
	}
	return list
}

//export nar_to_tuple
func nar_to_tuple(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_tuple_t {
	rt := runtime.GetById(runtime.Id(rtid))
	objects, err := runtime.ToTuple(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	tuple := C.nar_tuple_t{
		size:  C.nar_size_t(len(objects)),
		items: (*C.nar_object_t)(nar_alloc(rtid, C.nar_size_t(len(objects)*C.sizeof_nar_object_t))),
	}
	itemsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(tuple.items))
	for i, item := range objects {
		itemsSlice[i] = C.nar_object_t(item)
	}
	return tuple
}

//export nar_to_bool
func nar_to_bool(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_bool_t {
	rt := runtime.GetById(runtime.Id(rtid))
	v, err := runtime.ToBool(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	if v {
		return C.nar_bool_t(1)
	}
	return C.nar_bool_t(0)
}

//export nar_to_option
func nar_to_option(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_option_t {
	rt := runtime.GetById(runtime.Id(rtid))
	n, optValues, err := runtime.ToOption(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	optName, err := runtime.ToString(rt, n)
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	opt := C.nar_option_t{
		name:   stringToC(optName),
		size:   C.nar_size_t(len(optValues)),
		values: (*C.nar_object_t)(nar_alloc(rtid, C.nar_size_t(len(optValues)*C.sizeof_nar_object_t))),
	}
	valuesSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(opt.values))
	for i, item := range optValues {
		valuesSlice[i] = C.nar_object_t(item)
	}
	return opt
}

//export nar_to_func
func nar_to_func(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_func_t {
	rt := runtime.GetById(runtime.Id(rtid))
	f, a, err := runtime.ToFunc(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	return C.nar_func_t{ptr: C.nar_ptr_t(f), arity: C.nar_size_t(a)}
}

//export nar_to_native
func nar_to_native(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_native_t {
	rt := runtime.GetById(runtime.Id(rtid))
	ptr, size, err := runtime.ToNative(rt, runtime.Object(obj))
	if err != nil {
		setLastError(rt, runtime.TString(err.Error()))
	}
	return C.nar_native_t{
		ptr: C.nar_ptr_t(ptr),
		cmp: C.nar_cmp_native_fn_t(size),
	}
}

//export nar_alloc
func nar_alloc(rtid C.nar_runtime_t, size C.nar_size_t) unsafe.Pointer {
	rt := runtime.GetById(runtime.Id(rtid))
	ptr := C.malloc(C.size_t(size))
	rt.AppendFrameMemory(ptr)
	return ptr
}

//export nar_free_frame_memory
func nar_free_frame_memory(rtid C.nar_runtime_t) {
	rt := runtime.GetById(runtime.Id(rtid))
	rt.FreeFrameMemory(free)
}

//export nar_print
func nar_print(rtid C.nar_runtime_t, message C.nar_string_t) {
	fmt.Println(stringToGo(message))
}

//helpers --------------------------------------------

func free(ptr unsafe.Pointer) {
	C.free(ptr)
}

var lastError map[*runtime.Runtime]runtime.TString

func setLastError(rt *runtime.Runtime, err runtime.TString) {
	lastError[rt] = err
}

func getLastError(rt *runtime.Runtime) runtime.TString {
	err, ok := lastError[rt]
	if !ok {
		return ""
	}
	return err
}

func stringToGo(s C.nar_string_t) runtime.TString {
	strSlice := (*[1 << 30]C.char)(unsafe.Pointer(s))
	sz := 0
	for {
		if strSlice[sz] == 0 {
			break
		}
		sz += stringStride
	}
	str := C.GoStringN((*C.char)(unsafe.Pointer(s)), C.int(sz))
	res, _ := stringDecoder.String(str)
	return runtime.TString(res)
}

func stringToC(s runtime.TString) C.nar_string_t {
	encoded, _ := stringEncoder.String(string(s) + "\x00")
	cs := C.CBytes([]byte(encoded))
	return C.nar_string_t(unsafe.Pointer(cs))
}

var stringEncoder = func() *encoding.Encoder {
	if goruntime.GOOS == "windows" {
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	} else {
		return utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM).NewEncoder()
	}
}()

var stringDecoder = func() *encoding.Decoder {
	if goruntime.GOOS == "windows" {
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	} else {
		return utf32.UTF32(utf32.LittleEndian, utf32.IgnoreBOM).NewDecoder()
	}
}()

var stringStride = func() int {
	if goruntime.GOOS == "windows" {
		return 2
	} else {
		return 4
	}
}()
