package usbevent_test

import (
	"log"

	usbevent "github.com/christowolf/usb-event"
)

func ExampleRegister() {
	n, err := usbevent.Register()
	if err != nil {
		log.Printf("registration failed: %v\n", err)
		return
	}
	go func() {
		for e := range n.Channel {
			log.Printf("%v\n\n", e)
			log.Printf("Device name: %s\n\n", e.DeviceName)
		}
	}()
	n.Run()
}
