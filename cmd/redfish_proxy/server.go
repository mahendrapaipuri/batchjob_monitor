package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httputil"
	_ "net/http/pprof" // #nosec
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/exporter-toolkit/web"
)

// RedfishProxyServer struct implements HTTP server for proxy.
type RedfishProxyServer struct {
	logger    *slog.Logger
	server    *http.Server
	webConfig *web.FlagConfig
	redfish   *Redfish
}

// NewRedfishProxyServer creates new RedfishProxyServer struct instance.
func NewRedfishProxyServer(c *Config) *RedfishProxyServer {
	router := mux.NewRouter()
	server := &RedfishProxyServer{
		logger:  c.Logger,
		redfish: c.Redfish,
		server: &http.Server{
			Addr:              c.Web.Addresses[0],
			Handler:           router,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			ReadHeaderTimeout: 2 * time.Second, // slowloris attack: https://app.deepsource.com/directory/analyzers/go/issues/GO-S2112
		},
		webConfig: &web.FlagConfig{
			WebListenAddresses: &c.Web.Addresses,
			WebSystemdSocket:   &c.Web.WebSystemdSocket,
			WebConfigFile:      &c.Web.WebConfigFile,
		},
	}

	// If EnableDebugServer is true add debug endpoints
	if c.Web.EnableDebugServer {
		// pprof debug end points. Expose them only on localhost
		router.PathPrefix("/debug/").Handler(http.DefaultServeMux).Methods(http.MethodGet).Host("localhost")
	}

	// Handle metrics path
	router.PathPrefix("/").Handler(server.newProxyHandler())

	return server
}

// Start launches CEEMS exporter HTTP server.
func (s *RedfishProxyServer) Start() error {
	s.logger.Info("Starting " + appName)

	if err := web.ListenAndServe(s.server, s.webConfig, s.logger); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("Failed to Listen and Serve HTTP server", "err", err)

		return err
	}

	return nil
}

// Shutdown stops CEEMS exporter HTTP server.
func (s *RedfishProxyServer) Shutdown(ctx context.Context) error {
	s.logger.Info("Stopping " + appName)

	// First shutdown HTTP server to avoid accepting any incoming
	// connections
	// Do not return error here as we SHOULD ENSURE to close collectors
	// that might release any system resources
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("Failed to stop exporter's HTTP server")

		return err
	}

	return nil
}

// newProxyHandler creates a new handler for proxying requests to redfish targets.
func (s *RedfishProxyServer) newProxyHandler() *httputil.ReverseProxy {
	config := &rpConfig{
		logger:  s.logger.With("subsystem", "rp"),
		redfish: s.redfish,
	}

	return NewMultiHostReverseProxy(config)
}