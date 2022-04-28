package repository

import (
	"context"
	"main/model/entity"
	"time"
)

const (
	TagMakeBet                   = "MAKE BET"
	TagGetBetById                = "GET BET BY ID"
	TagCountBetsBySessionAndUser = "COUNT BETS BY SESSION AND USER"
	TagGetDeltaBetTimeForUser    = "GET DELTA BET TIME FOR USER"
	TagCountSessionParticipants  = "COUNT SESSION PARTICIPANTS"
)

type BetRepository interface {
	MakeBet(context.Context, entity.BetData) (int64, error)
	GetBetBySessionId(ctx context.Context, sessionId int64) ([]*entity.Bet, error)
	CountBetsBySessionAndUser(ctx context.Context, sessionId int64, profileID int64) (int64, error)
	GetDeltaBetTimeForUser(ctx context.Context, sessionId int64, profileID int64) (time.Duration, error)
	CountSessionParticipants(ctx context.Context, sessionId int64) (int64, error)
}
