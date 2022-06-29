//go:build !test

package main

import (
	"context"
	"log"
	"main/v2/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	startserver()
}

func startserver() {
	//var q deque.Deque
	//q.PushBack(1)
	ctx := context.Background()

	serverDoneChan := make(chan os.Signal, 1)
	signal.Notify(serverDoneChan, os.Interrupt, syscall.SIGTERM)

	srv := server.New(":9091")

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	log.Println("server started")

	<-serverDoneChan

	srv.Shutdown(ctx)
	log.Println("server stopped")
}
