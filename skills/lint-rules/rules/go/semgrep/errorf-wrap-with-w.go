package example

import "fmt"

var ErrUnavailable = fmt.Errorf("unavailable")

func load() error {
	err := fetch()
	if err != nil {
		// ok: errorf-wrap-with-w
		return fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	if err != nil {
		// ruleid: errorf-wrap-with-w
		return fmt.Errorf("load config: %v", err)
	}
	if err != nil {
		// ruleid: errorf-wrap-with-w
		return fmt.Errorf("rollback failed: %s", err)
	}
	// ok: errorf-wrap-with-w
	return fmt.Errorf("load config: %w", err)
}

func fetch() error { return nil }
