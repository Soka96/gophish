package models

// EchoAttach contains the fields and methods for
// an email attachment
type EchoAttach struct {
	Id         int64  `json:"-"`
	EchoEmailId int64  `json:"-"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	Name       string `json:"name"`
}
