package repository

import (
	"context"
	"main/model/entity"
)

const (
	TagMakeBet    = "MAKE BET"
	TagGetBetById = "GET BET BY ID"
)

type BetRepository interface {
	MakeBet(context.Context, entity.BetData) (int64, error)
	GetBetBySessionId(ctx context.Context, sessionId int64) ([]*entity.Bet, error)
}
