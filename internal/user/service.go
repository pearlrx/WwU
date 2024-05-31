package user

import (
	"context"
	"techpoint/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

func (s Service) Create(ctx context.Context, dto createUserDTO) (u User, err error) {
	//TODO for a next one
	return
}
