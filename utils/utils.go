package utils

import "log"

func Must[T any](x T, err error) T {
	if err != nil {
		log.Fatal(err)
	}

	return x
}
