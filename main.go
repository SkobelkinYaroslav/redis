package main

import (
	"context"
	"fmt"
	"log"
	"redis/client"
	"redis/server"
)

func main() {
	go func() {
		server := server.New(server.Config{})

		server.Start()
	}()

	for i := 0; i < 10; i++ {
		client := client.New("localhost:8080")

		if err := client.Set(context.Background(), fmt.Sprintf("foo_%d", i), fmt.Sprintf("bar_%d", i)); err != nil {
			log.Fatal(err)
		}
		fmt.Println("set complete")
		str, err := client.Get(context.Background(), fmt.Sprintf("foo_%d", i))
		if err != nil {
			log.Fatal(err)
		}

		log.Println("got from DB <= ", str)
		client.Close()

	}

}
