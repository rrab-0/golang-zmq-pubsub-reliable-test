/**
 *
 * Simple subscribers to demonstrate XSUB/XPUB usage.
 */

package main

import (
	"log"

	"github.com/pebbe/zmq4"
)

func sub(subName string) {
	subscriber, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		log.Fatal(err)
	}
	defer subscriber.Close()

	subscriber.Connect("tcp://localhost:5556")
	subscriber.SetSubscribe("") // Subscribe to all topics
	log.Printf("SUB%s listens on :5556\n", subName)

	for {
		message, err := subscriber.Recv(0)
		if err != nil {
			log.Printf("SUB%s Error: %s\n", subName, err)
			continue
		}
		log.Printf("SUB%s Received: %s\n", subName, message)
	}
}

func main() {
	go sub("1")
	go sub("2")
	go sub("3")
	go sub("4")
	go sub("5")
	go sub("6")
	go sub("7")
	go sub("8")
	go sub("9")
	go sub("10")

	select {}
}
