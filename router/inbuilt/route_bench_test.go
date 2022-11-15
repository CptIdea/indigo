package inbuilt

import (
	"github.com/fakefloordiv/indigo/http"
	"strings"
	"testing"

	methods "github.com/fakefloordiv/indigo/http/method"
)

func nopRender(_ http.Response) error {
	return nil
}

func BenchmarkRequestRouting(b *testing.B) {
	longURIRequest, _ := getRequest()
	longURIRequest.Method = methods.GET
	longURIRequest.Path = "/" + strings.Repeat("a", 255)

	shortURIRequest, _ := getRequest()
	shortURIRequest.Method = methods.GET
	shortURIRequest.Path = "/" + strings.Repeat("a", 15)

	unknownURIRequest, _ := getRequest()
	unknownURIRequest.Method = methods.GET
	unknownURIRequest.Path = "/" + strings.Repeat("b", 255)

	unknownMethodRequest, _ := getRequest()
	unknownMethodRequest.Method = methods.POST
	unknownMethodRequest.Path = longURIRequest.Path

	router := NewRouter()
	router.Get(longURIRequest.Path, nopHandler)
	router.Get(shortURIRequest.Path, nopHandler)

	router.OnStart()

	b.Run("LongURI", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			router.OnRequest(longURIRequest)
		}
	})

	b.Run("ShortURI", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			router.OnRequest(shortURIRequest)
		}
	})

	b.Run("UnknownURI", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			router.OnRequest(unknownURIRequest)
		}
	})

	b.Run("UnknownMethod", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			router.OnRequest(unknownMethodRequest)
		}
	})
}
