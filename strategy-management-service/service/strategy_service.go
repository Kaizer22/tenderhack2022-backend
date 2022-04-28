package service

import (
	"context"
	"fmt"
	"main/logging"
	"main/model/entity"
	"main/repository"
	"main/utils"
	"time"
)

func NewStrategyService(ctx context.Context,
	qsRepo repository.QuotationSessionRepository,
	bRepo repository.BetRepository) StrategyService {
	return StrategyService{
		quotationSessionRepo: qsRepo,
		betRepo:              bRepo,
		ctx:                  ctx,
	}
}

type StrategyJob struct {
	entity.CurrentSessionState
	QuotationSessionId int64
	S                  entity.Strategy
	QuitChannel        chan bool
	ParentService      *StrategyService
}

func (j StrategyJob) Run(qsRepo repository.QuotationSessionRepository, betRepo repository.BetRepository) error {
	go func() {
		counter := 0
		for {
			select {
			case <-j.QuitChannel:
				return
			default:
				{
					err := j.recalculateCurrentState(qsRepo, betRepo)
					if err != nil {
						logging.ErrorFormat("Cannot recalculate current state params "+
							"for runner session-%d-user-%d: %s", j.QuotationSessionId, j.UserId, err)
					}
					action := j.decide(j.S.BaseConditionSet)
					err = j.perform(action)
					if err != nil {
						logging.ErrorFormat("Cannot perform action for runner session-%d-user-%d: %s",
							j.QuotationSessionId, j.UserId, err)
					}

					logging.InfoFormat("Running %d; runner %d %d", counter,
						j.UserId, j.QuotationSessionId)

					time.Sleep(time.Duration(j.S.N) * time.Second)
				}
			}

		}
	}()
	return nil
}

func (j *StrategyJob) recalculateCurrentState(qsRepo repository.QuotationSessionRepository,
	betRepo repository.BetRepository) error {
	ctx := context.Background()
	session, err := qsRepo.GetSessionById(ctx, j.QuotationSessionId)
	if err != nil {
		logging.ErrorFormat("Cannot get session to recalculate current state runner session-%d-user-%d",
			j.QuotationSessionId, j.UserId)
		return err
	}
	j.IsOnAdditionalPurchase = session.IsInAdditionalPurchase
	j.CurrentWinnerId = session.LastBet.ProviderId
	j.MyCurrentBetNumber, err = betRepo.CountBetsBySessionAndUser(ctx, session.ID, j.UserId)
	if err != nil {
		logging.ErrorFormat("Cannot calculate MyCurrentBetNumber for runner session-%d-user-%d", j.QuotationSessionId, j.UserId)
	}
	j.CurrentStepNumber = int64(session.LastBet.BetNumber)
	j.StepsTillZero = int64(100/session.SessionStepPercent) - j.CurrentStepNumber
	j.CurrentPrice = session.CurrentPrice
	j.CurrentDiscount = session.StartPrice - session.CurrentPrice
	j.TimeSinceLastStep = time.Now().Sub(session.LastBet.Time)
	j.TimeSinceLastMyBet, err = betRepo.GetDeltaBetTimeForUser(ctx, session.ID, j.UserId)
	if err != nil {
		logging.ErrorFormat("Cannot calculate TimeSinceLastMyBet for runner session-%d-user-%d", j.QuotationSessionId, j.UserId)
	}
	j.StepSize = session.SessionStepPercent / 100 * session.StartPrice
	j.TimeSinceStart = time.Now().Sub(session.StartTime)
	j.TimeTillEnd = session.StartTime.Add(
		time.Duration(session.SessionDuration) * time.Minute).Sub(time.Now())
	j.ParticipantsCount, err = betRepo.CountSessionParticipants(ctx, session.ID)
	if err != nil {
		logging.ErrorFormat("Cannot calculate ParticipantsCount for runner session-%d-user-%d", j.QuotationSessionId, j.UserId)
	}

	logging.DebugFormat("Current State: %v+", j.CurrentSessionState)

	return nil
}

func (j *StrategyJob) decide(set *entity.ConditionSet) entity.Action {
	return set.Define(j.CurrentSessionState)
}

func (j *StrategyJob) perform(action entity.Action) error {
	logging.InfoFormat("Runner session-%d-user-%d: Performing action %s", j.QuotationSessionId,
		j.UserId, action)
	switch action {
	case entity.ActionBet:
		err := j.ParentService.MakeBet(j.QuotationSessionId, j.UserId)
		if err != nil {
			return err
		}
	}
	return nil
}

var runnersPool []*StrategyJob

type StrategyService struct {
	quotationSessionRepo repository.QuotationSessionRepository
	betRepo              repository.BetRepository
	ctx                  context.Context
}

func (s *StrategyService) RunStrategyRunner(params entity.StrategyParams) error {
	strat := entity.Strategy{}
	switch params.Str {
	case "aggressive":
		strat = utils.AggressiveStrategy
	case "waiting":
		strat = utils.WaitingStrategy
	case "progressive":
		strat = utils.ProgressiveStrategy
	}
	job := StrategyJob{
		CurrentSessionState: entity.CurrentSessionState{
			UserId:          params.UserId,
			MinimalPrice:    params.MinimalPrice,
			AcceptablePrice: params.AcceptablePrice,
			PreferablePrice: params.PreferablePrice,
		},
		QuotationSessionId: params.QuotationSessionId,
		S:                  strat,
		QuitChannel:        make(chan bool),
		ParentService:      s,
	}
	if sJ, i := s.findBySessionIdAndUserId(params.QuotationSessionId, params.UserId); sJ != nil || i != -1 {
		return fmt.Errorf("runner session-%d-user-%d is already running",
			params.QuotationSessionId, params.UserId)
	}
	runnersPool = append(runnersPool, &job)
	err := s.MakeBet(params.QuotationSessionId, params.UserId)
	if err != nil {
		logging.ErrorFormat("Cannot make initial bet to run the strategy %s", err)
		return err
	}
	err = job.Run(s.quotationSessionRepo, s.betRepo)
	if err != nil {
		return err
	}
	return nil
}

func (s *StrategyService) StopRunner(sessionId int64, userId int64) error {
	job, index := s.findBySessionIdAndUserId(sessionId, userId)
	if job != nil {
		job.QuitChannel <- true
		logging.InfoFormat("Quitting runner for session %d : user %d", sessionId, userId)
		runnersPool = append(runnersPool[:index], runnersPool[index+1:]...)
		return nil
	} else {
		return fmt.Errorf("runner for session - %d, user - %d not found", sessionId, userId)
	}

}

func (s *StrategyService) findBySessionIdAndUserId(sessionId int64, userId int64) (*StrategyJob, int) {
	for i, job := range runnersPool {
		if job.UserId == userId && job.QuotationSessionId == sessionId {
			return job, i
		}
	}
	return nil, -1
}

func (s *StrategyService) MakeBet(sessionId int64, profileId int64) error {
	bet, err := s.betRepo.MakeBet(s.ctx, entity.BetData{
		QuotationSessionID: sessionId,
		ProviderId:         profileId,
		Bot:                true,
	})
	if err != nil {
		return err
	}
	logging.InfoFormat("Bet %d has made", bet)
	return nil
}
