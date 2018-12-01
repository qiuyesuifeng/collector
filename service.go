package main

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

// Init service
func (s *Service) InitService(configPath string) error {
	// init conf
	var err error
	conf := NewConfig()
	err = conf.Parse(configPath)
	if err != nil {
		fmt.Println("Init config fail - parse fail - " + err.Error())
		return errors.Trace(err)
	}

	err = conf.Unmarshal(GlobalConf)
	if err != nil {
		fmt.Println("Init config fail - Unmarshal fail - " + err.Error())
		return errors.Trace(err)
	}

	if GlobalConf.LogRotate == "day" {
		SetRotateByDay()
	} else if GlobalConf.LogRotate == "hour" {
		SetRotateByHour()
	}

	err = Log.SetOutputByName(GlobalConf.LogName)
	if err != nil {
		fmt.Println("Log SetOutputByName fail - " + err.Error())
		return errors.Trace(err)
	}

	Log.Debug(GlobalConf.String())

	// init pprof
	pprofs := strings.Split(GlobalConf.PprofHost, ",")
	InitPprof(pprofs)
	Log.Info("[%s]Pprof service ok!", GlobalConf.AppName)

	return nil
}
