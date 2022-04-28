package repository

import (
	"context"
	"main/model/entity"
)

const (
	TagInsQS         = "INSERT SESSION"
	TagUpdQS         = "UPDATE SESSION"
	TagDelQS         = "DELETE SESSION"
	TagGetQSByStatus = "GET SESSION BY STATUS"
)

type QuotationSessionRepository interface {
	NewQuotationSession(ctx context.Context, quotationSession entity.QuotationSession) (int64, error)
	GetSessionsByStatus(ctx context.Context, status entity.SessionStatus) ([]*entity.QuotationSession, error)
	GetSessionById(ctx context.Context, sessionId int64) (entity.QuotationSession, error)
	UpdateQuotationSession(ctx context.Context, quotationSession entity.QuotationSession) error
	DeleteQuotationSession(ctx context.Context, quotationSession entity.QuotationSession) error
}
