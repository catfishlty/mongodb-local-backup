package test

import (
    "io/fs"
    "time"
)

type MockFileInfo struct {
}

const mockFileInfoName = "mockFileInfoName"
const mockFileInfoSize = 0

func (MockFileInfo) Name() string {
    return mockFileInfoName
}

func (MockFileInfo) Size() int64 {
    return mockFileInfoSize
}

func (MockFileInfo) Mode() fs.FileMode {
    return fs.ModeDir
}

func (MockFileInfo) ModTime() time.Time {
    return time.Now()
}

func (MockFileInfo) IsDir() bool {
    return true
}

func (MockFileInfo) Sys() interface{} {
    return mockFileInfoName
}
