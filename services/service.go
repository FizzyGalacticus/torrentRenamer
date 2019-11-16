package services

import (
	"torrentRenamer"
	"torrentRenamer/config"
)

var (
	registeredServices []Service
	defaultService     *Service
)

type Service interface {
	Name() string
	Search(*torrentRenamer.Video) (torrentRenamer.Video, error)
	IsAvailable() bool
	GetNewName(*torrentRenamer.Video) (string, error)
	GetServiceName() string
}

func RegisterService(service Service) {
	registeredServices = append(registeredServices, service)
}

func GetRegistedServices() []Service {
	return registeredServices
}

func IsDefault(service Service) bool {
	config := config.GetConfig()

	return config.DefaultService == service.GetServiceName()
}

func GetDefaultService() *Service {
	if defaultService == nil {
		for _, service := range GetRegistedServices() {
			if IsDefault(service) {
				defaultService = &service

				break
			}
		}
	}

	return defaultService
}

func GetDefaultServiceResults(video *torrentRenamer.Video) (string, error) {
	service := (*GetDefaultService())

	var result string
	var err error

	if service != nil {
		if service.IsAvailable() {
			result, err = service.GetNewName(video)
		}
	}

	return result, err
}
