package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func CloseServices() {
	if err := NATS_CONN.Drain(); err != nil {
		log.Fatal("Error during drain:", err)
	}
	NATS_CONN.Close()
}

func Shutdown() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s

		log.Println("shutting down...")
		CloseServices()

		os.Exit(0)
	}()
}
