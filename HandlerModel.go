package uhttp

import (
	"context"
	"net/http"

	"github.com/dunv/uhttp/middlewares"
	"github.com/dunv/uhttp/params"
)

// Handler configured
type Handler struct {
	Pattern        string
	PostHandler    http.HandlerFunc
	PostModel      interface{}
	GetHandler     http.HandlerFunc
	GetModel       interface{}
	DeleteHandler  http.HandlerFunc
	DeleteModel    interface{}
	RequiredGet    params.R
	OptionalGet    params.R
	AddMiddlewares []Middleware
	AddMiddleware  *Middleware
	PreProcess     func(ctx context.Context) error
}

func (h Handler) WsReady() Middleware {
	chain := Chain(
		middlewares.ParseModel(h.PostModel, h.GetModel, h.DeleteModel),
		middlewares.GetParams(h.OptionalGet, h.RequiredGet),
	)

	// Add contexts
	for key, value := range requestContext {
		chain = Chain(chain, middlewares.WithContext(key, value))
	}

	// Add global middlewares
	for key := range additionalMiddlewares {
		chain = Chain(chain, additionalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range h.AddMiddlewares {
		chain = Chain(chain, h.AddMiddlewares[key])
	}

	// Add handler-specified middleware
	if h.AddMiddleware != nil {
		chain = Chain(chain, *h.AddMiddleware)
	}

	// Add preProcess
	return Chain(chain, middlewares.PreProcess(h.PreProcess))
}

func (h Handler) HandlerFunc() http.HandlerFunc {
	chain := Chain(
		middlewares.SetCors(config.CORS),
		middlewares.SetJSONResponse,
		middlewares.ParseModel(h.PostModel, h.GetModel, h.DeleteModel),
		middlewares.GetParams(h.OptionalGet, h.RequiredGet),
		middlewares.AddLogging,
	)

	// Add contexts
	for key, value := range requestContext {
		chain = Chain(chain, middlewares.WithContext(key, value))
	}

	// Add global middlewares
	for key := range additionalMiddlewares {
		chain = Chain(chain, additionalMiddlewares[key])
	}

	// Add handler-specified middlewares
	for key := range h.AddMiddlewares {
		chain = Chain(chain, h.AddMiddlewares[key])
	}

	// Add handler-specified middleware
	if h.AddMiddleware != nil {
		chain = Chain(chain, *h.AddMiddleware)
	}

	// Add preProcess
	chain = Chain(chain, middlewares.PreProcess(h.PreProcess))

	return SelectMethod(chain, h)
}
