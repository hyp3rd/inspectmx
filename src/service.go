package inspectmx

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"io.hyperd/inspectmx/config"
	"io.hyperd/inspectmx/helpers"
	"io.hyperd/inspectmx/logger"
)

var (
	once        sync.Once
	imxInstance *inspectMX
)

type inspectMX struct {
	AllowedProviders []string `json:"allowed_providers"`
	AllowedServers   []string `json:"allowed_servers"`
}

func instance() *inspectMX {
	once.Do(func() {
		imxInstance = &inspectMX{
			AllowedProviders: config.Instance().Email.AllowedProviders,
		}

		for _, provider := range config.Instance().Email.AllowedProviders {
			mxrecords, err := helpers.MxLookup(provider)
			if err != nil {
				logger.Error("failed to resolve a provider", logger.WithFields{
					"provider": provider,
				})
			}
			for _, mx := range mxrecords {
				imxInstance.AllowedServers = append(imxInstance.AllowedServers, mx.Host)
			}

			imxInstance.AllowedServers = append(imxInstance.AllowedServers, config.Instance().Email.ExtraMX...)
		}

		logger.Info("allowed email providers", logger.WithFields{
			"list": imxInstance.AllowedProviders,
		})

		logger.Info("allowed MX servers", logger.WithFields{
			"list": imxInstance.AllowedServers,
		})

	})
	return imxInstance
}

type Service interface {
	GetAllowedProviders(ctx context.Context) []string
	GetAllowedServers(ctx context.Context) []string
	Verify(email string, ctx context.Context) (*string, error)
}

type service struct{}

// New returns a concrete Service
func New() (Service, error) {
	// return
	return &service{}, nil
}

func (s service) GetAllowedProviders(ctx context.Context) []string {
	return instance().AllowedProviders
}

func (s service) GetAllowedServers(ctx context.Context) []string {
	return instance().AllowedServers
}

func (s service) Verify(email string, ctx context.Context) (*string, error) {
	var err error
	at := strings.LastIndex(email, "@")
	if at >= 0 {
		username, domain := email[:at], email[at+1:]

		logger.Debug("verifying address", logger.WithFields{
			"email":    email,
			"domain":   domain,
			"username": username,
		})

		mxrecords, err := helpers.MxLookup(domain)

		if err != nil {
			logger.Error("MX lookup failed", logger.WithFields{
				"domain": domain,
			})
			err = fmt.Errorf("mx lookup failed for %s", domain)
			return nil, err
		}

		for _, m := range mxrecords {
			if helpers.FindElementInArray(instance().AllowedServers, m.Host) {
				return &m.Host, nil
			}
		}
		err = fmt.Errorf("%s is not a valid provider", domain)
		return &domain, err

	} else {
		logger.Error("invalid email address", logger.WithFields{
			"email": email,
		})
		err = fmt.Errorf("invalid email address: %s", email)
	}

	return nil, err
}
