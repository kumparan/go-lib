package webserver

import (
	"encoding/json"
	"github.com/kumparan/go-lib/logger"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"

	"github.com/kumparan/go-lib/env"
	"github.com/kumparan/go-lib/router"
	"github.com/prometheus/client_golang/prometheus"
)

type Options struct {
	Address string
	Timeout time.Duration
}

type WebServer struct {
	router *router.Router
	port   string
	birth  time.Time
}

func checkOptions(opt *Options) {
	if opt.Timeout == time.Duration(0) {
		opt.Timeout = time.Second * 3
	}
	if len(opt.Address) > 0 {
		_, err := strconv.Atoi(opt.Address)
		// address is pure port number
		if err == nil {
			// check if first string is ":"
			if opt.Address[:1] != ":" {
				opt.Address = ":" + opt.Address
			}
		}
	} else {
		opt.Address = ":9000"
	}
}

func New(opt Options) *WebServer {
	checkOptions(&opt)
	address := opt.Address

	r := router.New(router.Options{Timeout: opt.Timeout})
	// provide metrics endpoint for prometheus metrics
	r.Handle("/metrics", prometheus.Handler())
	// provide service healthcheck
	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	// provide status of service when running
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		currentEnv := env.GetCurrentServiceEnv()
		configDir := env.GetConfigDir()
		buildNumber := env.GetCurrentBuild()
		logLevel := log.GetLevel()
		goVersion := env.GetGoVersion()

		response := map[string]interface{}{
			"environemnt": currentEnv,
			"config":      configDir,
			"log_level":   logLevel,
			"go_version":  goVersion,
			"build":       buildNumber,
		}
		jsonResp, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "application/json")
		w.Write(jsonResp)
	})
	web := WebServer{
		router: r,
		port:   address,
		birth:  time.Now(),
	}
	return &web
}

func (w *WebServer) Router() *router.Router {
	return w.router
}

func (w *WebServer) Run() error {
	logger.Infof("Webserver serving on: %s", w.port)
	return http.ListenAndServe(w.port, w.router)
}
