package main

import (
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

// InitSignal register signals handler.
func InitSignal() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	return c
}

// HandleSignal fetch signal from chan then do exit or reload.
func HandleSignal(c chan os.Signal, arg interface{}) {
	// Block until a signal is received.
	for {
		s := <-c
		Log.Info("Get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			// Exit
			return
		case syscall.SIGHUP:
			// TODO
			// Reload
			// return
		default:
			return
		}
	}
}

// InitPprof start http pprof.
func InitPprof(pprofBind []string) {
	pprofServeMux := http.NewServeMux()
	pprofServeMux.HandleFunc("/debug/pprof/", pprof.Index)
	pprofServeMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pprofServeMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pprofServeMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	for _, addr := range pprofBind {
		go func() {
			if err := http.ListenAndServe(addr, pprofServeMux); err != nil {
				panic(err)
			}
		}()
	}
}
