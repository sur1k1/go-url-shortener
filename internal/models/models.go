package models

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	Reslut string `json:"result"`
}

type URLData struct {
	UUID string `json:"uuid"`
	ShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
