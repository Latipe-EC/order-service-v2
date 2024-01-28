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

	//subscriber
	var wg sync.WaitGroup
	wg.Add(1)
	go serv.OrderTransactionSubscriber().ListenOrderEventQueue(&wg)

	//api handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serv.App().Listen(serv.Config().Server.Port); err != nil {
			fmt.Printf("%s", err)
		}
	}()

	wg.Wait()
}
