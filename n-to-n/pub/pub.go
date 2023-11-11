/**
 *
 * Simple publishers to demonstrate XSUB/XPUB usage.
 */

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pebbe/zmq4"
)

func pub(pubID int) {
	publisher, err := zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	publisher.Bind("tcp://*:5555") // Connect to XSUB
	log.Printf("PUB%d connects to XSUB at :5555\n", pubID)
	count := 0

	for {
		flag, err := publisher.Send(fmt.Sprintf("count: %d", count), 0)
		if err != nil {
			fmt.Printf("PUB%d Error: %s\n", pubID, err)
			break
		}
		fmt.Printf("PUB%d [flag, sends]: [%d, %d]\n", pubID, flag, count)

		count++
		time.Sleep(1 * time.Second)
	}
}

func main() {
	totalPublisher := 10
	for i := 0; i < totalPublisher; i++ {
		go pub(i)
	}

	select {}
}
