package http

import (
	"context"
	"fmt"
	"general-agent/config"
	"general-agent/extension/logz"
	httpConfig "general-agent/framework/http/config"
	"net/http"

	"github.com/go-errors/errors"
)

// HTTPServer http server .
type HTTPServer struct {
	server  *http.Server
	options *httpConfig.Options
}

func NewEngineServer(engine *Engine, conf *config.Config) *HTTPServer {
	return NewHTTPServer(engine, &conf.HTTP)
}

// NewHTTPServer Deprecated new http server .
func NewHTTPServer(handler http.Handler, opt *httpConfig.Options) *HTTPServer {
	return &HTTPServer{
		options: opt,
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%d", opt.Host, opt.Port),
			Handler:      handler,
			ReadTimeout:  opt.ReadTimeout,
			WriteTimeout: opt.WriteTimeout,
		},
	}
}

// Serve start server
func (h *HTTPServer) Serve(ctx context.Context) {
	logz.Info(ctx, fmt.Sprintf("Start HTTP Server on %s, Access URL http://localhost:%d", h.server.Addr, h.options.Port))
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logz.Error(ctx, "Start HTTP Server Failed...", logz.Err(err))
		panic(err)
	}
}

func (h *HTTPServer) ServeTLS(ctx context.Context, certFile, keyFile string) {
	logz.Info(ctx, fmt.Sprintf("Start HTTPS Server on %s, Access URL https://localhost:%d", h.server.Addr, h.options.Port))
	if err := h.server.ListenAndServeTLS(certFile, keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logz.Error(ctx, "Start HTTPS Server Failed...", logz.Err(err))
		panic(err)
	}
}

// Shutdown end server
func (h *HTTPServer) Shutdown(ctx context.Context) error {
	if err := h.server.Shutdown(ctx); err != nil {
		logz.Error(ctx, "Stop HTTP Server Failed...", logz.Err(err))
		return err
	}
	logz.Info(ctx, "Stop HTTP Server Done")
	return nil
}
