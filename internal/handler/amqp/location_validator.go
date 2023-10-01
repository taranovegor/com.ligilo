package amqp

import (
	kontrakto "github.com/taranovegor/com.kontrakto"
	"github.com/taranovegor/com.ligilo/internal/service"
	amqp "github.com/taranovegor/pkg.amqp"
)

type LocationValidatorConsumer struct {
	amqp.Consumer
	validator *service.Validator
}

func NewLocationValidatorConsumer(
	validator *service.Validator,
) *LocationValidatorConsumer {
	return &LocationValidatorConsumer{
		validator: validator,
	}
}

func (c LocationValidatorConsumer) Name() string {
	return kontrakto.AmqpShortLinkValidateLocation
}

func (c LocationValidatorConsumer) Handle(body amqp.Body) amqp.Handled {
	contract := kontrakto.ValidateShortLinkLocation{}
	body.To(&contract)

	return amqp.HandledSuccessfully().WithReply(c.validator.ValidateLocation(contract.Location))
}
