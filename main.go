package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/Shopify/sarama"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	flag.Usage = Usage
	conf := flag.String("c", "", "Conf file")
	flag.Parse()

	if *conf == "" {
		Usage()
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	err := service.InitService(*conf)
	if err != nil {
		Log.Fatal("[Init service failed]%s", err)
	}

	Log.Info("[%s]Read conf file ok!", GlobalConf.AppName)

	pidFile := GlobalConf.PidFilePath + "/" + GlobalConf.AppName + ".pid"
	err = ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0666)
	if err != nil {
		Log.Fatal("[WriteFile failed]%s", err)
	}

	Log.Info("[%s]Create pid file ok!", GlobalConf.AppName)

	engine := negroni.New()
	recoveryMid := negroni.NewRecovery()
	engine.Use(recoveryMid)

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/api/v1/data/kafka", DataGetHandler).Methods("GET")

	engine.UseHandler(router)
	Log.Info("[%s]Server start!", GlobalConf.AppName)

	addr := fmt.Sprintf(":%d", GlobalConf.Port)
	go http.ListenAndServe(addr, engine)

	go service.FetchKafkaMsgs(GlobalConf.KafkaTopic, sarama.OffsetOldest, 10)

	signal := InitSignal()
	HandleSignal(signal, nil)
}
