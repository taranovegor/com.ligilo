package amqp

import (
	kontrakto "github.com/taranovegor/com.kontrakto"
	"github.com/taranovegor/com.ligilo/internal/domain"
	"github.com/taranovegor/com.ligilo/internal/transformer"
	amqp "github.com/taranovegor/pkg.amqp"
)

type LinkGetterConsumer struct {
	amqp.Consumer
	repository  domain.LinkRepository
	transformer *transformer.LinkQueueTransformer
}

func NewLinkGetterConsumer(
	repository domain.LinkRepository,
	transformer *transformer.LinkQueueTransformer,
) *LinkGetterConsumer {
	return &LinkGetterConsumer{
		repository:  repository,
		transformer: transformer,
	}
}

func (c LinkGetterConsumer) Name() string {
	return kontrakto.AmqpShortLinkReadGet
}

func (c LinkGetterConsumer) Handle(body amqp.Body) amqp.Handled {
	contract := kontrakto.GetShortLink{}
	body.To(&contract)

	link, err := c.repository.GetByToken(contract.WithToken)
	if err != nil {
		return amqp.HandledNotSuccessfully(false)
	}

	return amqp.HandledSuccessfully().WithReply(c.transformer.Transform(link))
}
