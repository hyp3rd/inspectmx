package service

import (
	"sync"

	"io.hyperd/inspectmx"
	"io.hyperd/inspectmx/logger"
)

var (
	once            sync.Once
	serviceInstance Service
)

// service implements the orangeideas.Repository
type Service struct {
	Inspector inspectmx.Service
}

func Instance() Service {
	once.Do(func() {
		logger.Info("bootstraping the Inspector service", nil)
		srv, err := inspectmx.New()

		if err != nil {
			logger.Fatal("failed instanciating the Inspector service", logger.WithFields{
				"error": err,
			})
		}

		serviceInstance = Service{
			Inspector: srv,
		}
	})

	return serviceInstance
}
