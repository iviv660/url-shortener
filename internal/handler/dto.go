package handler

// ShortenRequest — тело запроса на создание короткой ссылки.
type ShortenRequest struct {
	URL string `json:"url" example:"https://github.com/iviv660"`
}

// ShortenResponse — ответ с коротким кодом и полной короткой ссылкой.
type ShortenResponse struct {
	ShortCode string `json:"short_code" example:"iHwL01z-"`
	ShortURL  string `json:"short_url"  example:"http://localhost:3000/iHwL01z-"`
}

// ErrorResponse — единый формат ошибки.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid request body"`
}
