package impl

import (
	"context"
	"github.com/go-pg/pg/v10"
	"main/logging"
	"main/model/entity"
	"main/repository"
	"main/utils"
)

func NewPgOrmQuotationSessionRepository(pgOrm *pg.DB) repository.QuotationSessionRepository {
	return quotationSessionRepository{pgOrm: pgOrm}
}

type quotationSessionRepository struct {
	pgOrm *pg.DB
}

func (q quotationSessionRepository) GetAllSessions(ctx context.Context) ([]*entity.QuotationSession, error) {
	var res []*entity.QuotationSession
	err := utils.RunWithProfiler(repository.TagGetAllQS, func() error {
		tx, err := q.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get All Sessions transaction: %s", err)
			return err
		}
		defer tx.Rollback()
		err = tx.Model(&res).
			Relation("Creator").
			Relation("LastBet").
			Relation("Products").Select()
		if err != nil {
			logging.ErrorFormat("Cannot Get All sessions : %s",
				err.Error())
			return err
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return res, err
	}
	return res, nil
}

func (q quotationSessionRepository) GetSessionById(ctx context.Context, sessionId int64) (entity.QuotationSession, error) {
	var res entity.QuotationSession
	err := utils.RunWithProfiler(repository.TagGetQSByStatus, func() error {
		tx, err := q.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get By Status transaction: %s", err)
			return err
		}
		defer tx.Rollback()
		err = tx.Model(&res).
			Relation("Creator").
			Relation("LastBet").
			Relation("Products").Where("quotation_session.id = ?", sessionId).Select()
		if err != nil {
			logging.ErrorFormat("Cannot Get session by id %s: %s", sessionId,
				err.Error())
			return err
		}

		pr := entity.Product{}
		for i, pJ := range res.Products {
			err = tx.Model(&pr).Where("id = ?", pJ.ProductId).Select()
			if pr.Id != 0 {
				res.Products[i].Product = &pr
			}
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return res, err
	}
	return res, nil
}

func (q quotationSessionRepository) GetSessionsByStatus(ctx context.Context, status entity.SessionStatus) ([]*entity.QuotationSession, error) {
	var res []*entity.QuotationSession
	err := utils.RunWithProfiler(repository.TagGetQSByStatus, func() error {
		tx, err := q.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Get By Status transaction: %s", err)
			return err
		}
		defer tx.Rollback()
		err = tx.Model(&res).
			Relation("Creator").
			Relation("LastBet").Where("status = ?", status).Select()
		if err != nil {
			logging.ErrorFormat("Cannot Get session by status %s: %s", status,
				err.Error())
			return err
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return res, err
	}
	return res, nil
}

func (q quotationSessionRepository) DeleteQuotationSession(ctx context.Context, quotationSession entity.QuotationSession) error {
	err := utils.RunWithProfiler(repository.TagDelQS, func() error {
		tx, err := q.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Insert transaction: %s", err)
			return err
		}
		defer tx.Rollback()
		_, err = tx.Model(&quotationSession).
			WherePK().Delete()
		if err != nil {
			logging.ErrorFormat("Cannot Delete session %d: %s", quotationSession.ID,
				err.Error())
			return err
		}
		_, err = tx.Model(&entity.Bet{}).
			Where("quotation_session = ?", quotationSession.ID).Delete()
		if err != nil {
			logging.ErrorFormat("Cannot Delete bets for session %d: %s", quotationSession.ID,
				err.Error())
			return err
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (q quotationSessionRepository) NewQuotationSession(ctx context.Context, quotationSession entity.QuotationSession) (int64, error) {
	err := utils.RunWithProfiler(repository.TagInsQS, func() error {
		tx, err := q.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Insert transaction: %s", err)
			return err
		}
		defer tx.Rollback()
		_, err = tx.Model(&quotationSession).
			Returning("id").Insert()
		if err != nil {
			logging.ErrorFormat("Cannot Insert session %d: %s", quotationSession.ID,
				err.Error())
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
	return quotationSession.ID, nil
}

func (q quotationSessionRepository) UpdateQuotationSession(ctx context.Context, quotationSession entity.QuotationSession) error {
	err := utils.RunWithProfiler(repository.TagUpdQS, func() error {
		tx, err := q.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Update transaction: %s", err)
			return err
		}
		defer tx.Rollback()
		_, err = tx.Model(&quotationSession).WherePK().Update()
		if err != nil {
			logging.ErrorFormat("Cannot Update session %d: %s", quotationSession.ID,
				err.Error())
			return err
		}

		if err = tx.Commit(); err != nil {
			logging.Error("could not commit a transaction")
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil

}
