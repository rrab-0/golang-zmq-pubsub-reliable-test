/**
 *
 * Simple publisher with socket-monitor.
 * "Socket-Monitor" is used to count subscribers
 * available, when one subscriber disconnects,
 * publisher will stop sending messages.
 */

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pebbe/zmq4"
)

const (
	EXPECTED_SUB_COUNT = 10
)

var (
	queue        = make([]int, 0)
	isAbleToSend = false
	subCount     = 0
)

// Increment subCount if a sub connected,
// Decrement if a sub disconnected.
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
		eventType, addr, value, err := s.RecvEvent(0)
		if err != nil {
			log.Println(err)
			break
		}

		if eventType == zmq4.EVENT_HANDSHAKE_SUCCEEDED {
			subCount++
		}

		if eventType == zmq4.EVENT_DISCONNECTED {
			subCount--
		}

		log.Println(eventType, addr, value)
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
	err = publisher.Monitor("inproc://monitor.pub", zmq4.EVENT_ALL)
	if err != nil {
		log.Fatalln(err)
	}
	go pub_socket_monitor("inproc://monitor.pub")

	publisher.Bind("tcp://*:5555")
	count := 0

	// Main program loop
	for {
		if subCount == EXPECTED_SUB_COUNT {
			isAbleToSend = true
		}

		for isAbleToSend {
			if subCount != EXPECTED_SUB_COUNT {
				break
			}

			// Send "lost" messages
			// (when subs connected is not EXPECTED_SUB_COUNT,
			// messages are "lost" and sent into a queue)
			for len(queue) > 0 {
				if subCount != EXPECTED_SUB_COUNT {
					break
				}

				flag, err := publisher.Send(fmt.Sprintf("count: from queue %d", queue[0]), 0)
				if err != nil {
					fmt.Println("Queue Publisher Error:", err)
					break
				}
				fmt.Printf("[flag, sends from queue]: [%d, %d]\n", flag, queue[0])

				queue = queue[1:]
				time.Sleep(1 * time.Second)
			}

			if subCount != EXPECTED_SUB_COUNT {
				break
			}

			// Send message
			flag, err := publisher.Send(fmt.Sprintf("count: %d", count), 0)
			if err != nil {
				fmt.Println("Publisher Error:", err)
				break
			}
			fmt.Printf("[flag, sends]: [%d, %d]\n", flag, count)

			count++
			time.Sleep(1 * time.Second)
		}

		// For when a sub disconnect
		queue = append(queue, count)
		fmt.Println("Send to queue:", count)
		count++
		time.Sleep(1 * time.Second)
	}
}
