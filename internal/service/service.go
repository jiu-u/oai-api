package service

import (
	"github.com/jiu-u/oai-api/internal/repository"
	"github.com/jiu-u/oai-api/pkg/sid"
)

type Service struct {
	sid *sid.Sid
	tm  repository.Transaction
}

func NewService(
	sid *sid.Sid,
	tm repository.Transaction,
) *Service {
	return &Service{
		sid: sid,
		tm:  tm,
	}
}
