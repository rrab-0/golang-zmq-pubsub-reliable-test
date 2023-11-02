package main

import (
	"fmt"
	"log"

	"github.com/pebbe/zmq4"
)

func sub_socket_monitor(addr string) {
	s, err := zmq4.NewSocket(zmq4.PAIR)
	if err != nil {
		log.Fatalln(err)
	}
	err = s.Connect(addr)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		a, b, c, err := s.RecvEvent(0)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(a, b, c)
	}
	s.Close()
}

func main() {
	// Subscriber
	subscriber, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		log.Fatal(err)
	}
	defer subscriber.Close()

	// SUB socket monitor, all events
	err = subscriber.Monitor("inproc://monitor.sub", zmq4.EVENT_ALL)
	if err != nil {
		log.Fatalln(err)
	}
	go sub_socket_monitor("inproc://monitor.sub")

	// Connect initially
	subscriber.Connect("tcp://localhost:5555")
	subscriber.SetSubscribe("") // Subscribe to all topics

	for {
		message, err := subscriber.Recv(0)
		if err != nil {
			log.Printf("SUB Error: %s\n", err)
			continue
		}
		fmt.Println("Received:", message)
	}
}
