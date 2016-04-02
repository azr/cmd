package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

//code next was taken from
//https://github.com/golang/go/blob/fdba5a7544e54227c910ae3b26511c718df786a1/src/cmd/dist/util.go
//etc.
//
//it's to update when the system pkg implements CopyFile

const (
	writeExec = 1 << iota
	writeSkipSame
)

// ReadFile returns the content of the named file.
func ReadFile(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("ReadFile: %v", err)
	}
	return string(data)
}

// WriteFile writes b to the named file, creating it if needed.
// if exec is non-zero, marks the file as executable.
// If the file already exists and has the expected content,
// it is not rewritten, to avoid changing the time stamp.
func WriteFile(b, file string, flag int) {
	new := []byte(b)
	if flag&writeSkipSame != 0 {
		old, err := ioutil.ReadFile(file)
		if err == nil && bytes.Equal(old, new) {
			return
		}
	}
	mode := os.FileMode(0666)
	if flag&writeExec != 0 {
		mode = 0777
	}
	err := ioutil.WriteFile(file, new, mode)
	if err != nil {
		log.Fatalf("WriteFile: %v", err)
	}
}

// copy copies the file src to dst, via memory (so only good for small files).
func CopyFile(dst, src string, flag int) {
	WriteFile(ReadFile(src), dst, flag)
}
