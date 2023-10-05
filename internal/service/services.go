package service

import (
	"context"

	"go.uber.org/zap"
)

type Services struct {
	logger *zap.Logger
}

func NewServices(ctx context.Context, logger *zap.Logger) *Services {
	return &Services{
		logger: logger,
	}
}
