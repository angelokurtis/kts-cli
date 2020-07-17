package common

import (
	"github.com/pkg/errors"
	"log"
	"os"
)

func Exit(err error) {
	if err.Error() == "interrupt" {
		os.Exit(1)
	} else if _, ok := err.(stackTracer); ok {
		log.Fatalf("%+v", err)
	} else {
		log.Panic(err.Error())
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
