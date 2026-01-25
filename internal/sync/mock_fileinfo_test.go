package sync

import "os"

type mockFileInfo struct {
	os.FileInfo
	isDir bool
}

func (m *mockFileInfo) IsDir() bool {
	return m.isDir
}
