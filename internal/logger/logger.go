package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func New(lvl string) (*zap.Logger, error) {
	const op = "logger.New"

	level, err := zap.ParseAtomicLevel(lvl)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse logger level: %v", op, err)
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = level

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build logger: %v", op, err)
	}

	return logger, nil
}