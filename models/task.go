package models

type Status int

const (
	Incomplete Status = iota
	Complete
)

type Task struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status"`
}
