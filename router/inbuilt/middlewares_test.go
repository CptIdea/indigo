package inbuilt

import (
	"github.com/indigo-web/indigo/http/decode"
	"github.com/indigo-web/indigo/http/method"
	"github.com/indigo-web/indigo/http/status"
	"github.com/indigo-web/indigo/internal/server/tcp/dummy"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/indigo-web/indigo/http"
	"github.com/indigo-web/indigo/internal/parser/http1"
	"github.com/indigo-web/indigo/settings"

	"github.com/indigo-web/indigo/router/inbuilt/types"

	"github.com/indigo-web/indigo/http/headers"
	"github.com/indigo-web/indigo/http/query"
)

/*
This file is separated because it is a bit specific and contains a lot
of specific stuff for testing only middlewares. Decided it's better to
separate it from all the other tests
*/

type middleware uint8

const (
	m1 middleware = iota + 1
	m2
	m3
	m4
	m5
	m6
	m7
)

type callstack struct {
	chain []middleware
}

func (c *callstack) Push(ware middleware) {
	c.chain = append(c.chain, ware)
}

func (c *callstack) Chain() []middleware {
	return c.chain
}

func (c *callstack) Clear() {
	c.chain = c.chain[:0]
}

func getMiddleware(mware middleware, stack *callstack) types.Middleware {
	switch mware {
	case m1:
		return func(next types.Handler, request *http.Request) http.Response {
			stack.Push(m1)
			return next(request)
		}
	case m2:
		return func(next types.Handler, request *http.Request) http.Response {
			stack.Push(m2)
			return next(request)
		}
	case m3:
		return func(next types.Handler, request *http.Request) http.Response {
			stack.Push(m3)
			return next(request)
		}
	case m4:
		return func(next types.Handler, request *http.Request) http.Response {
			stack.Push(m4)
			return next(request)
		}
	case m5:
		return func(next types.Handler, request *http.Request) http.Response {
			stack.Push(m5)
			return next(request)
		}
	case m6:
		return func(next types.Handler, request *http.Request) http.Response {
			stack.Push(m6)
			return next(request)
		}
	case m7:
		return func(next types.Handler, request *http.Request) http.Response {
			stack.Push(m7)
			return next(request)
		}
	}

	return func(next types.Handler, request *http.Request) http.Response {
		stack.Push(mware)

		return next(request)
	}
}

func getRequest() *http.Request {
	q := query.NewQuery(nil)
	bodyReader := http1.NewBodyReader(dummy.NewNopClient(), settings.Default().Body)

	return http.NewRequest(
		headers.NewHeaders(nil), q, http.NewResponse(), dummy.NewNopConn(),
		http.NewBody(bodyReader, decode.NewDecoder()), nil, false,
	)
}

func TestMiddlewares(t *testing.T) {
	stack := new(callstack)
	r := New()
	r.Use(getMiddleware(m1, stack))
	r.Get("/", http.Respond, getMiddleware(m2, stack))

	api := r.Group("/api")
	api.Use(getMiddleware(m3, stack))

	v1 := api.Group("/v1")
	v1.Get("/hello", http.Respond, getMiddleware(m6, stack))
	v1.Use(getMiddleware(m4, stack))

	v2 := api.Group("/v2")
	v2.Get("/world", http.Respond, getMiddleware(m7, stack))
	v2.Use(getMiddleware(m5, stack))

	require.NoError(t, r.OnStart())

	t.Run("/", func(t *testing.T) {
		request := getRequest()
		request.Method = method.GET
		request.Path.String = "/"

		response := r.OnRequest(request)
		require.Equal(t, status.OK, response.Code)
		require.Equal(t, []middleware{m1, m2}, stack.Chain())
		stack.Clear()
	})

	t.Run("/api/v1/hello", func(t *testing.T) {
		request := getRequest()
		request.Method = method.GET
		request.Path.String = "/api/v1/hello"

		response := r.OnRequest(request)
		require.Equal(t, status.OK, response.Code)
		require.Equal(t, []middleware{m1, m3, m4, m6}, stack.Chain())
		stack.Clear()
	})

	t.Run("/api/v2/world", func(t *testing.T) {
		request := getRequest()
		request.Method = method.GET
		request.Path.String = "/api/v2/world"

		response := r.OnRequest(request)
		require.Equal(t, status.OK, response.Code)
		require.Equal(t, []middleware{m1, m3, m5, m7}, stack.Chain())
		stack.Clear()
	})
}
