/*
 * Copyright (c) 2019. Temple3x (temple3x@gmail.com)
 *
 * Use of this source code is governed by the MIT License
 * that can be found in the LICENSE file.
 */

package fnc

import "os"

// OpenFile opens a file with O_NOATIME flag.
// For some cases we don't need access time record.
func OpenFile(path string, flag int, perm os.FileMode) (f *os.File, err error) {

	flag |= O_NOATIME

	return os.OpenFile(path, flag, perm)
}

// Exist returns a file existed or not.
// Ignore error.
func Exist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// SyncDir syncs the given directory.
// Call it after rename, create new file etc if you want persist fs metadata.
// e.g. XFS uses delayed logging, may need SyncDir.
func SyncDir(dir string) (err error) {

	f, err := os.Open(dir)
	if err != nil {
		return
	}
	defer f.Close()

	return f.Sync()
}

const (
	syncFileRangeWaitBefore = 1
	syncFileRangeWrite      = 2
	syncFileRangeWaitAfter  = 4
)

// Flush flushes page_cache to disk in sync mode.
//
// OS may create a burst of write I/O when dirty pages hit a threshold,
// so flush it under users' control maybe a better choice in sometime.
func Flush(f *os.File, offset, size int64) (err error) {

	flags := syncFileRangeWaitBefore | syncFileRangeWrite | syncFileRangeWaitAfter
	return syncRange(f, offset, size, flags)
}

// FlushHint flushes in async mode.
//
// Warn: it can be stalled too in some situations.
func FlushHint(f *os.File, offset, size int64) (err error) {

	flags := syncFileRangeWrite
	return syncRange(f, offset, size, flags)
}

const (
	posixFadvRandom   = 1
	posixFadvDontneed = 4
)

// DropCache drops page_cache in range.
func DropCache(f *os.File, offset, size int64) (err error) {

	return fadvise(f, offset, size, posixFadvDontneed)
}

// DisableReadAhead disables file readahead entirely.
func DisableReadAhead(f *os.File) (err error) {

	return fadvise(f, 0, 0, posixFadvRandom)
}

// Preallocate allocates space for a new file.
// Avoid modify metadata & allocating space in future writing.
func PreAllocate(f *os.File, size int64) (err error) {
	return preAllocate(f, size)
}
