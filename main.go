package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func main() {
	opts := &server.Options{

	}

	ns, err := server.NewServer(opts)
	if err != nil {
		panic(err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		panic("not ready for connection")
	}

	// connect to the server
	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		panic(err)
	}

	subject := `my-subject`

	nc.Subscribe(subject, func(msg *nats.Msg) {
		data := string(msg.Data)
		fmt.Println(data)

		ns.Shutdown()
	})

	nc.Publish(subject, []byte("Hello embedded NATS !"))

	ns.WaitForShutdown()
}