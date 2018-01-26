package router

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lab46/example/pkg/log"
	"github.com/prometheus/client_golang/prometheus"
)

//Example of router package, the router wrapper is based on gorilla/mux

// Router of vehicle-insurance
type Router struct {
	opt Options
	r   *mux.Router
}

// Options for router
type Options struct {
	Timeout time.Duration
}

// New router
func New(opt Options) *Router {
	muxRouter := mux.NewRouter()
	rtr := &Router{
		r:   muxRouter,
		opt: opt,
	}
	return rtr
}

// URLParam get param from rest request
func URLParam(r *http.Request, key string) string {
	params := mux.Vars(r)
	return params[key]
}

// SubRouter return a new Router with path prefix
func (rtr Router) SubRouter(pathPrefix string) *Router {
	muxSubrouter := rtr.r.PathPrefix(pathPrefix).Subrouter()
	return &Router{
		r:   muxSubrouter,
		opt: rtr.opt,
	}
}

// timeout middleware
// the timeout middleware should cover timeout budget
func (rtr *Router) timeout(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opt := rtr.opt
		// cancel context
		if opt.Timeout > 0 {
			ctx, cancel := context.WithTimeout(r.Context(), opt.Timeout*time.Second)
			defer cancel()
			r = r.WithContext(ctx)
		}

		doneChan := make(chan bool)
		go func() {
			h(w, r)
			doneChan <- true
		}()
		select {
		case <-r.Context().Done():
			// only an example response
			resp := map[string]interface{}{
				"errors": []string{"Request timed out"},
			}
			jsonResp, _ := json.Marshal(resp)
			w.WriteHeader(http.StatusRequestTimeout)
			// only an example response
			w.Write(jsonResp)
			return
		case <-doneChan:
			return
		}
	}
}

// responseWriterDelegator to delegate the current writer
// this is a 100% from prometheus delegator with some modification
// the modification is needed because namespace is required
type responseWriterDelegator struct {
	http.ResponseWriter
	// handler, method string
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

func sanitizeStatusCode(status int) string {
	code := strconv.Itoa(status)
	return code
}

// Get function
func (rtr Router) Get(pattern string, h http.HandlerFunc) {
	log.Debugf("[router][get] %s", pattern)
	rtr.r.HandleFunc(pattern, prometheus.InstrumentHandlerFunc(pattern, rtr.timeout(h))).Methods("GET")
}

// Post function
func (rtr Router) Post(pattern string, h http.HandlerFunc) {
	log.Debugf("[router][post] %s", pattern)
	rtr.r.HandleFunc(pattern, prometheus.InstrumentHandlerFunc(pattern, rtr.timeout(h))).Methods("POST")
}

// Put function
func (rtr Router) Put(pattern string, h http.HandlerFunc) {
	log.Debugf("[router][put] %s", pattern)
	rtr.r.HandleFunc(pattern, prometheus.InstrumentHandlerFunc(pattern, rtr.timeout(h))).Methods("PUT")
}

// Delete function
func (rtr Router) Delete(pattern string, h http.HandlerFunc) {
	log.Debugf("[router][delete] %s", pattern)
	rtr.r.HandleFunc(pattern, prometheus.InstrumentHandlerFunc(pattern, rtr.timeout(h))).Methods("DELETE")
}

// Patch function
func (rtr Router) Patch(pattern string, h http.HandlerFunc) {
	log.Debugf("[router][patch] %s", pattern)
	rtr.r.HandleFunc(pattern, prometheus.InstrumentHandlerFunc(pattern, rtr.timeout(h))).Methods("PATCH")
}

// Handle function
func (rtr Router) Handle(pattern string, h http.Handler) {
	log.Debugf("[router][handle] %s", pattern)
	rtr.r.Handle(pattern, h)
}

// ServeHTTP function
func (rtr Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr.r.ServeHTTP(w, r)
}
