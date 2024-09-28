package helper

import (
	"go.uber.org/zap"
)

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func ErrorLog(err error, log *zap.Logger) {
	if err != nil {
		log.Error(err.Error(), zap.Error(err))
	}
}
