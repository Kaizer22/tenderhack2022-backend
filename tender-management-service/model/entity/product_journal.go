package entity

// ProductJournal example
type ProductJournal struct {
	RecordID int64 `pg:"id,pk" json:"record_id"`

	ProductId int64    `pg:"product_id" json:"product_id"`
	Product   *Product `pg:"rel:has-one" json:"product,omitempty"`

	QuotationSessionId int64 `pg:"quotation_session_id" json:"session_id"`

	Count int32 `pg:"count" json:"count"`
}

//ProductJournalData example
type ProductJournalData struct {
	ProductId int64 `json:"product_id"`
	Count     int32 `json:"count"`
}
