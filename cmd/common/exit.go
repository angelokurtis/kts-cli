package common

import (
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"log"
	"os"
)

func Exit(err error) {
	if err.Error() == "interrupt" {
		os.Exit(1)
	} else if _, ok := err.(stackTracer); ok {
		log.Fatal(color.Danger.Sprintf("%+v", err))
	} else {
		log.Fatal(color.Danger.Sprint(err.Error()))
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}
