package ffi

/*
#cgo CFLAGS: -I${SRCDIR}/../../rust/addrgen/include
#cgo LDFLAGS: -L${SRCDIR}/../../rust/addrgen/target/release -ljuno_addrgen

#include "juno_addrgen.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

var errNull = errors.New("addrgen: null response")

func DeriveJSON(ufvk string, index uint32) (string, error) {
	cUFVK := C.CString(ufvk)
	defer C.free(unsafe.Pointer(cUFVK))

	out := C.juno_addrgen_derive_json(cUFVK, C.uint32_t(index))
	if out == nil {
		return "", errNull
	}
	defer C.juno_addrgen_string_free(out)

	return C.GoString(out), nil
}

func BatchJSON(ufvk string, start uint32, count uint32) (string, error) {
	cUFVK := C.CString(ufvk)
	defer C.free(unsafe.Pointer(cUFVK))

	out := C.juno_addrgen_batch_json(cUFVK, C.uint32_t(start), C.uint32_t(count))
	if out == nil {
		return "", errNull
	}
	defer C.juno_addrgen_string_free(out)

	return C.GoString(out), nil
}
