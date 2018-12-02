package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

// Usage
func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

type Click struct {
	ClickID    int64  `json:"click_id"`
	UserID     int64  `json:"user_id"`
	AdsID      int64  `json:"ads_id"`
	ClickPrice int64  `json:"click_price"`
	CreateTime string `json:"create_time"`
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func randInt64(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func getDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func ProduceKafkaMsgs(hosts []string, topic string, count int) error {
	producer, err := sarama.NewSyncProducer(hosts, nil)
	if err != nil {
		return errors.Trace(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal("%v", err)
		}
	}()

	log.Infof("[create kafka sync producer ok!]\n")

	t := time.Now()
	for i := 0; i < count; i++ {
		clickID := int64(i + 1)
		userID := randInt64(101, 120)
		adsID := randInt64(1001, 1010)
		clickPrice := randInt64(1, 10) * 100
		click := Click{clickID, userID, adsID, clickPrice, getDateTime(t.Add(time.Second * time.Duration(i)))}

		data, err := json.Marshal(&click)
		if err != nil {
			log.Fatalf("[mock stream data marshal failed]%v", err)
		}

		msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(data), Partition: 0}
		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Fatal("%v", err)
		} else {
			log.Infof("[message sent to partition %d at offset %d]\n", partition, offset)
		}
	}

	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Usage = Usage
	confFile := flag.String("c", "", "Conf file")
	flag.Parse()

	if *confFile == "" {
		Usage()
	}

	var err error
	conf := NewConfig()
	err = conf.Parse(*confFile)
	if err != nil {
		log.Fatalf("Init config fail - parse fail - " + err.Error())
	}

	err = conf.Unmarshal(GlobalConf)
	if err != nil {
		log.Fatalf("Init config fail - Unmarshal fail - " + err.Error())
	}

	ProduceKafkaMsgs([]string{GlobalConf.Host}, GlobalConf.Topic, GlobalConf.Count)
}
