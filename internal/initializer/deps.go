package initializer

import (
	"context"
	"io"
	"os"

	"github.com/noxhalley/funken/config"
	"github.com/noxhalley/funken/internal/infrastructure/log"
	"github.com/noxhalley/funken/internal/infrastructure/mongodb"
	"github.com/noxhalley/funken/internal/infrastructure/pubsub"
	"github.com/noxhalley/funken/internal/infrastructure/repository"

	"go.uber.org/fx"
)

type initParams struct {
	fx.In
	writer io.Writer
	cfg    *config.Config
	keys   []string
}

func initLog(p initParams) {
	log.Initialize(p.writer, p.cfg, p.keys)
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			config.NewConfig,
			func() io.Writer { return os.Stdout },
			func() []string { return []string{} },
		),
		fx.Invoke(initLog),
		fx.Provide(mongo),
		fx.Provide(
			fx.Annotate(
				jetstreamManager,
				fx.As(
					new(pubsub.Publisher),
					new(pubsub.Subcriber),
					new(pubsub.PubSub),
					new(pubsub.StreamConsumerManager),
					new(pubsub.PubSubStreamManager),
				),
			),
		),

		// repositories
		fx.Provide(repository.NewGroupRepository),
		fx.Provide(repository.NewMemberGroupRepository),
		fx.Provide(repository.NewGroupNGFilterRepository),
		fx.Provide(repository.NewMessageRepository),
	)
}

func mongo(lc fx.Lifecycle, cfg *config.Config) *mongodb.MongoDB {
	mdb := mongodb.NewOrGetSingleton(cfg)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return mdb.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			mdb.Close(ctx)
			return nil
		},
	})
	return mdb
}

func jetstreamManager(lc fx.Lifecycle, cfg *config.Config) *pubsub.JetStreamManager {
	jsm := pubsub.NewOrGetSingleton(cfg)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			jsm.Close()
			return nil
		},
	})
	return jsm
}
