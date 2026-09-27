package handlers

import (
	"fmt"
	"strings"
	"time"
)

type CreateListingRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	City        string  `json:"city"`
}

type CreateListingResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type validationError struct {
	Field string
	Msg   string
}

func (e *validationError) Error() string {
	return fmt.Sprintf("%s:%s", e.Field, e.Msg)
}
func (req CreateListingRequest) validate() error {
	if strings.TrimSpace(req.Title) == "" {
		return &validationError{Field: "title", Msg: "must not be empty"}
	}
	return nil
}
