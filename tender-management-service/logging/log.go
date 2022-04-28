package logging

import (
	"github.com/op/go-logging"
	"os"
	"sync"
)

var lock sync.RWMutex
var log = logging.MustGetLogger("all")
var format = logging.MustStringFormatter(
	`[%{time:2006-01-02T15:04:05.999}][level=%{level}][class=%{shortpkg}:%{shortfile}]%{message}`,
)

func init() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend1, format)
	backend1Leveled := logging.AddModuleLevel(backend1)

	logLevel := os.Getenv("LOG_LEVEL")

	backend1Leveled.SetLevel(parseLogLevel(logLevel), "")
	logging.SetBackend(backend1Leveled, backend2Formatter)
	logging.SetLevel(logging.DEBUG, "all")
}

func parseLogLevel(logLevel string) logging.Level {
	switch logLevel {
	case "DEBUG":
		return logging.DEBUG
	case "INFO":
		return logging.INFO
	case "NOTICE":
		return logging.NOTICE
	case "WARNING":
		return logging.WARNING
	case "ERROR":
		return logging.ERROR
	case "CRITICAL":
		return logging.CRITICAL
	default:
		return logging.INFO
	}
}

func InfoFormat(message string, params ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	log.Infof(message, params...)
}
func Info(params ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	log.Info(params...)
}
func ErrorFormat(message string, params ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	log.Errorf(message, params...)
}
func Error(params ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	log.Error(params...)
}
func DebugFormat(message string, params ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	log.Debugf(message, params...)
}
func Debug(params ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	log.Debug(params...)
}

func FatalFormat(message string, params ...interface{}) {
	lock.Lock()
	defer lock.Unlock()
	log.Fatalf(message, params...)
}
