package util

import "goauth/log"

func LogError(err error, logger log.Logger)  {
	if err != nil {
		logger.Error(err.Error())
	}
}