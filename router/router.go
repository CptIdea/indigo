package router

import (
	"github.com/fakefloordiv/indigo/http/encodings"
	"github.com/fakefloordiv/indigo/types"
)

// Router is a general interface for any router compatible with indigo
// OnRequest called every time headers are parsed and ready to be processed
// OnError called once, and if it called, it means that connection will be
//         closed anyway. So you can process the error, send some response,
//         and when you are ready, just notify core that he can safely close
//         the connection (even if it's already closed from client side)
type Router interface {
	OnRequest(request *types.Request, writer types.ResponseWriter) error
	OnError(request *types.Request, writer types.ResponseWriter, err error)
}

// OnStart called when server is initialized and started. Can be implemented
// optionally
type OnStart interface {
	OnStart()
}

// GetContentEncodings returns a struct with custom-set content encodings.
// Can be implemented optionally
type GetContentEncodings interface {
	GetContentEncodings() encodings.ContentEncodings
}
