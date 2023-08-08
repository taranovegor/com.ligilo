package service

import (
	"fmt"
	"github.com/taranovegor/com.ligilo/internal/domain"
	"strings"
)

type UrlProvider struct {
	baseUrl string
}

func NewUrlProvider(
	baseUrl string,
) *UrlProvider {
	return &UrlProvider{
		baseUrl: strings.TrimRight(baseUrl, "/"),
	}
}

func (p UrlProvider) Provide(l domain.Link) string {
	return fmt.Sprintf("%s/%s", p.baseUrl, l.Token)
}
