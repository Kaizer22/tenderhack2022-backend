package response

type AccountCreated struct {
	Msg       string `json:"msg"`
	Id        int64  `json:"account_id"`
	ProfileId int64  `json:"profile_id"`
}
