package http

import (
	"net/http"
)

type Handler interface {
}

func MapRouter(
	router *http.ServeMux,
	link *LinkHandler,
) {
	router.HandleFunc("GET /{token}", link.Get)
}

func redirect(rw http.ResponseWriter, req *http.Request, url string) {
	http.Redirect(rw, req, url, http.StatusMovedPermanently)
}
