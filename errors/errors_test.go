package errors

import (
	"errors"
	"reflect"
	"testing"
)

func TestErrs(t *testing.T) {
	t.Parallel()
	cases := []struct {
		err    *Errs
		expect *Errs
	}{
		{
			err: New("New error without anything"),
			expect: &Errs{
				err: errors.New("New error without anything"),
			},
		},
		{
			err: New("Error with fields", Fields{"first": "two", "satu": 2}),
			expect: &Errs{
				err:    errors.New("Error with fields"),
				fields: Fields{"first": "two", "satu": 2},
			},
		},
	}

	for _, val := range cases {
		if !reflect.DeepEqual(val.err, val.expect) {
			t.Errorf("Expect %+v but got %+v", val.err, val.expect)
		}
	}
}

func TestMessages(t *testing.T) {
	t.Parallel()
	err := New("Some error", []string{"stack1", "stack2"})
	if len(err.GetMessages()) != 2 {
		t.Errorf("Expect %d but got %d", 2, len(err.GetMessages()))
	}
	err = New(err, []string{"field1", "field2"})
	if len(err.GetMessages()) != 4 {
		t.Errorf("Expect %d but got %d after append", 4, len(err.GetMessages()))
	}
}

func TestMatch(t *testing.T) {
	t.Parallel()
	cases := []struct {
		err1        error
		err2        error
		expectMatch bool
	}{
		{
			err1:        New(errors.New("This is new error")),
			err2:        nil,
			expectMatch: false,
		},
		{
			err1:        New(errors.New("This is new error")),
			err2:        errors.New("This is new error"),
			expectMatch: true,
		},
		{
			err1:        New(errors.New("This is new error")),
			err2:        errors.New("Something is different"),
			expectMatch: false,
		},
	}

	for _, val := range cases {
		if match := Match(val.err1, val.err2); match != val.expectMatch {
			t.Errorf("TestMatch: Expecting %v but got %v", val.expectMatch, match)
		}
	}
}
