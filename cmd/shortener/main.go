package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sur1k1/go-url-shortener/internal/app/rest"
	"github.com/sur1k1/go-url-shortener/internal/config"
	"github.com/sur1k1/go-url-shortener/internal/logger"
	storage "github.com/sur1k1/go-url-shortener/internal/repository/memstorage"
	"github.com/sur1k1/go-url-shortener/internal/service"
)

func main() {
	// Getting a configuration
	cf := config.MustGetConfig()

	// Logger init
	log, err := logger.New(cf.LogLevel)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
	}
	
	// Load storage
	s, err := storage.NewStorage(log, cf.FilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
	}

	// Init service repository
	serviceRepo := service.New(s)

	// Init application
	application := rest.New(log, serviceRepo, cf)

	// Start server
	go func() {
		application.MustRun()
	}()

	// Gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<- stop

	// Close file on shutdown
	err = s.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
	}

	log.Info("Gracefully stopped")
}
