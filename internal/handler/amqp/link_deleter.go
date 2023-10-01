package amqp

import (
	kontrakto "github.com/taranovegor/com.kontrakto"
	"github.com/taranovegor/com.ligilo/internal/domain"
	amqp "github.com/taranovegor/pkg.amqp"
)

type LinkDeleterConsumer struct {
	amqp.Consumer
	repository domain.LinkRepository
}

func NewLinkDeleterConsumer(
	repository domain.LinkRepository,
) *LinkDeleterConsumer {
	return &LinkDeleterConsumer{
		repository: repository,
	}
}

func (c LinkDeleterConsumer) Name() string {
	return kontrakto.AmqpShortLinkWriteDelete
}

func (c LinkDeleterConsumer) Handle(body amqp.Body) amqp.Handled {
	contract := kontrakto.DeleteShortLink{}
	body.To(&contract)

	err := c.repository.DeleteByToken(contract.WithToken)
	if err != nil {
		return amqp.HandledNotSuccessfully(false)
	}

	return amqp.HandledSuccessfully()
}
