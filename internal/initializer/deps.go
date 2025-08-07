package initializer

import (
	"context"

	"github.com/noxhalley/funken/config"
	"github.com/noxhalley/funken/internal/infrastructure/mongodb"
	"github.com/noxhalley/funken/internal/infrastructure/repository"
	"go.uber.org/fx"
)

func Build() fx.Option {
	return fx.Options(
		fx.Provide(config.NewConfig),
		fx.Provide(mongo),

		// repositories
		fx.Provide(repository.NewGroupRepository),
        fx.Provide(repository.NewMemberGroupRepository),
		fx.Provide(repository.NewMessageRepository),
	)
}

func mongo(lc fx.Lifecycle, cfg *config.Config) (*mongodb.MongoDB, error) {
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
	return mdb, nil
}
