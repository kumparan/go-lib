package errors

// errors package inspired and a subset copy of errors package in upspin project

import (
	"errors"
	"runtime"

	"log"
)

var runtimeOutput bool

// SetRuntimeOutput will provide error information where the error is happened
func SetRuntimeOutput(b bool) {
	runtimeOutput = b
}

// IsRuntimeEnabled to check the status of runtimeOutput
func IsRuntimeOuputEnabled() bool {
	return runtimeOutput
}

type Fields map[string]interface{}

func (f Fields) ToArrayInterface() []interface{} {
	fieldsLength := len(f)
	if fieldsLength == 0 {
		return nil
	}
	// because fields is key value, we need to create an array with len * 2
	intf := make([]interface{}, fieldsLength*2)
	counter := 0
	for key, val := range f {
		intf[counter] = key
		intf[counter+1] = val
		counter += 2
	}
	return intf
}

// Errs struct
type Errs struct {
	err error
	// Codes used for Errs to identify known errors in the application
	// If the error is expected by Errs object, the errors will be shown as listed in Codes
	code    Codes
	message string

	// Traces used to add function traces to errors, this is different from context
	// While context is used to add more information about the error, traces is used
	// for easier function tracing purposes without hurting heap too much
	traces []string

	// Fields is a fields context similar to logrus.Fields
	// Can be used for adding more context to the errors
	fields Fields

	// Messages is a field to add stack of messages to error
	// this is used to simplify error message stack
	messages []string

	// var for runtime output
	file string
	line int
}

var _ error = (*Errs)(nil)

// New Errs
func New(args ...interface{}) *Errs {
	var (
		er    error
		isBad bool
	)
	err := &Errs{err: errors.New("Unknown error")}
	for _, arg := range args {
		switch arg.(type) {
		case string:
			er = errors.New(arg.(string))
		case *Errs:
			// copy and put the errors back
			errcpy := *arg.(*Errs)
			err = &errcpy
		// error should be placed below *Errs
		// implementation of Error() string will detect *Errs as error
		case error:
			er = arg.(error)
		case Codes:
			err.code = arg.(Codes)
			errString, _ := err.code.ErrorAndCode()
			er = errors.New(errString)
		// Fields cannot be appended
		// new fields will always replace the old fields
		case Fields:
			err.fields = arg.(Fields)
		// []string is detected as Errs.Messages
		// Messages can be appended, but might need to create a different type in the future
		case []string:
			if err.messages == nil {
				err.messages = make([]string, 0)
			}
			msgs := arg.([]string)
			err.messages = append(err.messages, msgs...)
		default:
			// the default error is unknown
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.Errs: bad call from %s:%d: %v", file, line, args)
		}
	}
	// if er have value then set errrors.error to er
	if er != nil {
		err.err = er
	}
	// only get the runtime file and line if err is defined
	if runtimeOutput && !isBad {
		_, err.file, err.line, _ = runtime.Caller(1)
	}
	return err
}

// WithCodes give a safer passing of codes to errors as compiler/linter will check the interface{} implementation
func WithCodes(codes Codes) *Errs {
	return New(codes)
}

func (e *Errs) Error() string {
	return e.err.Error()
}

// SetMessage for error
func (e *Errs) SetMessage(message string) {
	e.message = message
}

// GetMessage return message for error
func (e *Errs) GetMessage() string {
	return e.message
}

// GetTrace return traces
func (e *Errs) GetTrace() []string {
	return e.traces
}

// GetFields return available fields in errors
func (e *Errs) GetFields() Fields {
	return e.fields
}

// GetMessages return array of errors, this is depends by what kind of messages can be exists in the stack.
func (e *Errs) GetMessages() []string {
	return e.messages
}

// GetFileAndLine is part of runtimeOutput, as runtime will give file and line information
// will give empty string and 0 if runtimeOutput is false
func (e *Errs) GetFileAndLine() (string, int) {
	return e.file, e.line
}

/*
Match will match two strings error through a fuzzy matching
Need some improvement in fuzzy matching, not all cases is covered
*/

// Match error
func Match(errs1, errs2 error) bool {
	if errs1 == nil && errs2 == nil {
		return true
	}

	if errs1 != nil {
		err1, ok := errs1.(*Errs)
		if ok {
			errs1 = err1.err
		}
	} else {
		errs1 = errors.New("nil")
	}

	if errs2 != nil {
		err2, ok := errs2.(*Errs)
		if ok {
			errs2 = err2.err
		}
	} else {
		errs2 = errors.New("nil")
	}

	if errs1.Error() != errs2.Error() {
		return false
	}
	return true
}

// Codes is interface to define error custom code.
// It have two function called ErrorAndCode which return string of error and httpcode desired from the error
// Err will return the error of code itself, so error can be implemented directly in Codes
type Codes interface {
	ErrorAndCode() (string, int)
	Err() error
}
