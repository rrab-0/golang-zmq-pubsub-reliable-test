/**
 *
 * Simple XSUB/XPUB usage.
 */

package main

import (
	"fmt"
	"log"

	"github.com/pebbe/zmq4"
)

func main() {
	xsub, err := zmq4.NewSocket(zmq4.XSUB)
	if err != nil {
		log.Fatal(err)
	}
	defer xsub.Close()
	xsub.Connect("tcp://localhost:5555")

	xpub, err := zmq4.NewSocket(zmq4.XPUB)
	if err != nil {
		log.Fatal(err)
	}
	defer xpub.Close()
	xpub.Bind("tcp://*:5556")

	// msgChan := make(chan string)
	// go func() {
	// 	for {
	// 		message, err := xsub.Recv(zmq4.DONTWAIT)
	// 		if err != nil {
	// 			log.Printf("XSUB Error: %s\n", err)
	// 			continue
	// 		}
	// 		msgChan <- message
	// 		log.Printf("XSUB Received: %s\n", message)
	// 		log.Println("XSUB Forwarding message to XPUB...")
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		_, err := xpub.Send(<-msgChan, zmq4.DONTWAIT)
	// 		if err != nil {
	// 			fmt.Printf("XPUB Error: %s\n", err)
	// 			continue
	// 		}
	// 		fmt.Printf("XPUB sent message to SUBs: %s\n", <-msgChan)
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	fmt.Println("Intermediary is up and running...")
	err = zmq4.Proxy(xsub, xpub, nil)
	if err != nil {
		log.Fatal(err)
	}
}
