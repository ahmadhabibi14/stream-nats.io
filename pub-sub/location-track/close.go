package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func CloseServices() {
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