package utils

import (
	"main/logging"
	"time"
)

func RunWithProfiler(tag string, p func() error) error {
	startTime := time.Now()
	err := p()
	if err != nil {
		return err
	}
	endTime := time.Now()
	logging.DebugFormat("Run function %s for %.3f ms",
		tag, float64(endTime.Sub(startTime).Nanoseconds())/1000000.0)
	return nil
}
