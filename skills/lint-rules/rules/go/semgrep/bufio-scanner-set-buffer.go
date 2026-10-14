package example

import (
	"bufio"
	"io"
)

func lines(r io.Reader) {
	// ruleid: bufio-scanner-set-buffer
	s := bufio.NewScanner(r)
	for s.Scan() {
	}
}

func linesOK(r io.Reader) {
	// ok: bufio-scanner-set-buffer
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0, 64*1024), 1<<20)
	for s.Scan() {
	}
}
