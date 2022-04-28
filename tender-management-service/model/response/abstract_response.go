package response

type AbstractResponse struct {
	Code    int    `json:"code"`
	Payload string `json:"payload"`
}
