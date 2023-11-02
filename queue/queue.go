package main

import (
	"github.com/pebbe/zmq4"
)

func main() {
	context, _ := zmq4.NewContext()

	xpub, _ := context.NewSocket(zmq4.XPUB)
	xsub, _ := context.NewSocket(zmq4.XSUB)

	xpub.Bind("tcp://*:5556")
	xsub.Bind("tcp://*:5555")

	go zmq4.Proxy(xpub, xsub, nil)
	select {}
}
