package transformer

import (
	kontrakto "github.com/taranovegor/com.kontrakto"
	"github.com/taranovegor/com.ligilo/internal/domain"
	"github.com/taranovegor/com.ligilo/internal/service"
)

type LinkQueueTransformer struct {
	urlProvider *service.UrlProvider
}

func NewLinkQueueTransformer(
	urlProvider *service.UrlProvider,
) *LinkQueueTransformer {
	return &LinkQueueTransformer{
		urlProvider: urlProvider,
	}
}

func (t LinkQueueTransformer) Transform(link domain.Link) kontrakto.ShortLink {
	return kontrakto.ShortLink{
		Token:    link.Token,
		Location: link.Location,
		Url:      t.urlProvider.Provide(link),
	}
}
