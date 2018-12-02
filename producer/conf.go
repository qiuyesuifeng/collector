package main

import "fmt"

var GlobalConf = &Conf{}

type Conf struct {
	Host  string `goconf:"kafka:host"`
	Count int    `goconf:"kafka:count"`
	Topic string `goconf:"kafka:topic"`
}

func (c *Conf) String() string {
	if c == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Conf(%+v)", *c)
}
