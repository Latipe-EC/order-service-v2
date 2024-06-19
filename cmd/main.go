package main

import (
	"fmt"
	server "latipe-order-service-v2/internal"
	"log"
	"runtime"
	"sync"
)

func main() {
	fmt.Println("Init application")
	defer log.Fatalf("[Info] Application has closed")

	numCPU := runtime.NumCPU()
	fmt.Printf("Number of CPU cores: %d\n", numCPU)

	serv, err := server.New()
	if err != nil {
		log.Fatalf("%s", err)
	}

	// Wrap each goroutine with a recovery handler
	var wg sync.WaitGroup

	wg.Add(1)
	go runWithRecovery(func() {
		defer wg.Done()
		serv.OrderTransactionSubscriber().ListenOrderEventQueue(&wg)
	})

	wg.Add(1)
	go runWithRecovery(func() {
		defer wg.Done()
		serv.RatingItemSubscriber().ListenRatingRatingQueue(&wg)
	})

	// API handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serv.App().Listen(serv.Config().Server.Port); err != nil {
			fmt.Printf("%s", err)
		}
	}()

	wg.Wait()
}

func runWithRecovery(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()
	fn()
}
