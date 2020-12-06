//+build !debug

package logger

import (
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var APILog *logrus.Entry
var ServiceLog *logrus.Entry
var InitLog *logrus.Entry

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

	APILog = log.WithFields(logrus.Fields{"API_Service": "API"})
	ServiceLog = log.WithFields(logrus.Fields{"API_Service": "Service"})
	InitLog = log.WithFields(logrus.Fields{"API_Service": "Init"})

}

func SetLogLevel(level logrus.Level) {
	log.SetLevel(level)
}

func SetReportCaller(bool bool) {
	log.SetReportCaller(bool)
}
