package kanaco

/*
#cgo CFLAGS: -I.
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

func Byte(bytes []byte, mode string) []byte {
	return []byte(String(string(bytes), mode))
}

func String(str, mode string) string {
	s := C.CString(str)
	slen := C.int(len(str))
	m := C.CString(mode)
	mlen := C.int(len(mode))
	ret := C.convert(s, slen, m, mlen)
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
