package repository

import (
	"context"
	"main/model/entity"
)

const (
	TagAddPrJ = "ADD PRODUCT JOURNAL"
)

type ProductJournalRepository interface {
	AddProductJournal(ctx context.Context, sessionId int64, data []*entity.ProductJournalData) error
}
