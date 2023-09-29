package client

import (
	"log"

	"go.temporal.io/sdk/client"
)

func NewClient() client.Client {
	return newLocalClient()
}

func GetNamespace() string {
	return "default"
}

func newLocalClient() client.Client {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	return c
}
