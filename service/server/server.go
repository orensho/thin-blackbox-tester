package server

import (
	"context"
	"github.com/orensho/thin-slack-blackbox-tester/service/config"
	"github.com/orensho/thin-slack-blackbox-tester/service/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

type BlackboxServer struct {
	http.Server

	shutdownChan chan bool
	isInShutdown uint32

	serverSettings *config.ServerSettings
}

func Must(err error) {
	if err != nil {
		log.WithError(err).Panic("oops...")
	}
}

func NewFgBlackboxServer(serverSettings *config.ServerSettings, serviceFactory *service.Factory) *BlackboxServer {
	log.Info("Starting thin-blackbox-tester ...")

	server := &BlackboxServer{
		Server: http.Server{
			Addr: serverSettings.LocalListenIP + ":" + serverSettings.LocalListenPort,
		},
		shutdownChan:   make(chan bool),
		serverSettings: serverSettings,
	}

	return server
}

// WaitForShutdown will block until either SIGINT/SIGTERM is received or
// until /shutdown is called.
// Once one of the above happens the server will start shutting down gracefully,
// meaning it will wait for all requests on connections to complete and start closing
// those that are idle. It will wait at most SHUTDOWN_GRACE_PERIOD before the server
// this function returns.
func (as *BlackboxServer) WaitForShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	// Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-irqSig:
		log.WithField("signal", sig).Info("Received termination signal - shutting down.")
	case <-as.shutdownChan:
		log.Info("Received shutdown request.")
	}

	log.Info("Stopping http server ...")

	// Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), as.serverSettings.ParsedShutDownGracePeriod)
	defer cancel()

	// shutdown the server: wait for requests to complete and close idle connections
	if err := as.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server shutdown procedure failed.")
	}
}

func (as *BlackboxServer) Start() error {
	log.Infof("starting http server on %s:%s", as.serverSettings.LocalListenIP, as.serverSettings.LocalListenPort)

	return as.ListenAndServe()
}

func (as *BlackboxServer) TriggerShutdown() {
	// Do nothing if shutdown request already issued:
	// if as.isInShutdown == 0 then set to 1 and return true
	// if as.isInShutdown != 0 do nothing and return false
	if !atomic.CompareAndSwapUint32(&as.isInShutdown, 0, 1) {
		log.Info("Shutdown is already in progress...")

		return
	}

	go func() {
		as.shutdownChan <- true
	}()
}
