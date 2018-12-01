package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/juju/errors"
)

var service = NewService()

type Service struct {
	client sarama.Client
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

func (s *Service) FetchKafkaMsgs(topic string, offset int64, count int64) ([]KafkaStreamData, error) {
	msgs := []KafkaStreamData{}

	hosts := strings.Split(GlobalConf.KafkaHost, ",")
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_1_0_0

	var err error
	s.client, err = sarama.NewClient(hosts, config)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer s.client.Close()

	Log.Info("[create kafka client ok!]\n")

	consumer, err := sarama.NewConsumerFromClient(s.client)
	if err != nil {
		return nil, errors.Trace(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			Log.Fatal("%v", err)
		}
	}()

	Log.Info("[create kafka consumer ok!][hosts]%v[topic]%s\n", hosts, topic)

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, offset)
	if err != nil {
		return nil, errors.Trace(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			Log.Fatal("%v", err)
		}
	}()

	Log.Info("[kafka consumer start!]\n")

	consumed := 0
	for {
		select {
		case m := <-partitionConsumer.Messages():
			Log.Info("[index]%d[event]%s[offset]%d\n", consumed, m.Value, m.Offset)

			offset = m.Offset

			msg := KafkaStreamData{Data: string(m.Value), Offset: offset}
			msgs = append(msgs, msg)

			consumed++
		case <-time.After(time.Second * 3):
			Log.Info("Ticker...")
		}
	}

	return nil, nil
}
