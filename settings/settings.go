package settings

import "math"

const (
	defaultMaxHeaders           uint8  = math.MaxUint8
	defaultSockReadBuffSize     uint16 = 2048
	defaultMaxBodyLength        uint32 = math.MaxUint32
	defaultMaxURILength         uint16 = 4096
	defaultMaxHeaderKeyLength   uint8  = 100 // just like Apache
	defaultMaxHeaderValueLength uint16 = 8192
	defaultMaxBodyChunkLength   uint32 = math.MaxUint32
	defaultInfoLineBuffSize     uint16 = 30
	defaultHeadersBuffSize      uint16 = 500
)

type Settings struct {
	// MaxHeaders is a max number of headers allowed to keep, in case of exceeding this value
	// connection will be closed with StatusBadRequest code. By default, the value is 255,
	// and it cannot be more. IMHO nobody even needs more as 255 is already a hell
	MaxHeaders uint8

	// SockReadBuffSize is a size of buffer to which one we are reading from socket
	SockReadBuffSize uint16

	// MaxBodyLength is a maximal value accepted in Content-Length header
	MaxBodyLength uint32

	// MaxURILength is a maximal length of request path is accepted
	MaxURILength uint16

	// MaxHeaderKeyLength is a maximal length of header key is allowed (colon is not included)
	MaxHeaderKeyLength uint8

	// MaxHeaderValueLength is a maximal length of header value
	MaxHeaderValueLength uint16

	// MaxBodyChunkLength is a maximal length for body chunk (in case of chunked transfer encoding)
	MaxBodyChunkLength uint32

	// DefaultInfoLineBuffSize is a default capacity of newly allocated buffer for info line
	DefaultInfoLineBuffSize uint16

	// DefaultHeadersBuffSize is a default capacity of newly allocated buffer for headers line
	DefaultHeadersBuffSize uint16
}

func Prepare(settings Settings) Settings {
	if settings.MaxHeaders == 0 {
		settings.MaxHeaders = defaultMaxHeaders
	}
	if settings.SockReadBuffSize == 0 {
		settings.SockReadBuffSize = defaultSockReadBuffSize
	}
	if settings.MaxBodyLength == 0 {
		settings.MaxBodyLength = defaultMaxBodyLength
	}
	if settings.MaxURILength == 0 {
		settings.MaxURILength = defaultMaxURILength
	}
	if settings.MaxHeaderKeyLength == 0 {
		settings.MaxHeaderKeyLength = defaultMaxHeaderKeyLength
	}
	if settings.MaxHeaderValueLength == 0 {
		settings.MaxHeaderValueLength = defaultMaxHeaderValueLength
	}
	if settings.MaxBodyChunkLength == 0 {
		settings.MaxBodyChunkLength = defaultMaxBodyChunkLength
	}
	if settings.DefaultInfoLineBuffSize == 0 {
		settings.DefaultInfoLineBuffSize = defaultInfoLineBuffSize
	}
	if settings.DefaultHeadersBuffSize == 0 {
		settings.DefaultHeadersBuffSize = defaultHeadersBuffSize
	}

	return settings
}

func Default() Settings {
	return Prepare(Settings{})
}
