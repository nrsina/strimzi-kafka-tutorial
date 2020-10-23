package log

import (
	"go.uber.org/zap"
)

//Global Zap Sugared Logger
var Logger *zap.SugaredLogger

type LogProfile string

const (
	Production LogProfile = "production"
	Development = "development"
)

func InitLogger(profile LogProfile) error {
	var logger *zap.Logger
	var err error
	switch profile {
	case Production:
		logger, err = zap.NewProduction()
	case Development:
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return err
	}
	Logger = logger.Sugar()
	return nil
}
