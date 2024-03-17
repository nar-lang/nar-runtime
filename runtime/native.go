package runtime

/*
#include <stdlib.h>
#include "../lib/nar.h"

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
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	goruntime "runtime"
	"unsafe"
)

//export xnar_fail
func xnar_fail(rtid C.nar_runtime_t, message C.nar_string_t) {
	msg := stringToGo(message)
	setLastError(GetById(Id(rtid)), msg)
	if KDebug {
		panic(msg)
	}
}

//export xnar_get_last_error
func xnar_get_last_error(rtid C.nar_runtime_t) C.nar_string_t {
	rt := GetById(Id(rtid))
	return stringToC(getLastError(rt))
}

//export xnar_load_bytecode
func xnar_load_bytecode(size C.nar_size_t, data *C.nar_byte_t) C.nar_bytecode_t {
	buf := (*[1 << 30]byte)(unsafe.Pointer(data))
	btc, err := bytecode.Read(bytes.NewReader(buf[:size]))
	if err != nil {
		setLastError(nil, "failed to load bytecode")
		return C.nar_bytecode_t(nil)
	}
	return C.nar_bytecode_t(btc)
}

//export xnar_create_runtime
func xnar_create_runtime(btc C.nar_bytecode_t, libsPath C.nar_string_t) C.nar_runtime_t {
	rt, err := NewRuntime((*bytecode.Binary)(btc), string(stringToGo(libsPath)))
	if err != nil {
		setLastError(nil, TString(err.Error()))
		return 0
	}
	return C.nar_runtime_t(rt.Id())
}

//export xnar_destroy_runtime
func xnar_destroy_runtime(rtid C.nar_runtime_t) {
	xnar_free_frame_memory(rtid)
	rt := GetById(Id(rtid))
	rt.Destroy()
}

//export xnar_register_def
func xnar_register_def(rtid C.nar_runtime_t, module_name C.nar_string_t, def_name C.nar_string_t, def C.nar_object_t) {
	rt := GetById(Id(rtid))
	rt.RegisterDef(bytecode.QualifiedIdentifier(stringToGo(module_name)), ast.Identifier(stringToGo(def_name)), Object(def))
}

//export xnar_apply
func xnar_apply(rtid C.nar_runtime_t, name C.nar_string_t, num_args C.nar_size_t, args *C.nar_object_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	args_ := make([]Object, 0, num_args)
	argsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(args))
	for i := C.nar_size_t(0); i < num_args; i++ {
		args_ = append(args_, Object(argsSlice[i]))
	}
	res, err := rt.Apply(bytecode.FullIdentifier(stringToGo(name)), args_...)
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	return C.nar_object_t(res)
}

//export xnar_apply_func
func xnar_apply_func(rtid C.nar_runtime_t, fn C.nar_object_t, num_args C.nar_size_t, args *C.nar_object_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	args_ := make([]Object, 0, num_args)
	argsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(args))
	for i := C.nar_size_t(0); i < num_args; i++ {
		args_ = append(args_, Object(argsSlice[i]))
	}
	res, err := rt.ApplyFunc(Object(fn), args_...)
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	return C.nar_object_t(res)
}

//export xnar_get_object_kind
func xnar_get_object_kind(_ C.nar_runtime_t, obj C.nar_object_t) C.nar_object_kind_t {
	return C.nar_object_kind_t(Object(obj).Kind())
}

//export xnar_new_unit
func xnar_new_unit(rtid C.nar_runtime_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	return C.nar_object_t(rt.NewUnit())
}

//export xnar_new_char
func xnar_new_char(rtid C.nar_runtime_t, value C.nar_char_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	return C.nar_object_t(rt.NewChar(TChar(value)))
}

//export xnar_new_int
func xnar_new_int(rtid C.nar_runtime_t, value C.nar_int_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	return C.nar_object_t(rt.NewInt(TInt(value)))
}

//export xnar_new_float
func xnar_new_float(rtid C.nar_runtime_t, value C.nar_float_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	return C.nar_object_t(rt.NewFloat(TFloat(value)))
}

//export xnar_new_string
func xnar_new_string(rtid C.nar_runtime_t, value C.nar_string_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	return C.nar_object_t(rt.NewString(stringToGo(value)))
}

//export xnar_new_record
func xnar_new_record(rtid C.nar_runtime_t, size C.nar_size_t, keys *C.nar_string_t, values *C.nar_object_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	fields_ := make([]TString, 0, size)
	keys_ := make([]Object, 0, size)
	keysSlice := (*[1 << 30]C.nar_string_t)(unsafe.Pointer(keys))
	valuesSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(values))
	for i := C.nar_size_t(0); i < size; i++ {
		fields_ = append(fields_, stringToGo(keysSlice[i]))
		keys_ = append(keys_, Object(valuesSlice[i]))
	}

	return C.nar_object_t(rt.NewRecord(fields_, keys_))
}

//export xnar_new_list
func xnar_new_list(rtid C.nar_runtime_t, size C.nar_size_t, items *C.nar_object_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	if size <= 0 {
		return C.nar_object_t(rt.NewList())
	}
	items_ := make([]Object, 0, size)
	itemsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(items))
	for i := C.nar_size_t(0); i < size; i++ {
		items_ = append(items_, Object(itemsSlice[i]))
	}
	return C.nar_object_t(rt.NewList(items_...))
}

//export xnar_new_list_cons
func xnar_new_list_cons(rtid C.nar_runtime_t, head C.nar_object_t, tail C.nar_object_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	return C.nar_object_t(rt.NewListItem(Object(head), Object(tail)))
}

//export xnar_new_tuple
func xnar_new_tuple(rtid C.nar_runtime_t, size C.nar_size_t, items *C.nar_object_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	items_ := make([]Object, 0, size)
	itemsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(items))
	for i := C.nar_size_t(0); i < size; i++ {
		items_ = append(items_, Object(itemsSlice[i]))
	}
	return C.nar_object_t(rt.NewTuple(items_...))
}

//export xnar_new_bool
func xnar_new_bool(rtid C.nar_runtime_t, value C.nar_bool_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	return C.nar_object_t(rt.NewBool(value != 0))
}

//export xnar_new_option
func xnar_new_option(rtid C.nar_runtime_t, name C.nar_string_t, size C.nar_size_t, items *C.nar_object_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	items_ := make([]Object, 0, size)
	itemsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(items))
	for i := C.nar_size_t(0); i < size; i++ {
		items_ = append(items_, Object(itemsSlice[i]))
	}
	optName := stringToGo(name)
	return C.nar_object_t(rt.NewOption(optName, items_...))
}

var funcWrappers = map[C.nar_ptr_t]unsafe.Pointer{}

//export xnar_new_func
func xnar_new_func(rtid C.nar_runtime_t, fn C.nar_ptr_t, arity C.nar_size_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	wrapper, ok := funcWrappers[fn]
	if !ok {
		switch int(arity) {
		case 0:
			{
				wfn := func() Object {
					return Object(C.call_func0(C.func0(fn), rtid))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 1:
			{
				wfn := func(a Object) Object {
					return Object(C.call_func1(C.func1(fn), rtid, C.nar_object_t(a)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 2:
			{
				wfn := func(a Object, b Object) Object {
					return Object(C.call_func2(C.func2(fn), rtid, C.nar_object_t(a), C.nar_object_t(b)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 3:
			{
				wfn := func(a Object, b Object, c Object) Object {
					return Object(C.call_func3(C.func3(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 4:
			{
				wfn := func(a Object, b Object, c Object, d Object) Object {
					return Object(C.call_func4(C.func4(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c), C.nar_object_t(d)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 5:
			{
				wfn := func(a Object, b Object, c Object, d Object, e Object) Object {
					return Object(C.call_func5(C.func5(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c), C.nar_object_t(d), C.nar_object_t(e)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 6:
			{
				wfn := func(a Object, b Object, c Object, d Object, e Object, f Object) Object {
					return Object(C.call_func6(C.func6(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c), C.nar_object_t(d), C.nar_object_t(e), C.nar_object_t(f)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 7:
			{
				wfn := func(a Object, b Object, c Object, d Object, e Object, f Object, g Object) Object {
					return Object(C.call_func7(C.func7(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c), C.nar_object_t(d), C.nar_object_t(e), C.nar_object_t(f), C.nar_object_t(g)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		case 8:
			{
				wfn := func(a Object, b Object, c Object, d Object, e Object, f Object, g Object, h Object) Object {
					return Object(C.call_func8(C.func8(fn), rtid, C.nar_object_t(a), C.nar_object_t(b), C.nar_object_t(c), C.nar_object_t(d), C.nar_object_t(e), C.nar_object_t(f), C.nar_object_t(g), C.nar_object_t(h)))
				}
				wrapper = unsafe.Pointer(&wfn)
			}
		}
		funcWrappers[fn] = wrapper
	}
	return C.nar_object_t(rt.NewFunc(wrapper, uint8(arity)))
}

//export xnar_new_native
func xnar_new_native(rtid C.nar_runtime_t, ptr C.nar_ptr_t, cmp C.nar_cmp_native_fn_t) C.nar_object_t {
	rt := GetById(Id(rtid))
	return C.nar_object_t(rt.NewNative(unsafe.Pointer(ptr), unsafe.Pointer(cmp)))
}

//export xnar_to_unit
func xnar_to_unit(rtid C.nar_runtime_t, obj C.nar_object_t) {
	rt := GetById(Id(rtid))
	_, err := ToUnit(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
}

//export xnar_to_char
func xnar_to_char(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_char_t {
	rt := GetById(Id(rtid))
	v, err := ToChar(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	return C.nar_char_t(v)
}

//export xnar_to_int
func xnar_to_int(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_int_t {
	rt := GetById(Id(rtid))
	v, err := ToInt(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	return C.nar_int_t(v)
}

//export xnar_to_float
func xnar_to_float(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_float_t {
	rt := GetById(Id(rtid))
	v, err := ToFloat(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	return C.nar_float_t(v)
}

//export xnar_to_string
func xnar_to_string(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_string_t {
	rt := GetById(Id(rtid))
	v, err := ToString(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	return stringToC(v)
}

//export xnar_to_record
func xnar_to_record(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_record_t {
	rt := GetById(Id(rtid))
	fields, err := ToRecordFields(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	rec := C.nar_record_t{
		size:   C.nar_size_t(len(fields)),
		keys:   (*C.nar_string_t)(xnar_alloc(rtid, C.nar_size_t(len(fields)*C.sizeof_nar_string_t))),
		values: (*C.nar_object_t)(xnar_alloc(rtid, C.nar_size_t(len(fields)*C.sizeof_nar_object_t))),
	}
	keysSlice := (*[1 << 30]C.nar_string_t)(unsafe.Pointer(rec.keys))
	valuesSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(rec.values))
	for i, f := range fields {
		keysSlice[i] = stringToC(f.K)
		valuesSlice[i] = C.nar_object_t(f.V)
	}
	return rec
}

//export xnar_to_list
func xnar_to_list(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_list_t {
	rt := GetById(Id(rtid))
	objects, err := ToList(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	list := C.nar_list_t{
		size:  C.nar_size_t(uint64(len(objects))),
		items: (*C.nar_object_t)(xnar_alloc(rtid, C.nar_size_t(len(objects)*C.sizeof_nar_object_t))),
	}
	itemsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(list.items))
	for i, item := range objects {
		itemsSlice[i] = C.nar_object_t(item)
	}
	return list
}

//export xnar_to_tuple
func xnar_to_tuple(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_tuple_t {
	rt := GetById(Id(rtid))
	objects, err := ToTuple(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	tuple := C.nar_tuple_t{
		size:  C.nar_size_t(len(objects)),
		items: (*C.nar_object_t)(xnar_alloc(rtid, C.nar_size_t(len(objects)*C.sizeof_nar_object_t))),
	}
	itemsSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(tuple.items))
	for i, item := range objects {
		itemsSlice[i] = C.nar_object_t(item)
	}
	return tuple
}

//export xnar_to_bool
func xnar_to_bool(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_bool_t {
	rt := GetById(Id(rtid))
	v, err := ToBool(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	if v {
		return C.nar_bool_t(1)
	}
	return C.nar_bool_t(0)
}

//export xnar_to_option
func xnar_to_option(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_option_t {
	rt := GetById(Id(rtid))
	n, optValues, err := ToOption(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	optName, err := ToString(rt, n)
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	opt := C.nar_option_t{
		name:   stringToC(optName),
		size:   C.nar_size_t(len(optValues)),
		values: (*C.nar_object_t)(xnar_alloc(rtid, C.nar_size_t(len(optValues)*C.sizeof_nar_object_t))),
	}
	valuesSlice := (*[1 << 30]C.nar_object_t)(unsafe.Pointer(opt.values))
	for i, item := range optValues {
		valuesSlice[i] = C.nar_object_t(item)
	}
	return opt
}

//export xnar_to_func
func xnar_to_func(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_func_t {
	rt := GetById(Id(rtid))
	f, a, err := ToFunc(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	return C.nar_func_t{ptr: C.nar_ptr_t(f), arity: C.nar_size_t(a)}
}

//export xnar_to_native
func xnar_to_native(rtid C.nar_runtime_t, obj C.nar_object_t) C.nar_native_t {
	rt := GetById(Id(rtid))
	ptr, size, err := ToNative(rt, Object(obj))
	if err != nil {
		setLastError(rt, TString(err.Error()))
	}
	return C.nar_native_t{
		ptr: C.nar_ptr_t(ptr),
		cmp: C.nar_cmp_native_fn_t(size),
	}
}

//export xnar_alloc
func xnar_alloc(rtid C.nar_runtime_t, size C.nar_size_t) unsafe.Pointer {
	rt := GetById(Id(rtid))
	ptr := C.malloc(C.size_t(size))
	rt.AppendFrameMemory(ptr)
	return ptr
}

//export xnar_free_frame_memory
func xnar_free_frame_memory(rtid C.nar_runtime_t) {
	rt := GetById(Id(rtid))
	rt.FreeFrameMemory(free)
}

//export xnar_print
func xnar_print(rtid C.nar_runtime_t, message C.nar_string_t) {
	fmt.Println(stringToGo(message))
}

//helpers --------------------------------------------

func free(ptr unsafe.Pointer) {
	C.free(ptr)
}

var lastError map[*Runtime]TString

func setLastError(rt *Runtime, err TString) {
	lastError[rt] = err
}

func getLastError(rt *Runtime) TString {
	err, ok := lastError[rt]
	if !ok {
		return ""
	}
	return err
}

func stringToGo(s C.nar_string_t) TString {
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
	return TString(res)
}

func stringToC(s TString) C.nar_string_t {
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
