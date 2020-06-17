//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
package registry

import (
	"github.com/sirupsen/logrus"
)

//counterfeiter:generate . RegistryAdder
type RegistryAdder interface {
	AddToRegistry(AddToRegistryRequest) error
}

func NewRegistryAdder(logger *logrus.Entry) RegistryAdder {
	return RegistryUpdater{
		Logger: logger,
	}
}

//counterfeiter:generate . RegistryDeleter
type RegistryDeleter interface {
	DeleteFromRegistry(DeleteFromRegistryRequest) error
}

func NewRegistryDeleter(logger *logrus.Entry) RegistryDeleter {
	return RegistryUpdater{
		Logger: logger,
	}
}
