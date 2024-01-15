package main

import (
	"fmt"
	server "latipe-order-service-v2/internal"
	"log"
)

func main() {
	fmt.Println("Init application")
	defer log.Fatalf("[Info] Application has closed")

	serv, err := server.New()
	if err != nil {
		log.Fatalf("%s", err)
	}

	if err := serv.App().Listen(serv.Config().Server.Port); err != nil {
		fmt.Printf("%s", err)
	}
}
