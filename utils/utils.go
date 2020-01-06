package utils

import (
	"fmt"
	"os"
)

func MustRun(fn func() error) {
	if err := fn(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
