package utils

import (
	"fmt"
	"os"
)

func MustRun(fn func() error) {
	if err := fn(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
