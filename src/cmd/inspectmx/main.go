package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"io.hyperd/inspectmx/config"
	"io.hyperd/inspectmx/logger"
	svc "io.hyperd/inspectmx/service"
	httpTransport "io.hyperd/inspectmx/transport/http"
)

var (
	Version   string
	BuildTime string
)

func init() {

	// load the config
	config.Instance()

	// scaffold the service
	svc.Instance()

	// setup logger
	logger.SetupLogger(config.Instance().Environment.LogLevel)

	logger.Info("orange ideas", logger.WithFields{
		"version":       Version,
		"build time":    BuildTime,
		"http port":     config.Instance().Server.HTTPPort,
		"development":   config.Instance().Environment.Development,
		"logging level": config.Instance().Environment.LogLevel,
	})
}

func main() {
	server := &http.Server{Addr: ":" + strconv.Itoa(config.Instance().Server.HTTPPort), Handler: httpTransport.Service()}

	logger.Info("http Server", logger.WithFields{"serving on port": server.Addr})

	// server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		logger.Info("http Server received signal", logger.WithFields{"action": "stop"})
		// shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal("http server graceful shutdown timed out: forcing exit", logger.WithFields{"action": "shutdown"})
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Fatal("http server exited without proper shutdown", logger.WithFields{"error": err.Error()})
		}
		serverStopCtx()
	}()

	// create channel for when http server gets shutdown for any reason
	shut := make(chan error)
	// start the server in the background and pass the channel for error handling
	go func(srv *http.Server, err chan error) {
		err <- srv.ListenAndServe()
	}(server, shut)

	// wait for server context to be stopped
	<-serverCtx.Done()
}
