package entity

import (
	"time"
)

// Bet example
type Bet struct {
	ID                 int64     `pg:"id,pk" json:"id"`
	QuotationSessionID int64     `pg:"quotation_session" json:"quotation_session_id"`
	ProviderId         int64     `pg:"provider_id" json:"provider_id"`
	BetNumber          int       `pg:"bet_number,use_zero" json:"bet_number"`
	Time               time.Time `pg:"time" json:"time"`
	Bot                bool      `pg:"bot,use_zero" json:"bot"`
	NewPrice           float64   `pg:"new_price" json:"new_price"`
}

// BetData example
type BetData struct {
	QuotationSessionID int64 `json:"quotation_session_id"`
	ProviderId         int64 `json:"provider_id"`
	Bot                bool  `json:"bot"`
}
