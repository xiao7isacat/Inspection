package signal

import (
	"context"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
)

func SetupStopSignalContext() (*errgroup.Group, <-chan struct{}) {
	group, ctx := errgroup.WithContext(Context(SetupStopSignalHandler()))
	return group, ctx.Done()
}

func Context(signal <-chan struct{}) context.Context {
	ret, cancel := context.WithCancel(context.Background())
	go func() {
		<-signal
		cancel()
	}()
	return ret
}

var onlyOneSignalHandler = make(chan struct{})
var shutdownHandler chan os.Signal

func SetupStopSignalHandler() <-chan struct{} {
	close(onlyOneSignalHandler) // panics when called twice

	shutdownHandler = make(chan os.Signal, 2)

	stop := make(chan struct{})
	signal.Notify(shutdownHandler, shutdownSignals...)
	go func() {
		<-shutdownHandler
		close(stop)
		<-shutdownHandler
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
