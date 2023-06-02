//Package middleware
package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	//"time"
	"rainier/internal/config"

	"github.com/rs/zerolog/log"
)

// HTTPHandlerContext provides context for passing global values to handlers
// such as http thread pools or database handlers
//
// SEE: https://drstearns.github.io/tutorials/gohandlerctx/
type MiddleWareContext struct {
	hConfig *config.AppConfig
}

//CTXHandlerContext constructs a new HandlerContext,
//ensuring that the dependencies are valid values
func CTXMiddlewareContext(config *config.AppConfig) *MiddleWareContext {
	return &MiddleWareContext{config}
}

// WPLog is a middleware that wraps the http.Handler and it records
// how long the handler took to run, which path was called, and the status code.
// This method is going to be used with gorilla/mux.
func (mw *MiddleWareContext) WPLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()

		delegate := &responseWriterDelegator{ResponseWriter: w}
		rw := delegate

		next.ServeHTTP(rw, r) // call original

		code := sanitizeCode(delegate.status)
		method := sanitizeMethod(r.Method)

		// Throw into a go routine so it does not block, but probably is alreayd in a go routine...have to check
		go log.Info().Str("uri", r.RequestURI).Str("method", method).Str("response_time", time.Since(begin).String()).Str("status", code).Msg("")
	})
}

func (mw *MiddleWareContext) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	begin := time.Now()

	delegate := &responseWriterDelegator{ResponseWriter: w}
	rw := delegate

	next(rw, r) // call original

	code := sanitizeCode(delegate.status)
	method := sanitizeMethod(r.Method)

	go log.Info().Str("uri", r.RequestURI).Str("method", method).Str("response_time", time.Since(begin).String()).Str("status", code).Msg("")
}

type responseWriterDelegator struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (r *responseWriterDelegator) WriteHeader(code int) {
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}
func (r *responseWriterDelegator) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}

func sanitizeMethod(m string) string {
	return strings.ToLower(m)
}

func sanitizeCode(s int) string {
	return strconv.Itoa(s)
}
