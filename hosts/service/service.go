package service

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ninnemana/fuz/hosts"
	"github.com/pkg/errors"
)

// New creates a new hosts service
func New() (*service, error) {

	h := hosts.Hosts{
		Path: os.ExpandEnv(filepath.FromSlash(hostsFilePath)),
	}

	err := h.Load()
	if err != nil {
		return nil, err
	}

	return &service{
		Hosts: h,
	}, nil
}

func (s *service) Get(context.Context, hosts.GetParams) (*hosts.Record, error) {
	return nil, errors.Errorf("not implemented")
}

func (s *service) Set(context.Context, hosts.Record) (*hosts.Record, error) {
	return nil, errors.Errorf("not implemented")
}

func (s *service) List(context.Context, *hosts.ListParams) ([]hosts.Record, error) {
	return s.Hosts.Records, nil
}

func (s *service) Delete(context.Context, hosts.GetParams) error {
	return errors.Errorf("not implemented")
}

type service struct {
	Hosts hosts.Hosts
}
