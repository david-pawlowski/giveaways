package models

type Giveaway struct {
	Game    string `json:"game"`
	Code    string `json:"code"`
	Claimed bool   `json:"claimed"`
}
