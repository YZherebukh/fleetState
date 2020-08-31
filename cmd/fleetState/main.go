package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/fleetState/queue"
	"github.com/fleetState/web"

	"github.com/fleetState/model"
	"github.com/fleetState/store"

	"github.com/fleetState/config"
	"github.com/fleetState/logger"
	"github.com/fleetState/web/middleware"
	vhcl "github.com/fleetState/web/vehicle"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call: %+v \n", oscall)
		cancel()
	}()

	err := loadService(ctx)
	if err != nil {
		cancel()
		panic(err)
	}
}

func loadService(ctx context.Context) error {
	config, err := config.New()
	if err != nil {
		return fmt.Errorf("load application failed. error %s", err)
	}

	logger := logger.New(config.Logger())

	state := store.New()
	vehicle := model.New(state)

	stream := queue.NewStream(ctx)
	stateQ := queue.NewState(ctx, *stream)

	router := mux.NewRouter().StrictSlash(true)
	resp := web.NewResponse(logger)
	middleware := middleware.New(logger)

	vhcl.NewHandler(ctx, router, logger, middleware, resp, vehicle, *stateQ, *stream)

	server := &http.Server{
		Addr:    config.Service().Port,
		Handler: handlers.CORS()(router),
	}

	errCh := make(chan error)

	go startServer(ctx, logger, errCh, server)

	select {
	case <-ctx.Done():
		logger.Infof(ctx, "stopping service")
		err = server.Shutdown(ctx)
		logger.Infof(ctx, "service has been stopped")
		return nil
	case err = <-errCh:
		logger.Errorf(ctx, "service has been stopped with error: %s", err.Error())
		return err
	}
}

func startServer(ctx context.Context, l logger.Logger, errCh chan error, server *http.Server) chan error {
	l.Infof(ctx, "Starting HTTP listener...")
	errCh <- server.ListenAndServe()
	return errCh
}
