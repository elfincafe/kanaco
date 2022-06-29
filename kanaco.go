package kanaco

/*
#cgo CFLAGS: -I./c
#cgo LDFLAGS: -L.build -lkanaco
#include <stdlib.h>
#include "kanaco.h"
*/
import "C"
import (
	"bufio"
	"fmt"
	"io"
	"unsafe"
)

type Reader struct {
	r    *bufio.Reader
	mode string
}

func Byte(b []byte, mode string) []byte {
	return []byte(String(string(b), mode))
}

func String(s, mod string) string {
	str := C.CString(s)
	str_len := C.int(len(s))
	mode := C.CString(mod)
	mode_len := C.int(len(mod))
	ret := C.convert(str, str_len, mode, mode_len)
	// defer C.freeMemory(unsafe.Pointer(ret))
	defer C.free(unsafe.Pointer(ret))
	return C.GoString(ret)
}

func NewReader(r io.Reader, mode string) *Reader {
	reader := new(Reader)
	reader.r = bufio.NewReader(r)
	reader.mode = mode
	return reader
}

func (r *Reader) Read(p []byte) (int, error) {
	line, err := r.r.ReadBytes('\n')
	if err == io.EOF {
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	line = Byte(line, r.mode)
	if len(p) < len(line) {
		return 0, fmt.Errorf("buffer size is not enough")
	}
	n := copy(p, line)
	return n, nil
}
