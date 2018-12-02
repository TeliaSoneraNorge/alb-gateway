package gateway

import (
	"context"
	"testing"

	"github.com/getas/alb-gateway/events"
	"github.com/tj/assert"
	"io/ioutil"
)

func TestNewRequest_path(t *testing.T) {
	e := events.LambdaTargetGroupRequest{
		Path: "/pets/luna",
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, "GET", r.Method)
	assert.Equal(t, `/pets/luna`, r.URL.Path)
	assert.Equal(t, `/pets/luna`, r.URL.String())
}

func TestNewRequest_method(t *testing.T) {
	e := events.LambdaTargetGroupRequest{
		HTTPMethod: "DELETE",
		Path:       "/pets/luna",
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, "DELETE", r.Method)
}

func TestNewRequest_queryString(t *testing.T) {
	e := events.LambdaTargetGroupRequest{
		HTTPMethod: "GET",
		Path:       "/pets",
		QueryStringParameters: map[string]string{
			"order":  "desc",
			"fields": "name,species",
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `/pets?fields=name%2Cspecies&order=desc`, r.URL.String())
	assert.Equal(t, `desc`, r.URL.Query().Get("order"))
}

func TestNewRequest_header(t *testing.T) {
	e := events.LambdaTargetGroupRequest{
		HTTPMethod: "POST",
		Path:       "/pets",
		Body:       `{ "name": "Tobi" }`,
		Headers: map[string]string{
			"Content-Type":      "application/json",
			"X-Foo":             "bar",
			"Host":              "example.com",
			"x-amzn-trace-id":   "Root=1-5bdb40ca-556d8b0c50dc66f0511bf520",
			"x-forwarded-for":   "72.21.198.66",
			"x-forwarded-port":  "443",
			"x-forwarded-proto": "https",
			"x-totally-random":  "randomheader",
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, `example.com`, r.Host)
	assert.Equal(t, `Root=1-5bdb40ca-556d8b0c50dc66f0511bf520`, r.Header.Get("x-amzn-trace-id"))
	assert.Equal(t, `72.21.198.66`, r.Header.Get("x-forwarded-for"))
	assert.Equal(t, `443`, r.Header.Get("x-forwarded-port"))
	assert.Equal(t, `https`, r.Header.Get("x-forwarded-proto"))
	assert.Equal(t, `randomheader`, r.Header.Get("x-totally-random"))
	assert.Equal(t, `18`, r.Header.Get("Content-Length"))
	assert.Equal(t, `application/json`, r.Header.Get("Content-Type"))
	assert.Equal(t, `bar`, r.Header.Get("X-Foo"))
}

func TestNewRequest_body(t *testing.T) {
	e := events.LambdaTargetGroupRequest{
		HTTPMethod: "POST",
		Path:       "/pets",
		Body:       `{ "name": "Tobi" }`,
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, `{ "name": "Tobi" }`, string(b))
}

func TestNewRequest_bodyBinary(t *testing.T) {
	e := events.LambdaTargetGroupRequest{
		HTTPMethod:      "POST",
		Path:            "/pets",
		Body:            `aGVsbG8gd29ybGQK`,
		IsBase64Encoded: true,
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, "hello world\n", string(b))
}

func TestNewRequest_remoteHost(t *testing.T) {
	e := events.LambdaTargetGroupRequest{
		HTTPMethod: "GET",
		Path:       "/pets",
		Headers: map[string]string{
			"x-amzn-trace-id": "Root=1-5bdb40ca-556d8b0c50dc66f0511bf520",
			"x-forwarded-for": "72.21.198.66",
		},
	}

	r, err := NewRequest(context.Background(), e)
	assert.NoError(t, err)

	assert.Equal(t, "72.21.198.66", r.RemoteAddr)
}

func TestNewRequest_context(t *testing.T) {
	e := events.LambdaTargetGroupRequest{}
	ctx := context.WithValue(context.Background(), "key", "value")
	r, err := NewRequest(ctx, e)
	assert.NoError(t, err)
	v := r.Context().Value("key")
	assert.Equal(t, "value", v)
}
