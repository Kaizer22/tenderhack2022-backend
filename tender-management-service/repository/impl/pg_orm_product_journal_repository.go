package impl

import (
	"context"
	"github.com/go-pg/pg/v10"
	"main/logging"
	"main/model/entity"
	"main/repository"
	"main/utils"
)

func NewPgOrmProductJournalRepository(ctx context.Context,
	db *pg.DB) repository.ProductJournalRepository {
	return pgOrmProductJournalRepository{
		pgOrm: db,
	}
}

type pgOrmProductJournalRepository struct {
	pgOrm *pg.DB
}

func (p pgOrmProductJournalRepository) AddProductJournal(ctx context.Context, sessionId int64, data []*entity.ProductJournalData) error {

	err := utils.RunWithProfiler(repository.TagAddPrJ, func() error {
		tx, err := p.pgOrm.Begin()
		if err != nil {
			logging.ErrorFormat("Cannot open Insert transaction: %s", err)
			return err
		}
		defer tx.Rollback()
		for _, pJData := range data {
			//TODO handle returned id
			productJournal := entity.ProductJournal{
				RecordID:           0,
				ProductId:          pJData.ProductId,
				Product:            nil,
				QuotationSessionId: sessionId,
				Count:              pJData.Count,
			}

			_, err = tx.Model(&productJournal).Returning("id").Insert()
			if err != nil {
				logging.ErrorFormat("Cannot Insert product journal %d: %s", productJournal.RecordID,
					err.Error())
				return err
			}
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
