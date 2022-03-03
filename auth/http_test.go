package auth_test

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gol4ng/fiberware/auth"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestFromHeader(t *testing.T) {
	tests := []struct {
		request            *fiber.Request
		expectedCredential string
	}{
		{
			request:            nil,
			expectedCredential: "",
		},
		{
			request: fiberRequest(map[string]string{
				"Authorization": "foo",
			}),
			expectedCredential: "foo",
		},
		{
			request: fiberRequest(map[string]string{
				"X-Authorization": "foo",
			}),
			expectedCredential: "foo",
		},
		{
			request: fiberRequest(map[string]string{
				"Authorization":   "foo",
				"X-Authorization": "bar",
			}),
			expectedCredential: "foo",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.Equal(t, auth.Credential(tt.expectedCredential), auth.FromHeader(tt.request)())
		})
	}
}

func TestAddHeader(t *testing.T) {
	t.Run("nil request", func(t *testing.T) {
		auth.AddHeader(nil)("foo")
	})

	t.Run("request", func(t *testing.T) {
		req := &fiber.Request{}

		credSetter := auth.AddHeader(req)
		credSetter("foo")
		assert.Equal(t, "foo", string(req.Header.Peek("Authorization")))
	})
}

func fiberRequest(headers map[string]string) *fiber.Request {
	header := fasthttp.RequestHeader{}
	for key, val := range headers {
		header.Set(key, val)
	}
	return &fiber.Request{Header: header}
}
