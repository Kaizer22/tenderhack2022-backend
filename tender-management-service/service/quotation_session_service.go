package service

import (
	"context"
	"github.com/robfig/cron/v3"
	"main/logging"
	"main/model/entity"
	"main/repository"
	"main/utils"
	"time"
)

func NewQuotationSessionService(ctx context.Context,
	qsRepo repository.QuotationSessionRepository) QuotationSessionService {
	return QuotationSessionService{
		quotationSessionRepo: qsRepo,
		ctx:                  ctx,
	}
}

type QuotationSessionService struct {
	quotationSessionRepo repository.QuotationSessionRepository
	ctx                  context.Context
}

func (s QuotationSessionService) RunQuotationService() error {
	if utils.GetEnvBool(utils.RunSessionCron, true) {
		scheduler := cron.New()
		schedule := utils.GetEnv(utils.SessionUpdateFrequencyEnvKey, "*/1 * * * *")
		logging.InfoFormat("Starting quotation service with schedule %s", schedule)
		_, err := scheduler.AddFunc(schedule,
			s.updateQuotationSessionsStatus)
		if err != nil {
			return err
		}
		scheduler.Start()
	}
	return nil
}

func (s *QuotationSessionService) updateQuotationSessionsStatus() {
	logging.InfoFormat("Updating sessions statuses...")
	if sessions, err := s.quotationSessionRepo.GetSessionsByStatus(s.ctx, entity.StatusActive); err == nil {
		logging.InfoFormat("Found %d active sessions", len(sessions))
		for _, session := range sessions {
			currentTime := time.Now()
			endTime := session.StartTime.Add(time.Duration(session.SessionDuration) * time.Minute)
			if endTime.Before(currentTime) {
				//if session.IsInAdditionalPurchase {
				//	if session.LastBet != nil && session.LastBet.Time.After(
				//		currentTime.Add(5*time.Minute)) {
				//		logging.InfoFormat("Finishing session %d", session.ID)
				//		session.Status = entity.StatusFinished
				//	}
				//	err := s.quotationSessionRepo.UpdateQuotationSession(s.ctx, *session)
				//	if err != nil {
				//		logging.ErrorFormat("Cannot update session %d", session.ID)
				//		continue
				//	}
				//} else {
				if session.LastBet != nil {
					logging.InfoFormat("Finishing session %d", session.ID)
					session.Status = entity.StatusFinished
				} else {
					logging.InfoFormat("Session %d didnt take place", session.ID)
					session.Status = entity.StatusDidntTakePlace
				}
				err := s.quotationSessionRepo.UpdateQuotationSession(s.ctx, *session)
				if err != nil {
					logging.ErrorFormat("Cannot update session %d", session.ID)
					continue
				}
			}
		}
	}
}
