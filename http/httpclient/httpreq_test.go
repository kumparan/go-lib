package httpclient

import "testing"

func TestValidateClientOptions(t *testing.T) {
	cases := []struct {
		c      ClientOptions
		expect error
	}{
		{
			c: ClientOptions{
				BaseURL:    "http://something1:9001",
				HostHeader: "someaddress",
			},
			expect: ErrHostHeaderLimited,
		},
		{
			c: ClientOptions{
				HostHeader: "someaddress",
			},
			expect: ErrBaseURLEmpty,
		},
		{
			c: ClientOptions{
				BaseURL:    "127.0.0.1:9001",
				HostHeader: "someaddress",
			},
			expect: nil,
		},
		{
			c: ClientOptions{
				BaseURL:    "localhost:9001",
				HostHeader: "someaddress",
			},
			expect: nil,
		},
	}

	for _, c := range cases {
		err := c.c.Validate()
		if err != c.expect {
			t.Errorf("Expecting %s but got %s", c.expect.Error(), err.Error())
		}
	}
}
