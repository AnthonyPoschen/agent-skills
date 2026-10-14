package example

import (
	"errors"
	"strings"
)

var ErrNotFound = errors.New("not found")

func classify(err error) bool {
	// ruleid: no-error-string-match
	if strings.Contains(err.Error(), "not found") {
		return true
	}
	// ruleid: no-error-string-match
	if err.Error() == "no such table" {
		return true
	}
	// ok: no-error-string-match
	return errors.Is(err, ErrNotFound)
}
