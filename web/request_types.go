package web

type genericResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type mountRequest struct {
	MediaID string `json:"id"`
}

type mountResponse struct {
	genericResponse
	Location string `json:"location"`
}

type unmountRequest struct {
	MediaID string `json:"id"`
}

type unmountResponse struct {
	genericResponse
}
