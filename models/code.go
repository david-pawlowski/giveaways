package models

type Code struct {
	Code    string `json:"code"`
	Claimed bool   `json:"claimed"`
}
