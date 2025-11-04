package dto

type WebResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}
