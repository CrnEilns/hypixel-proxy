//go:build !windows && !wasm && !illumos

package buf

import (
	"net"
	"syscall"
	"unsafe"
)

type posixReader struct {
	iovecs []syscall.Iovec
}

func (r *posixReader) Init(bs net.Buffers) {
	iovecs := r.iovecs
	if iovecs == nil {
		iovecs = make([]syscall.Iovec, 0, len(bs))
	}
	for idx, b := range bs {
		iovecs = append(iovecs, syscall.Iovec{Base: &b[0]})
		iovecs[idx].SetLen(8 * 1024)
	}
	r.iovecs = iovecs
}

func (r *posixReader) Read(fd uintptr) int32 {
	n, _, e := syscall.Syscall(syscall.SYS_READV, fd, uintptr(unsafe.Pointer(&r.iovecs[0])), uintptr(len(r.iovecs)))
	if e != 0 {
		return -1
	}
	return int32(n)
}

func (r *posixReader) Clear() {
	for idx := range r.iovecs {
		r.iovecs[idx].Base = nil
	}
	r.iovecs = r.iovecs[:0]
}

func newVectorizedReader() vectorizedReader {
	return &posixReader{}
}
