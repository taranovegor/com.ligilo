package http

import (
	"github.com/taranovegor/com.ligilo/internal/config"
	"github.com/taranovegor/com.ligilo/internal/domain"
	"net/http"
	"regexp"
)

type LinkHandler struct {
	Handler
	link        domain.LinkRepository
	fallbackUrl string
}

func NewLinkHandler(
	link domain.LinkRepository,
	fallbackUrl string,
) *LinkHandler {
	return &LinkHandler{
		link:        link,
		fallbackUrl: fallbackUrl,
	}
}

func (hdlr LinkHandler) Get(rw http.ResponseWriter, req *http.Request) {
	token := req.PathValue("token")
	if match, _ := regexp.MatchString(config.LinkTokenRegex, token); !match {
		hdlr.redirectToFallbackUrl(rw, req)

		return
	}

	link, err := hdlr.link.GetByToken(token)
	if err != nil {
		hdlr.redirectToFallbackUrl(rw, req)
	} else {
		redirect(rw, req, link.Location)
	}
}

func (hdlr LinkHandler) redirectToFallbackUrl(rw http.ResponseWriter, req *http.Request) {
	http.Redirect(rw, req, hdlr.fallbackUrl, http.StatusFound)
}
