package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/taranovegor/com.ligilo/cmd"
	"github.com/taranovegor/com.ligilo/internal/config"
	"github.com/taranovegor/com.ligilo/internal/container"
	"log"
	"net/http"
)

func main() {
	sc := cmd.Init("http")

	router := sc.Get(container.HttpRouter).(*chi.Mux)

	err := http.ListenAndServe(fmt.Sprintf(":%s", config.GetEnv(config.HttpPort)), router)
	if err != nil {
		log.Panic(err)
	}
}
