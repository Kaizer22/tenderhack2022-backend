package entity

import "time"

type SessionStatus string

const (
	StatusActive         = "ACTIVE"
	StatusFinished       = "FINISHED"
	StatusDidntTakePlace = "DIDNT_TAKE_PLACE"
)

type QuotationSession struct {
	ID        int64         `pg:"id,pk" json:"id"`
	Name      string        `pg:"name" json:"name"`
	CreatorId int64         `pg:"creator_id" json:"creator_id"`
	Creator   *Profile      `pg:"rel:has-one" json:"creator,omitempty"`
	Status    SessionStatus `pg:"status" json:"status"`
	//Duration in minutes
	SessionDuration        int               `pg:"session_duration" json:"session_duration"`
	StartPrice             float64           `pg:"start_price" json:"start_price"`
	CurrentPrice           float64           `pg:"current_price" json:"current_price"`
	SessionStepPercent     float64           `pg:"session_step_percent" json:"session_step_percent"`
	StartTime              time.Time         `pg:"start_time" json:"start_time"`
	IsInAdditionalPurchase bool              `pg:"is_in_additional_purchase,use_zero" json:"is_in_additional_purchase"`
	LastBetId              int64             `pg:"last_bet_id" json:"last_bet_id"`
	LastBet                *Bet              `pg:"rel:has-one" json:"last_bet,omitempty"`
	Products               []*ProductJournal `pg:"rel:has-many" json:"products"`
}

type QuotationSessionData struct {
	Name               string                `pg:"name" json:"name"`
	CreatorId          int64                 `json:"creator_id"`
	SessionDuration    int                   `pg:"session_duration" json:"session_duration"`
	StartPrice         float64               `pg:"start_price" json:"start_price"`
	SessionStepPercent float64               `pg:"session_step_percent" json:"session_step_percent"`
	Products           []*ProductJournalData `json:"products"`
}
