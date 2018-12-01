package main

import (
	"fmt"
)

var GlobalConf = &Conf{}

type Conf struct {
	AppName     string `goconf:"core:appName"`
	LogName     string `goconf:"core:logName"`
	LogRotate   string `goconf:"core:logRotate"`
	PidFilePath string `goconf:"core:pidFilePath"`
	Port        int64  `goconf:"core:port"`
	PprofHost   string `goconf:"core:pprofHost"`
	KafkaHost   string `goconf:"kafka:host"`
	KafkaTopic  string `goconf:"kafka:topic"`
}

func (c *Conf) String() string {
	if c == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Conf(%+v)", *c)
}
