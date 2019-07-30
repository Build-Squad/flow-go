// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package collect

import (
	"github.com/dapperlabs/bamboo-node/internal/roles/collect/config"
	"github.com/dapperlabs/bamboo-node/internal/roles/collect/controller"
	"github.com/dapperlabs/bamboo-node/internal/roles/collect/txpool"
	"github.com/sirupsen/logrus"
)

// Injectors from wire.go:

func InitializeServer() (*Server, error) {
	configConfig := config.New()
	logger := logrus.New()
	txPool := txpool.New()
	controllerController := controller.New(txPool, logger)
	server, err := NewServer(configConfig, logger, controllerController)
	if err != nil {
		return nil, err
	}
	return server, nil
}
