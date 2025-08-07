package repository

import (
	"github.com/noxhalley/funken/internal/infrastructure/log"
	"github.com/noxhalley/funken/internal/infrastructure/mongodb"
)

type MessageRepository interface {
}

type messageRepo struct {
	logger *log.Logger
	db     *mongodb.MongoDB
}

func NewMessageRepository(db *mongodb.MongoDB) MessageRepository {
	return &messageRepo{
		logger: log.With("repository", "message_repository"),
		db:     db,
	}
}
