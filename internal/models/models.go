package models

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Reslut string `json:"result"`
}