package amqp

import (
	kontrakto "github.com/taranovegor/com.kontrakto"
	"github.com/taranovegor/com.ligilo/internal/domain"
	"github.com/taranovegor/com.ligilo/internal/transformer"
	amqp "github.com/taranovegor/pkg.amqp"
)

type LinkPaginationConsumer struct {
	amqp.Consumer
	repository  domain.LinkRepository
	transformer *transformer.LinkQueueTransformer
}

func NewLinkPaginationConsumer(
	repository domain.LinkRepository,
	transformer *transformer.LinkQueueTransformer,
) *LinkPaginationConsumer {
	return &LinkPaginationConsumer{
		repository:  repository,
		transformer: transformer,
	}
}

func (c LinkPaginationConsumer) Name() string {
	return kontrakto.AmqpShortLinkReadPaginate
}

func (c LinkPaginationConsumer) Handle(body amqp.Body) amqp.Handled {
	contract := kontrakto.PaginateShortLinks{}
	body.To(&contract)

	links, err := c.repository.Paginate(contract.Paginator)
	if err != nil {
		return amqp.HandledNotSuccessfully(false)
	}

	var pagination kontrakto.ShortLinksPagination
	for _, link := range links {
		pagination.Items = append(pagination.Items, c.transformer.Transform(link))
	}

	return amqp.HandledSuccessfully().WithReply(pagination)
}
