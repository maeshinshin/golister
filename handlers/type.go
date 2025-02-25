package handlers

import "github.com/maeshinshin/go-multiapi/internal/database"

type Handlers struct {
	db database.Service
}

func NewHandlers(db database.Service) *Handlers {
	return &Handlers{db}
}

type HelloWorldresponse struct {
	Message      string  `json:"message"`
	Message2     string  `json:"message2"`
	Message3     string  `json:"message3"`
	ErrorMessage string  `json:"errormessage,omitempty"`
	Error        *Errors `json:"errors"`
}

type Errors struct {
	Error  string `json:"error"`
	Error2 string `json:"error2"`
}

type HealthResponse struct {
}
