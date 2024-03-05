package historylog

import "os"

type HFile struct {
	*os.File
	history []byte
}

func OpenFile(name string, flag int, perm os.FileMode) (*HFile, error) {
	f, err := os.OpenFile(name, flag, perm)
	return &HFile{
		File:    f,
		history: make([]byte, 0),
	}, err
}

func (hFile *HFile) Write(p []byte) (n int, err error) {
	hFile.history = append(hFile.history, p...)
	return hFile.File.Write(p)
}

func (hFile *HFile) TakeOutHistory() string {
	history := hFile.history
	hFile.history = hFile.history[:0]
	return string(history)
}
