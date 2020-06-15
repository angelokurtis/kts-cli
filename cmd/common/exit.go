package common

import "os"

func Exit(err error) {
	if err.Error() == "interrupt" {
		os.Exit(1)
	} else {
		panic(err)
	}
}
