package container

import (
	"github.com/go-chi/chi/v5"
	"github.com/sarulabs/di"
	kontrakto "github.com/taranovegor/com.kontrakto"
	"github.com/taranovegor/com.ligilo/internal/config"
	"github.com/taranovegor/com.ligilo/internal/domain"
	consumer "github.com/taranovegor/com.ligilo/internal/handler/amqp"
	"github.com/taranovegor/com.ligilo/internal/handler/http"
	"github.com/taranovegor/com.ligilo/internal/repository"
	"github.com/taranovegor/com.ligilo/internal/service"
	"github.com/taranovegor/com.ligilo/internal/transformer"
	amqp "github.com/taranovegor/pkg.amqp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	AmqpController = "amqp_controller"
	HttpRouter     = "http_router"
	Orm            = "orm"

	handlerAmqpLinkCreator           = "handler_amqp_link_creator"
	handlerAmqpLinkDeleter           = "handler_amqp_link_deleter"
	handlerAmqpLinkGetter            = "handler_amqp_link_getter"
	handlerAmqpLinkLocationValidator = "handler_amqp_link_location_validator"
	handlerAmqpLinkPagination        = "handler_amqp_link_pagination"
	handlerAmqpLinkTokenValidator    = "handler_amqp_link_token_validator"
	handlerAmqpLinkUpdater           = "handler_amqp_link_updater"
	handlerHttpLink                  = "handler_http_link"
	repositoryLink                   = "repository_link"
	serviceUrlProvider               = "service_url_provider"
	serviceValidator                 = "service_validator"
	transformerLinkQueue             = "transformer_link_queue"
)

type ServiceContainer interface {
	Get(name string) interface{}
}

type serviceContainer struct {
	ServiceContainer
	container di.Container
}

func Init() (ServiceContainer, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return nil, err
	}

	return &serviceContainer{
		container: build(builder),
	}, nil
}

func (sc serviceContainer) Get(name string) interface{} {
	return sc.container.Get(name)
}

func build(builder *di.Builder) di.Container {
	buildRepository(builder)
	buildService(builder)
	buildTransformer(builder)
	buildHandler(builder)

	return builder.Build()
}

func buildRepository(builder *di.Builder) {
	builder.Add(di.Def{
		Name: Orm,
		Build: func(ctx di.Container) (interface{}, error) {
			return gorm.Open(
				mysql.Open(config.GetEnv(config.DatabaseDsn)),
				&gorm.Config{},
			)
		},
	})

	builder.Add(di.Def{
		Name: repositoryLink,
		Build: func(ctx di.Container) (interface{}, error) {
			return repository.NewLinkRepository(
				ctx.Get(Orm).(*gorm.DB),
			), nil
		},
	})
}

func buildService(builder *di.Builder) {
	builder.Add(di.Def{
		Name: serviceUrlProvider,
		Build: func(ctx di.Container) (interface{}, error) {
			return service.NewUrlProvider(config.GetEnv(config.AppUrl)), nil
		},
	})

	builder.Add(di.Def{
		Name: serviceValidator,
		Build: func(ctx di.Container) (interface{}, error) {
			return service.NewValidator(ctx.Get(repositoryLink).(domain.LinkRepository)), nil
		},
	})
}

func buildTransformer(builder *di.Builder) {
	builder.Add(di.Def{
		Name: transformerLinkQueue,
		Build: func(ctx di.Container) (interface{}, error) {
			return transformer.NewLinkQueueTransformer(ctx.Get(serviceUrlProvider).(*service.UrlProvider)), nil
		},
	})
}

func buildHandler(builder *di.Builder) {
	buildHandlerHttp(builder)
	buildHandlerAmqp(builder)
}

func buildHandlerHttp(builder *di.Builder) {
	builder.Add(di.Def{
		Name: HttpRouter,
		Build: func(ctx di.Container) (interface{}, error) {
			router := chi.NewRouter()
			http.MapRouter(
				router,
				config.GetEnv(config.FallbackUrl),
				ctx.Get(handlerHttpLink).(*http.LinkHandler),
			)

			return router, nil
		},
	})

	builder.Add(di.Def{
		Name: handlerHttpLink,
		Build: func(ctx di.Container) (interface{}, error) {
			return http.NewLinkHandler(
				ctx.Get(repositoryLink).(domain.LinkRepository),
			), nil
		},
	})
}

func buildHandlerAmqp(builder *di.Builder) {
	builder.Add(di.Def{
		Name: AmqpController,
		Build: func(ctx di.Container) (interface{}, error) {
			return amqp.Init(
				config.AppName,
				config.GetEnv(config.AmqpDsn),
				kontrakto.AmqpConfig(),
				[]amqp.Consumer{
					ctx.Get(handlerAmqpLinkCreator).(amqp.Consumer),
					ctx.Get(handlerAmqpLinkDeleter).(amqp.Consumer),
					ctx.Get(handlerAmqpLinkGetter).(amqp.Consumer),
					ctx.Get(handlerAmqpLinkPagination).(amqp.Consumer),
					ctx.Get(handlerAmqpLinkUpdater).(amqp.Consumer),
					ctx.Get(handlerAmqpLinkLocationValidator).(amqp.Consumer),
					ctx.Get(handlerAmqpLinkTokenValidator).(amqp.Consumer),
				},
			)
		},
	})

	builder.Add(di.Def{
		Name: handlerAmqpLinkCreator,
		Build: func(ctx di.Container) (interface{}, error) {
			return consumer.NewLinkCreatorConsumer(
				ctx.Get(repositoryLink).(domain.LinkRepository),
				ctx.Get(serviceValidator).(*service.Validator),
				ctx.Get(transformerLinkQueue).(*transformer.LinkQueueTransformer),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: handlerAmqpLinkDeleter,
		Build: func(ctx di.Container) (interface{}, error) {
			return consumer.NewLinkDeleterConsumer(
				ctx.Get(repositoryLink).(domain.LinkRepository),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: handlerAmqpLinkGetter,
		Build: func(ctx di.Container) (interface{}, error) {
			return consumer.NewLinkGetterConsumer(
				ctx.Get(repositoryLink).(domain.LinkRepository),
				ctx.Get(transformerLinkQueue).(*transformer.LinkQueueTransformer),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: handlerAmqpLinkPagination,
		Build: func(ctx di.Container) (interface{}, error) {
			return consumer.NewLinkPaginationConsumer(
				ctx.Get(repositoryLink).(domain.LinkRepository),
				ctx.Get(transformerLinkQueue).(*transformer.LinkQueueTransformer),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: handlerAmqpLinkUpdater,
		Build: func(ctx di.Container) (interface{}, error) {
			return consumer.NewLinkUpdaterConsumer(
				ctx.Get(repositoryLink).(domain.LinkRepository),
				ctx.Get(serviceValidator).(*service.Validator),
				ctx.Get(transformerLinkQueue).(*transformer.LinkQueueTransformer),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: handlerAmqpLinkLocationValidator,
		Build: func(ctx di.Container) (interface{}, error) {
			return consumer.NewLocationValidatorConsumer(
				ctx.Get(serviceValidator).(*service.Validator),
			), nil
		},
	})

	builder.Add(di.Def{
		Name: handlerAmqpLinkTokenValidator,
		Build: func(ctx di.Container) (interface{}, error) {
			return consumer.NewTokenValidatorConsumer(
				ctx.Get(serviceValidator).(*service.Validator),
			), nil
		},
	})
}
