package amqp

import (
	kontrakto "github.com/taranovegor/com.kontrakto"
	"github.com/taranovegor/com.ligilo/internal/service"
	amqp "github.com/taranovegor/pkg.amqp"
)

type TokenValidatorConsumer struct {
	amqp.Consumer
	validator *service.Validator
}

func NewTokenValidatorConsumer(
	validator *service.Validator,
) *TokenValidatorConsumer {
	return &TokenValidatorConsumer{
		validator: validator,
	}
}

func (c TokenValidatorConsumer) Name() string {
	return kontrakto.AmqpShortLinkValidateToken
}

func (c TokenValidatorConsumer) Handle(body amqp.Body) amqp.Handled {
	contract := kontrakto.ValidateShortLinkToken{}
	body.To(&contract)

	return amqp.HandledSuccessfully().WithReply(c.validator.ValidateToken(contract.Token))
}
