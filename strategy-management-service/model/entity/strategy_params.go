package entity

type StrategyParams struct {
	UserId             int64   `json:"user_id"`
	QuotationSessionId int64   `json:"quotation_session_id"`
	MinimalPrice       float64 `json:"minimal_price"`
	AcceptablePrice    float64 `json:"acceptable_price"`
	PreferablePrice    float64 `json:"preferable_price"`
	Str                string  `json:"strategy"`
}
