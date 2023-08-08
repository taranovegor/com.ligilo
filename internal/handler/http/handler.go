package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/taranovegor/com.ligilo/internal/config"
	"log"
	"net/http"
)

type Handler interface {
}

func MapRouter(
	router *chi.Mux,
	fallbackUrl string,
	link *LinkHandler,
) {
	router.Get(fmt.Sprintf("/{token:%s}", config.LinkTokenRegex), link.Get)
	router.NotFound(func(rw http.ResponseWriter, req *http.Request) {
		http.Redirect(rw, req, fallbackUrl, http.StatusFound)
	})
}

func response(rw http.ResponseWriter, code int, body interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	if _, err := rw.Write([]byte("{}")); err != nil {
		log.Println(err.Error())
	}
}

func redirect(rw http.ResponseWriter, req *http.Request, url string) {
	http.Redirect(rw, req, url, http.StatusMovedPermanently)
}
