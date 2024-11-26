package models

type Game struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	DevelopedBy  string `json:"developedBy"`
	PrimaryImage string `json:"primaryImage"`
}
