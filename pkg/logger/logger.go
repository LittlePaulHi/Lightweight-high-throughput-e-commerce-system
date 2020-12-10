//+build !debug

package logger

import (
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var (
	log        *logrus.Logger
	InitLog    *logrus.Entry
	APILog     *logrus.Entry
	ServiceLog *logrus.Entry
	RedisLog   *logrus.Entry
)

func init() {
	log = logrus.New()
	log.SetReportCaller(true)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	InitLog = log.WithFields(logrus.Fields{"API_Service": "Init"})
	APILog = log.WithFields(logrus.Fields{"API_Service": "API"})
	ServiceLog = log.WithFields(logrus.Fields{"API_Service": "Service"})
	RedisLog = log.WithFields(logrus.Fields{"Redis": "Cache"})
}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
