package amqp

import (
	kontrakto "github.com/taranovegor/com.kontrakto"
	"github.com/taranovegor/com.ligilo/internal/domain"
	"github.com/taranovegor/com.ligilo/internal/service"
	"github.com/taranovegor/com.ligilo/internal/transformer"
	amqp "github.com/taranovegor/pkg.amqp"
)

type LinkUpdaterConsumer struct {
	amqp.Consumer
	repository  domain.LinkRepository
	validator   *service.Validator
	transformer *transformer.LinkQueueTransformer
}

func NewLinkUpdaterConsumer(
	repository domain.LinkRepository,
	validator *service.Validator,
	transformer *transformer.LinkQueueTransformer,
) *LinkUpdaterConsumer {
	return &LinkUpdaterConsumer{
		repository:  repository,
		validator:   validator,
		transformer: transformer,
	}
}

func (c LinkUpdaterConsumer) Name() string {
	return kontrakto.AmqpShortLinkWriteUpdate
}

func (c LinkUpdaterConsumer) Handle(body amqp.Body) amqp.Handled {
	contract := kontrakto.UpdateShortLink{}
	body.To(&contract)

	link, err := c.repository.GetByToken(contract.WithToken)
	if err != nil {
		return amqp.HandledNotSuccessfully(false)
	}

	link.Location = contract.Location
	err = c.repository.Update(&link)
	if err != nil {
		return amqp.HandledNotSuccessfully(false)
	}

	return amqp.HandledSuccessfully().WithReply(c.transformer.Transform(link))
}
