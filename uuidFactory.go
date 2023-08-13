package main

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type uuidFactory struct{}

var SharedUUIDFactory = uuidFactory{}

func (i uuidFactory) Create() (string, error) {

	uuidIdentifier, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	generatedIdentifier := uuidIdentifier.String()

	log.Debugf("Generated identifier: %v", generatedIdentifier)

	return generatedIdentifier, nil
}
