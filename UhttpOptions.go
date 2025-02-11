package uhttp

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dunv/ulog"
)

type UhttpOption interface {
	apply(*uhttpOptions)
}

type uhttpOptions struct {
	cors                    string
	log                     ulog.ULogger
	gzipCompressionLevel    int
	encodingErrorLogLevel   ulog.LogLevel
	parseModelErrorLogLevel ulog.LogLevel

	// Global middlewares
	globalMiddlewares []Middleware

	// If handler panics, send the information to the client
	sendPanicInfoToClient bool

	// Http-Server options
	address           string
	serveMux          *http.ServeMux
	readTimeout       time.Duration
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration

	// Static files encoding
	enableGzip   bool
	enableBrotli bool

	// TLS
	enableTLS      bool
	tlsErrorLogger *log.Logger
	tlsCertPath    *string
	tlsKeyPath     *string

	// Prometheus
	enableMetrics bool
	metricsSocket string
	metricsPath   string

	// Granular logging
	logHandlerCalls         bool
	logHandlerErrors        bool
	logHandlerRegistrations bool
}

type funcUhttpOption struct {
	f func(*uhttpOptions)
}

func (fdo *funcUhttpOption) apply(do *uhttpOptions) {
	fdo.f(do)
}

func newFuncUhttpOption(f func(*uhttpOptions)) *funcUhttpOption {
	return &funcUhttpOption{f: f}
}

func WithCORS(cors string) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.cors = cors
	})
}

func WithLogger(logger ulog.ULogger) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.log = logger
	})
}

func WithGzipCompressionLevel(level int) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.gzipCompressionLevel = level
	})
}

func WithEncodingErrorLogLevel(level ulog.LogLevel) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.encodingErrorLogLevel = level
	})
}

func WithParseModelErrorLogLevel(level ulog.LogLevel) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.parseModelErrorLogLevel = level
	})
}

func WithAddress(address string) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.address = address
	})
}

func WithReadTimeout(readTimeout time.Duration) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.readTimeout = readTimeout
	})
}

func WithReadHeaderTimeout(readHeaderTimeout time.Duration) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.readHeaderTimeout = readHeaderTimeout
	})
}

func WithWriteTimeout(writeTimeout time.Duration) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.writeTimeout = writeTimeout
	})
}

func WithIdleTimeout(idleTimeout time.Duration) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.idleTimeout = idleTimeout
	})
}

func WithGlobalMiddlewares(middlewares []Middleware) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.globalMiddlewares = middlewares
	})
}

func WithTLS(certPath string, keyPath string, tlsErrorLogger *log.Logger) UhttpOption {
	usedLogger := log.New(ioutil.Discard, "", log.Lshortfile)
	if tlsErrorLogger != nil {
		usedLogger = tlsErrorLogger
	}
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.enableTLS = true
		o.tlsErrorLogger = usedLogger
		o.tlsCertPath = &certPath
		o.tlsKeyPath = &keyPath
	})
}

func WithMetrics(metricsSocket string, metricsPath string) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.enableMetrics = true
		o.metricsSocket = metricsSocket
		o.metricsPath = metricsPath
	})
}

func WithSendPanicInfoToClient(sendPanicInfoToClient bool) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.sendPanicInfoToClient = sendPanicInfoToClient
	})
}

func WithStaticFileEncodings(enableGzip bool, enableBrotli bool) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.enableGzip = enableGzip
		o.enableBrotli = enableBrotli
	})
}

func WithGranularLogging(logHandlerCalls bool, logHandlerErrors bool, logHandlerRegistrations bool) UhttpOption {
	return newFuncUhttpOption(func(o *uhttpOptions) {
		o.logHandlerCalls = logHandlerCalls
		o.logHandlerErrors = logHandlerErrors
		o.logHandlerRegistrations = logHandlerRegistrations
	})
}
