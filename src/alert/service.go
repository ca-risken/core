package main

import (
	"github.com/CyberAgent/mimosa-core/proto/alert"
)

type alertService struct {
	repository alertRepository
}

func newAlertService() alert.AlertServiceServer {
	return &alertService{
		repository: newAlertRepository(),
	}
}
