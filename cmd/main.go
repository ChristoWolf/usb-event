package main

import (
	"log"

	usbevent "github.com/christowolf/usb-event"
)

func main() {
	n, err := usbevent.Register()
	if err != nil {
		log.Printf("registration failed: %v\n", err)
		return
	}
	go func() {
		for e := range n.Channel {
			log.Printf("%v", e)
		}
	}()
	n.Run()
}
