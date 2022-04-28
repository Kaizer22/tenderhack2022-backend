package impl

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"main/logging"
	"main/model/entity"
	"main/repository"
	"main/utils"
	"time"
)

func NewPgOrmBetRepository(ctx context.Context,
	db *pg.DB) repository.BetRepository {
	return pgOrmBetRepository{
		pgOrm: db,
	}
}

type pgOrmBetRepository struct {
	pgOrm *pg.DB
}

func (p pgOrmBetRepository) CountBetsBySessionAndUser(ctx context.Context, sessionId int64, profileID int64) (int64, error) {
	var res int
	err := utils.RunWithProfiler(repository.TagCountBetsBySessionAndUser, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Count bets by session and user transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		res, err = tx.Model(&entity.Bet{}).
			Where("quotation_session = ?", sessionId).
			Where("provider_id = ?", profileID).
			Count()
		if err != nil {
			logging.ErrorFormat("Error counting bets by session id and profile id: %s", err)
			return err
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return -1, err
	}
	return int64(res), nil
}

func (p pgOrmBetRepository) GetDeltaBetTimeForUser(ctx context.Context, sessionId int64, profileID int64) (time.Duration, error) {
	var res time.Duration
	err := utils.RunWithProfiler(repository.TagGetDeltaBetTimeForUser, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Count bets by session and user transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		buf := entity.Bet{}
		err = tx.Model(&buf).
			Where("quotation_session = ?", sessionId).
			Where("provider_id = ?", profileID).
			Order("time DESC").First()
		if err != nil {
			logging.ErrorFormat("Error counting bets by session id and profile id: %s", err)
			return err
		}

		res = time.Now().Sub(buf.Time)

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (p pgOrmBetRepository) CountSessionParticipants(ctx context.Context, sessionId int64) (int64, error) {
	var res int
	err := utils.RunWithProfiler(repository.TagCountSessionParticipants, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Count session participants transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		res, err = tx.Model(&entity.Bet{}).Where("quotation_session = ?", sessionId).
			Column("provider_id").Distinct().Count()
		if err != nil {
			logging.ErrorFormat("Error counting session particiapnts: %s", err)
			return err
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return -1, err
	}
	return int64(res), nil
}

func (p pgOrmBetRepository) GetBetBySessionId(ctx context.Context, sessionId int64) ([]*entity.Bet, error) {
	var res []*entity.Bet
	err := utils.RunWithProfiler(repository.TagGetBetById, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get bets By session Id transaction: %s", err)
			return err
		}
		defer tx.Rollback()

		err = tx.Model(&res).Where("quotation_session = ?0", sessionId).Select()
		if err != nil {
			logging.ErrorFormat("Error selecting bets by session id: %s", err)
			return err
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p pgOrmBetRepository) MakeBet(ctx context.Context, data entity.BetData) (int64, error) {
	bet := entity.Bet{
		ID:                 0,
		QuotationSessionID: data.QuotationSessionID,
		ProviderId:         data.ProviderId,
		BetNumber:          -1,
		Time:               time.Now(),
		Bot:                data.Bot,
	}
	err := utils.RunWithProfiler(repository.TagMakeBet, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Make Bet transaction: %s", err)
			return err
		}
		defer tx.Rollback()
		session := entity.QuotationSession{}
		err = tx.Model(&session).Relation("Creator").
			Relation("LastBet").
			Relation("Products").
			Where("quotation_session.id = ?", data.QuotationSessionID).Select()
		if err != nil {
			logging.ErrorFormat("Cannot Get session by ID  %d: %s", data.QuotationSessionID,
				err.Error())
			return err
		}
		if session.Status == entity.StatusActive && !session.IsInAdditionalPurchase {
			newPrice := session.CurrentPrice - session.StartPrice*(session.SessionStepPercent/100)
			if newPrice <= 0 {
				session.CurrentPrice = 0
				bet.NewPrice = 0
				session.Status = entity.StatusFinished
			} else {
				session.CurrentPrice = newPrice
				bet.NewPrice = newPrice
			}

			if session.LastBet != nil {
				bet.BetNumber = session.LastBet.BetNumber + 1
				if session.LastBet.ProviderId == bet.ProviderId {
					return fmt.Errorf("cannot make bet: this provider made the last bet")
				}
			} else {
				bet.BetNumber = 0
			}
			_, err = tx.Model(&bet).Returning("id").Insert()
			if err != nil {
				logging.ErrorFormat("Cannot Insert new bet %v+: %s", bet,
					err.Error())
				return err
			}
			session.LastBetId = bet.ID
			session.LastBet = &bet
			endTime := session.StartTime.Add(time.Duration(session.SessionDuration) * time.Minute)
			if session.Status == entity.StatusActive &&
				endTime.Sub(bet.Time) > (0*time.Millisecond) && endTime.Sub(bet.Time) < (5*time.Minute) {
				session.IsInAdditionalPurchase = true
				session.SessionDuration += 5
			}
		} else if session.Status == entity.StatusActive && session.IsInAdditionalPurchase {
			endTime := session.StartTime.Add(time.Duration(session.SessionDuration) * time.Minute)
			if endTime.Sub(bet.Time) > (0 * time.Millisecond) {
				newPrice := session.CurrentPrice - session.StartPrice*(session.SessionStepPercent/100)
				if newPrice <= 0 {
					session.CurrentPrice = 0
					bet.NewPrice = 0
					session.Status = entity.StatusFinished
				} else {
					session.CurrentPrice = newPrice
					bet.NewPrice = newPrice
				}

				if session.LastBet != nil {
					bet.BetNumber = session.LastBet.BetNumber + 1
					if session.LastBet.ProviderId == bet.ProviderId {
						return fmt.Errorf("cannot make bet: this provider made the last bet")
					}
				} else {
					bet.BetNumber = 0
				}
				_, err = tx.Model(&bet).Returning("id").Insert()
				if err != nil {
					logging.ErrorFormat("Cannot Insert new bet %v+: %s", bet,
						err.Error())
					return err
				}
				session.LastBetId = bet.ID
				session.LastBet = &bet
				if endTime.Sub(bet.Time) <= (5 * time.Minute) {
					session.SessionDuration += 5
				}
			} else {
				session.Status = entity.StatusFinished
			}

		} else {
			return fmt.Errorf("cannot make bet: session isn't active")
		}

		_, err = tx.Model(&session).WherePK().Update()
		if err != nil {
			logging.ErrorFormat("Cannot Update session during Bet transaction: %s", err)
			return err
		}
		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return -1, err
	}
	return bet.ID, nil
}
