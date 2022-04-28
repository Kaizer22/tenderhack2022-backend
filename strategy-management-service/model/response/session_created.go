package response

type SessionCreated struct {
	Msg string `json:"msg"`
	Id  int64  `json:"session_id"`
}
