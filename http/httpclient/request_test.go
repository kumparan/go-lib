package httpclient

import (
	"strings"
	"testing"
)

func TestParseURL(t *testing.T) {
	cases := []struct {
		url            string
		params         []string
		expectContains []string
	}{
		{
			url:            "test:900",
			params:         []string{"test1", "val1", "test2", "val2"},
			expectContains: []string{"test1=val1", "test2=val2"},
		},
		{
			url:            "http://something1:900",
			params:         []string{"test1", "val1", "test2", "val2"},
			expectContains: []string{"test1=val1", "test2=val2"},
		},
		{
			url:            "https://something2:900",
			params:         []string{"test1", "val1", "test2", "val2", "test3", "a very long text"},
			expectContains: []string{"test1=val1", "test2=val2", "test3=a+very+long+text"},
		},
	}

	for _, c := range cases {
		u, err := ParseURL(c.url, c.params...)
		if err != nil {
			t.Error(err)
		}
		for _, expect := range c.expectContains {
			if !strings.Contains(u, expect) {
				t.Errorf("Expect %s but url is %s", expect, u)
			}
		}
	}
}
