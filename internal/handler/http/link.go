package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/taranovegor/com.ligilo/internal/domain"
	"net/http"
)

type LinkHandler struct {
	Handler
	link domain.LinkRepository
}

func NewLinkHandler(
	link domain.LinkRepository,
) *LinkHandler {
	return &LinkHandler{
		link: link,
	}
}

func (hdlr LinkHandler) Get(rw http.ResponseWriter, req *http.Request) {
	token := chi.URLParam(req, "token")
	link, err := hdlr.link.GetByToken(token)
	if err != nil {
		response(rw, http.StatusNotFound, err)
	} else {
		redirect(rw, req, link.Location)
	}
}
