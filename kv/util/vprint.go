package util

import (
	"fmt"
	"os"
	"strconv"
)

const (
	infoColor    = "\033[1;34m"
	noticeColor  = "\033[1;36m"
	warningColor = "\033[1;33m"
	errorColor   = "\033[1;31m"
	debugColor   = "\033[0;36m"
	endColor     = "\033[0m"
	endLine     = "\033[0m\n"
)

// We basically only care about two levels of logging, at least right now
func SetVprint() func(a ...interface{}) {
	verbose, _ := strconv.ParseInt(os.Getenv("VPRINTV"), 10, 16)
	//errlog := log.New(os.Stderr, "", 0)

	return func(a ...interface{}) {
		if verbose > 0 {
			fmt.Fprint(os.Stderr, warningColor)
			fmt.Fprint(os.Stderr, a...)
			fmt.Fprint(os.Stderr, endLine)

		}
	}
}

var Vprint = SetVprint()
