package dto

type BookCopyRequest struct {
	Status string `json:"status"`
}

type BookCopyResponse struct {
	ID     uint   `json:"id"`
	Status string `json:"status"`
}
