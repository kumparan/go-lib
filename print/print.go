package print

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kumparan/go-lib/errors"
)

// this package is used to produce pretty print output
// mostly used by ssi-cli app
// all print should go to stdout instead of stderr

// var for printing prefix, this is naive and need a better implementation
var (
	prefixInfo  = func(prefix string) string { return color.GreenString("[INFO]" + prefix) }
	prefixDebug = func(prefix string) string { return color.YellowString("[DEBUG]" + prefix) }
	prefixWarn  = func(prefix string) string { return color.HiCyanString("[WARN]" + prefix) }
	prefixError = func(prefix string) string { return color.RedString("[ERROR]" + prefix) }
)

// debug var to identify if debug print is allowed or not
var isDebug bool

func SetDebug(debug bool) {
	isDebug = debug
}

func Debug(v ...interface{}) {
	// idiomatic debug
	if !isDebug {
		return
	}
	print(prefixDebug(""), v...)
}

func Info(v ...interface{}) {
	print(prefixInfo(""), v...)
}

func Warn(v ...interface{}) {
	print(prefixWarn(""), v...)
}

func Error(v ...interface{}) {
	print(prefixError(""), v...)
}

func Fatal(err error) {
	if err == nil {
		return
	}
	Error(err)
	os.Exit(1)
}

func print(prefix string, v ...interface{}) {
	// naively reject if only tag
	if len(v) == 0 {
		return
	}
	// return if parased argument is not valid
	parsedIntf := parseArgs(v...)
	if len(parsedIntf) == 0 {
		return
	}
	newIntf := []interface{}{prefix}
	newIntf = append(newIntf, parsedIntf...)
	fmt.Println(newIntf...)
}

// TODO: count on interface{} length and discard append to reduce memory use
func parseArgs(v ...interface{}) []interface{} {
	var newIntf []interface{}
	for key, val := range v {
		switch val.(type) {
		// dispatch if array of string
		case []string:
			arrOfString := val.([]string)
			for _, stringval := range arrOfString {
				newIntf = append(newIntf, stringval)
			}
			continue
		case nil:
			continue
		// pretty print errors if available
		case *errors.Errs:
			err := val.(*errors.Errs)
			newIntf = append(newIntf, err.Error())
			fields := err.GetFields()
			for key, val := range fields {
				s := fmt.Sprintf("%v=%v", key, val)
				newIntf = append(newIntf, s)
			}
			file, line := err.GetFileAndLine()
			if line != 0 {
				newIntf = append(newIntf, fmt.Sprintf("err_file=%s", file))
				newIntf = append(newIntf, fmt.Sprintf("err_line=%d", line))
			}
			continue
		}
		newIntf = append(newIntf, v[key])
	}
	return newIntf
}
