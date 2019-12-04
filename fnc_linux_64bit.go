// +build arm64 amd64 ppc64 ppc64le mips64 mips64le s390x
// +build linux

/*
 * Copyright (c) 2019. Temple3x (temple3x@gmail.com)
 *
 * Use of this source code is governed by the MIT License
 * that can be found in the LICENSE file.
 */

package fnc

import (
	"os"
	"syscall"
)

// Disable all File access time(atime) updates.
const O_NOATIME = syscall.O_NOATIME

func syncRange(f *os.File, offset int64, size int64, flags int) (err error) {
	return syscall.SyncFileRange(int(f.Fd()), offset, size, flags)
}

func fadvise(f *os.File, offset, size int64, advice int) (err error) {

	// discard partial pages are ignored.
	var align int64
	align = 1 << 12
	size = (size + align - 1) &^ (align - 1)

	_, _, errno := syscall.Syscall6(syscall.SYS_FADVISE64, f.Fd(), uintptr(offset), uintptr(size), uintptr(advice), 0, 0)
	if errno != 0 {
		err = errno
	}
	return
}

const fallocate_default = 0 // file size will change if off+len is greater than fsize

func preAllocate(f *os.File, size int64) error {

	err := syscall.Fallocate(int(f.Fd()), fallocate_default, 0, size)
	if err != nil {
		errno, ok := err.(syscall.Errno)
		if ok &&
			// Not support is rare, in case bad news here.
			(errno == syscall.ENOTSUP ||
				// Go does retry syscall(signal handler set SA_RESTART),
				// but some platform(caused by bugs) may still return EINTR.
				errno == syscall.EINTR) {
			return f.Truncate(size)
		}
	}
	return err
}
