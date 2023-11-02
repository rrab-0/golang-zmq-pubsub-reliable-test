/**
 *
 * Simple publisher with socket-monitor,
 * socket-monitor is used to handle sending messages
 * when subscriber disconnects.
 */

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pebbe/zmq4"
)

var (
	queue        = make([]int, 0)
	isAbleToSend = false
)

func pub_socket_monitor(addr string) {
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

		if a == zmq4.EVENT_ACCEPTED {
			isAbleToSend = true
		}

		if a == zmq4.EVENT_DISCONNECTED {
			isAbleToSend = false
		}
		log.Println(a, b, c)
	}
	s.Close()
}

func main() {
	// Publisher
	publisher, err := zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		fmt.Println("Publisher Error:", err)
	}
	defer publisher.Close()

	// PUB socket monitor, all events
	err = publisher.Monitor("inproc://monitor.sub", zmq4.EVENT_ALL)
	if err != nil {
		log.Fatalln(err)
	}
	go pub_socket_monitor("inproc://monitor.sub")

	publisher.Bind("tcp://*:5555")
	count := 0
	for {
		for isAbleToSend {
			for len(queue) > 0 {
				flag, err := publisher.Send(fmt.Sprintf("count: from queue %d", queue[0]), 0)
				if err != nil {
					fmt.Println("Queue Publisher Error:", err)
					break
				}
				fmt.Printf("[flag, sends from queue]: [%d, %d]\n", flag, queue[0])

				queue = queue[1:]
				time.Sleep(1 * time.Second)
			}

			flag, err := publisher.Send(fmt.Sprintf("count: %d", count), 0)
			if err != nil {
				fmt.Println("Publisher Error:", err)
				break
			}
			fmt.Printf("[flag, sends]: [%d, %d]\n", flag, count)

			count++
			time.Sleep(1 * time.Second)
		}

		queue = append(queue, count)
		fmt.Println("Send to queue:", count)
		count++
		time.Sleep(1 * time.Second)
	}
}
