package main

import (
	"github.com/taranovegor/com.ligilo/cmd"
	"github.com/taranovegor/com.ligilo/internal/container"
	amqp "github.com/taranovegor/pkg.amqp"
)

func main() {
	sc := cmd.Init("amqp")

	router := sc.Get(container.AmqpController).(*amqp.Controller)
	router.Consume()

	select {}
}
