/*
 * Copyright (c) 2019. Temple3x (temple3x@gmail.com)
 *
 * Use of this source code is governed by the MIT License
 * that can be found in the LICENSE file.
 */

package fnc

import "os"

// Disable all File access time(atime) updates,
// darwin doesn't have it.
const O_NOATIME = 0

func syncRange(f *os.File, off int64, n int64, flags int) (err error) {

	// Before Go1.12 on darwin even call Sync, the drive may not write dirty page to the media.
	return f.Sync()
}

func fadvise(f *os.File, offset, size int64, advice int) (err error) {
	return
}

func preAllocate(f *os.File, size int64) error {
	return f.Truncate(size)
}
