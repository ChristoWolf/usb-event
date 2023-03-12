package usbevent_test

import (
	"log"

	usbevent "github.com/christowolf/usb-event"
)

func ExampleRegister() {
	n, err := usbevent.Register()
	if err != nil {
		log.Fatalf("registration failed: %v\n", err)
	}
	go func() {
		for e := range n.Channel {
			log.Printf("%+v\n\n", e)
		}
	}()
	n.Run()
}
