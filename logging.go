package main

import (
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

func initLogger(logLevel string) {

	level, _ := log.ParseLevel(logLevel)

	formatter := &log.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05.000",
		PadLevelText:     true,
		QuoteEmptyFields: true,
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "sensemon.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      level,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})

	if err != nil {
		log.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	log.SetLevel(level)
	log.SetFormatter(formatter)

	if log.GetLevel() == log.TraceLevel {
		log.SetReportCaller(true)
		formatter.CallerPrettyfier = LogPrettyTrace
	}
	log.AddHook(rotateFileHook)
}

func LogPrettyTrace(f *runtime.Frame) (function string, file string) {
	//function = path.Base(f.Function)
	function = ""
	file = " " + path.Base(f.File) + ":" + strconv.Itoa(f.Line)
	return function, file
}
