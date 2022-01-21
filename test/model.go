package test

import (
	"io/fs"
	"time"
)

// MockFileInfo used for unit test
type MockFileInfo struct {
}

const mockFileInfoName = "mockFileInfoName"
const mockFileInfoSize = 0

// Name implement fs.FileInfo.Name()
func (MockFileInfo) Name() string {
	return mockFileInfoName
}

// Size implement fs.FileInfo.Size()
func (MockFileInfo) Size() int64 {
	return mockFileInfoSize
}

// Mode implement fs.FileInfo.Mode()
func (MockFileInfo) Mode() fs.FileMode {
	return fs.ModeDir
}

// ModTime implement fs.FileInfo.ModTime()
func (MockFileInfo) ModTime() time.Time {
	return time.Now()
}

// IsDir implement fs.FileInfo.IsDir()
func (MockFileInfo) IsDir() bool {
	return true
}

// Sys implement fs.FileInfo.Sys()
func (MockFileInfo) Sys() interface{} {
	return mockFileInfoName
}
