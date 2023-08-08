package amqp

import (
	kontrakto "github.com/taranovegor/com.kontrakto"
	"github.com/taranovegor/com.ligilo/internal/domain"
	"github.com/taranovegor/com.ligilo/internal/service"
	"github.com/taranovegor/com.ligilo/internal/transformer"
	amqp "github.com/taranovegor/pkg.amqp"
)

type LinkCreatorConsumer struct {
	amqp.Consumer
	repository  domain.LinkRepository
	validator   *service.Validator
	transformer *transformer.LinkQueueTransformer
}

func NewLinkCreatorConsumer(
	repository domain.LinkRepository,
	validator *service.Validator,
	transformer *transformer.LinkQueueTransformer,
) *LinkCreatorConsumer {
	return &LinkCreatorConsumer{
		repository:  repository,
		validator:   validator,
		transformer: transformer,
	}
}

func (c LinkCreatorConsumer) Name() string {
	return kontrakto.AmqpShortLinkWriteCreate
}

func (c LinkCreatorConsumer) Handle(body amqp.Body) amqp.Handled {
	contract := kontrakto.CreateShortLink{}
	body.To(&contract)

	var result kontrakto.CreateShortLinkResult
	var validation kontrakto.ValidationResult
	validation = c.validator.ValidateToken(contract.Token)
	result.TokenValidation = validation
	c.validator.ValidateLocation(contract.Location)
	result.LocationValidation = validation

	result.Success = result.TokenValidation.Success && result.LocationValidation.Success

	if result.Success {
		link := domain.NewLink(contract.Token, contract.Location)
		err := c.repository.Store(link)
		if err != nil {
			result.Success = false
			result.Message = err.Error()

			return amqp.HandledNotSuccessfully(false).WithReply(result)
		}

		result.Success = true
		result.ShortLink = c.transformer.Transform(*link)
	}

	return amqp.HandledSuccessfully().WithReply(result)
}
