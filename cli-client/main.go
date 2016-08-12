package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ganners/gossip"
)

func main() {
	// Create a new gossip server
	server, err := gossip.NewServer(
		"cli-client",
		"The client which launches syncronous client requests",
		"0.0.0.0", "8003",
		gossip.NewSilentLogger(),
	)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}

	// Launch a new client
	client := NewCLIClient(server)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Read or write? [w/r]: ")
			text, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Unable to read input:", err)
				continue
			}
			text = strings.TrimSpace(text)

			if text == "w" {
				fmt.Print("Enter payload to write: ")
				payload, err := reader.ReadString('\n')
				if err != nil {
					log.Println("Unable to read input:", err)
					continue
				}
				payload = strings.TrimSpace(payload)

				fmt.Print("Enter id to write value to: ")
				id, err := reader.ReadString('\n')
				if err != nil {
					log.Println("Unable to read input:", err)
					continue
				}
				id = strings.TrimSpace(id)

				key, err := client.Store([]byte(id), []byte(payload))
				if err != nil {
					log.Println("Could not store payload:", err)
					continue
				}
				log.Println("[Success] Your key to retrieve in future is: ", string(key))
			} else if text == "r" {
				fmt.Print("Enter your ID: ")
				id, err := reader.ReadString('\n')
				if err != nil {
					log.Println("Unable to read input: %s", err)
					continue
				}
				id = strings.TrimSpace(id)

				fmt.Print("Enter your encryption key: ")
				key, err := reader.ReadString('\n')
				if err != nil {
					log.Println("Unable to read input:", err)
					continue
				}
				key = strings.TrimSpace(key)

				payload, err := client.Retrieve([]byte(id), []byte(key))
				if err != nil {
					log.Println("Could not retrieve payload:", err)
					continue
				}
				log.Println("[Success] Your original payload was: ", string(payload))
			}
		}
	}()

	<-server.Start()
	log.Println("Terminating client server")
}
